package raft

import (
	"BD/pkg/database"
	"BD/pkg/parser"
	"BD/pkg/xlog"
	"encoding/gob"
	"fmt"
	"os"
	"sync"
	"time"
)

// состояние рафта
type RaftNode struct {
	// Базовые данные о узле
	ID         string // Идентификатор узла
	State      string // Текущее состояние: "Follower", "Candidate", "Leader"
	Term       int    // Текущий термин
	VotedFor   string // Кандидат, за которого проголосовал узел в текущем терминe
	dirPath    string // дира, где лежит log_file
	log_file   string // файл, который будет читать нода при запуске/писать при падении
	LeaderPeer string
	myPeer     string

	// Для синхронизации состояний
	Mutex             sync.Mutex // Защита состояния
	ResetElectionChan chan bool  // Канал для сброса таймера выборов

	// Для выборов и таймеров
	ElectionTimeout  time.Duration // Таймаут выборов
	HeartbeatTimeout time.Duration // Таймаут heartbeat для лидера

	// Для данных узла
	Log         []LogEntry // Журнал команд
	CommitIndex int        // Индекс последней закоммиченной записи
	LastApplied int        // Индекс последней примененной записи

	// Для репликации и взаимодействия с другими узлами
	Peers  []string                          // Список адресов других узлов
	AllDbs *map[string]database.DataBaseImpl // Данные базы
	parser *parser.ParserImpl                // парсер для применения логов

	NextIndex  map[string]int // узел -> индекс следующей записи
	MatchIndex map[string]int // Последний индекс, подтвержденный каждым узлом
}

type LogEntry struct {
	Index   int    `json:"index"`   // Индекс лога
	Term    int    `json:"term"`    // Термин записи
	Command string `json:"command"` // Команда для выполнения
	Data    any    `json:"data"`    // Данные команды
}

// голосовалка-запрос
type VoteRequest struct {
	Term         int    `json:"term"`         // Текущий термин кандидата
	CandidateID  string `json:"candidateID"`  // Идентификатор кандидата
	LastLogIndex int    `json:"lastLogIndex"` // Индекс последней записи в журнале кандидата
	LastLogTerm  int    `json:"lastLogTerm"`  // Термин последней записи
	LeaderPeer   string `json:"leaderPeer"`   // Откуда пришло сообщение
}

// голосовалка-ответ
type VoteResponse struct {
	Term        int  `json:"term"`        // Текущий термин узла
	VoteGranted bool `json:"voteGranted"` // Голос был предоставлен
}

type AppendEntriesRequest struct {
	Term         int        `json:"term"`         // Текущий термин лидера
	LeaderID     string     `json:"leaderID"`     // ID лидера, отправляющего запрос
	PrevLogIndex int        `json:"prevLogIndex"` // Индекс записи, предшествующей новой записи
	PrevLogTerm  int        `json:"prevLogTerm"`  // Термин записи, предшествующей новой записи
	Entries      []LogEntry `json:"entries"`      // Лог записи для сохранения (может быть пустым для heartbeat)
	LeaderCommit int        `json:"leaderCommit"` // Индекс последней закоммиченной записи у лидера
}

type AppendEntriesResponse struct {
	Term          int  `json:"term"`    // Текущий термин узла (для актуализации лидера)
	Success       bool `json:"success"` // Флаг успешности: true, если запись или heartbeat приняты
	CommitedIndex int  `json:"index"`
}

func sendHeartbeats(raftNode *RaftNode) {
	for raftNode.State == "Leader" {
		for _, peer := range raftNode.Peers {
			go sendAppendEntries(peer, raftNode)
		}
		time.Sleep(raftNode.HeartbeatTimeout)
	}
}

func initializeRaftNode(
	id string,
	peers []string,
	port string,
	pars *parser.ParserImpl,
) *RaftNode {
	node := &RaftNode{
		ID:                id,
		Term:              0,
		State:             "Follower",
		dirPath:           "raft_log_data",
		log_file:          fmt.Sprintf("/%v_raft_state.gob", port),
		Log:               []LogEntry{},
		Peers:             peers,
		Mutex:             sync.Mutex{},
		ElectionTimeout:   3 * time.Second,
		HeartbeatTimeout:  50 * time.Millisecond,
		ResetElectionChan: make(chan bool),
		parser:            pars,
		AllDbs:            &pars.Databases,
		NextIndex:         make(map[string]int, len(peers)),
		MatchIndex:        make(map[string]int, len(peers)),
		LeaderPeer:        "",
		myPeer:            os.Getenv("MYPEER"),
	}
	go electionTimeout(node) // запускаем выборы
	go startApplyLoop(node)  // запускаем применение логов

	gob.Register(LogEntry{})
	gob.Register([]LogEntry{})
	return node
}

func NewRaftNode(id string, peers []string, port string, impl *parser.ParserImpl) *RaftNode {
	if len(peers) < 1 {
		xlog.Error(
			"there must be another server in the system",
			xlog.Field("number of peers", len(peers)),
		)
		panic(
			fmt.Errorf(
				"there must be another server in the system, got %d",
				len(peers),
			),
		)
	}
	return initializeRaftNode(id, peers, port, impl)
}

func saveStateToDisk(node *RaftNode) error {
	xlog.Info("start saving to disk")
	err := os.MkdirAll(node.dirPath, os.ModePerm)
	if err != nil {
		xlog.Error("got error during create dir on disk", xlog.ErrorField(err))
		return err
	}
	file, err := os.Create(node.dirPath + node.log_file)
	if err != nil {
		xlog.Error("got error during save on disk", xlog.ErrorField(err))
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	state := map[string]interface{}{
		"Term":     node.Term,
		"VotedFor": node.VotedFor,
		"Log":      node.Log,
	}
	xlog.Info("end saving to disk")
	return encoder.Encode(state)
}

func loadStateFromDisk(node *RaftNode) error {
	xlog.Info("start loading from disk")
	file, err := os.Open(node.dirPath + node.log_file)
	if err != nil {
		xlog.Error("got error during loading on disk", xlog.ErrorField(err))
		return err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	state := map[string]interface{}{}
	if err := decoder.Decode(&state); err != nil {
		xlog.Error("got error during loading from disk", xlog.ErrorField(err))
		return err
	}

	node.Term = state["Term"].(int)
	node.VotedFor = state["VotedFor"].(string)
	node.Log = state["Log"].([]LogEntry)
	xlog.Info("end loading from disk")
	return nil
}
