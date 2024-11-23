package raft

import (
	"BD/pkg/database"
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
	Peers  []string                         // Список адресов других узлов
	AllDbs map[string]database.DataBaseImpl // Данные базы
}

type LogEntry struct {
	Term    int    // Термин записи
	Command string // Команда для выполнения
	Data    any    // Данные команды
}

// голосовалка-запрос
type VoteRequest struct {
	Term         int    // Текущий термин кандидата
	CandidateID  string // Идентификатор кандидата
	LastLogIndex int    // Индекс последней записи в журнале кандидата
	LastLogTerm  int    // Термин последней записи
}

// голосовалка-ответ
type VoteResponse struct {
	Term        int  // Текущий термин узла
	VoteGranted bool // Голос был предоставлен
}

type AppendEntriesRequest struct {
	Term         int        `json:"term"`           // Текущий термин лидера
	LeaderID     string     `json:"leader_id"`      // ID лидера, отправляющего запрос
	PrevLogIndex int        `json:"prev_log_index"` // Индекс записи, предшествующей новой записи
	PrevLogTerm  int        `json:"prev_log_term"`  // Термин записи, предшествующей новой записи
	Entries      []LogEntry `json:"entries"`        // Лог записи для сохранения (может быть пустым для heartbeat)
	LeaderCommit int        `json:"leader_commit"`  // Индекс последней закоммиченной записи у лидера
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

func initializeRaftNode(id string, peers []string) *RaftNode {
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
	}
	go electionTimeout(node)
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
