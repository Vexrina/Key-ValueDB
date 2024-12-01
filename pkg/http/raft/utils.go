package raft

import (
	"BD/pkg/database"
	"BD/pkg/parser"
	"fmt"
	"sync"
	"time"
)

// состояние рафта
type RaftNode struct {
	// Базовые данные о узле
	ID       string // Идентификатор узла
	State    string // Текущее состояние: "Follower", "Candidate", "Leader"
	Term     int    // Текущий термин
	VotedFor string // Кандидат, за которого проголосовал узел в текущем терминe

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
	Term    int  `json:"term"`    // Текущий термин узла (для актуализации лидера)
	Success bool `json:"success"` // Флаг успешности: true, если запись или heartbeat приняты
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
	pars *parser.ParserImpl,
) *RaftNode {
	node := &RaftNode{
		ID:                id,
		Term:              0,
		State:             "Follower",
		Log:               []LogEntry{},
		Peers:             peers,
		Mutex:             sync.Mutex{},
		ElectionTimeout:   150 * time.Millisecond,
		HeartbeatTimeout:  50 * time.Millisecond,
		ResetElectionChan: make(chan bool),
		parser:            pars,
		AllDbs:            &pars.Databases,
	}
	go electionTimeout(node) // запускаем выборы
	startApplyLoop(node)     // запускаем применение логов
	return node
}

func initializeLeaderState(raftNode *RaftNode) {
	raftNode.Mutex.Lock()
	defer raftNode.Mutex.Unlock()

	// Устанавливаем все необходимые параметры для состояния лидера
	raftNode.State = "Leader"

	// Инициализируем CommitIndex и LastApplied, если требуется
	raftNode.CommitIndex = 0
	raftNode.LastApplied = 0

	// Отправляем heartbeat всем узлам
	go func() {
		ticker := time.NewTicker(raftNode.HeartbeatTimeout)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				for _, peer := range raftNode.Peers {
					go sendAppendEntries(peer, raftNode)
				}
			case <-raftNode.ResetElectionChan:
				// Если сброс произошел, завершаем отправку heartbeat
				return
			}
		}
	}()
}

func NewRaftNode(id string, peers []string, impl *parser.ParserImpl) *RaftNode {
	if len(peers) < 1 {
		panic(
			fmt.Errorf(
				"there must be another server in the system, got %d",
				len(peers),
			),
		)
	}
	return initializeRaftNode(id, peers, impl)
}
