package raft

import (
	"BD/pkg/xlog"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func applyCommittedEntries(raftNode *RaftNode) {
	raftNode.Mutex.Lock()
	defer raftNode.Mutex.Unlock()

	for raftNode.CommitIndex > raftNode.LastApplied && raftNode.LastApplied < len(raftNode.Log) {
		entry := raftNode.Log[raftNode.LastApplied]
		applyEntry(raftNode, entry)
		raftNode.LastApplied++
	}
}

func applyEntry(raftNode *RaftNode, entry LogEntry) {
	xlog.Debug("Applying entry at index", xlog.Field("index", entry.Index), xlog.Field("entry", fmt.Sprintf("%+v", entry)))
	_, err := raftNode.parser.Parse(entry.Command)
	if err != nil {
		xlog.Error("raftNode.parser.Parse", xlog.ErrorField(err))
	}
}

func startApplyLoop(raftNode *RaftNode) {
	for {
		time.Sleep(3 * time.Second)
		xlog.Debug("startApplyLoop", xlog.Field("commitIndex", raftNode.CommitIndex), xlog.Field("lastApplied", raftNode.LastApplied), xlog.Field("lenLog", len(raftNode.Log)))
		if raftNode.CommitIndex > raftNode.LastApplied {
			xlog.Info("Apply entries")
			applyCommittedEntries(raftNode)
		}
	}
}

func sendAppendEntries(peer string, raftNode *RaftNode) {
	raftNode.Mutex.Lock()
	prevLogIndex := raftNode.NextIndex[peer] - 1
	prevLogTerm := 0
	if prevLogIndex >= 0 && len(raftNode.Log) > prevLogIndex {
		prevLogTerm = raftNode.Log[prevLogIndex].Term
	}
	xlog.Debug("logEntries", xlog.Field("entries", raftNode.Log))
	var entries []LogEntry
	if len(raftNode.Log) != 0 && raftNode.NextIndex[peer] < len(raftNode.Log) {
		entries = raftNode.Log[raftNode.NextIndex[peer]:]
	}

	leaderCommit := raftNode.CommitIndex
	raftNode.Mutex.Unlock()

	appendEntries := AppendEntriesRequest{
		Term:         raftNode.Term,
		LeaderID:     raftNode.ID,
		PrevLogIndex: prevLogIndex,
		PrevLogTerm:  prevLogTerm,
		Entries:      entries,
		LeaderCommit: leaderCommit,
	}

	// Отправка запроса
	body, err := json.Marshal(appendEntries)
	if err != nil {
		xlog.Error("Failed to marshal AppendEntries request", xlog.ErrorField(err))
		return
	}

	resp, err := http.Post("http://"+peer+"/api/internal/raft/append-entries", "application/json", bytes.NewBuffer(body))
	if err != nil {
		xlog.Error("Failed to send AppendEntries", xlog.Field("destination", peer), xlog.ErrorField(err))
		return
	}
	defer resp.Body.Close()

	var response AppendEntriesResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		xlog.Error("Failed to decode AppendEntries response", xlog.ErrorField(err))
		return
	}

	// Обработка ответа
	raftNode.Mutex.Lock()
	defer raftNode.Mutex.Unlock()

	if response.Success {
		// Обновляем NextIndex и MatchIndex
		raftNode.NextIndex[peer] = response.CommitedIndex
		raftNode.MatchIndex[peer] = raftNode.NextIndex[peer] - 1
	} else if response.Term > raftNode.Term {
		// Узел обнаружил, что его Term устарел
		raftNode.Term = response.Term
		raftNode.State = "Follower"
	} else {
		// Уменьшаем NextIndex для повторной отправки
		raftNode.NextIndex[peer] = max(0, raftNode.NextIndex[peer]-1)
	}
}

func handleAppendLog(node *RaftNode, aer AppendEntriesRequest) AppendEntriesResponse {
	node.Mutex.Lock()
	defer node.Mutex.Unlock()

	// Сравниваем PrevLogIndex и PrevLogTerm
	if aer.PrevLogIndex >= 0 && aer.PrevLogIndex < len(node.Log) {
		if node.Log[aer.PrevLogIndex].Term != aer.PrevLogTerm {
			// Конфликт, откатываем лог
			node.Log = node.Log[:aer.PrevLogIndex+1]
			return AppendEntriesResponse{Term: node.Term, Success: false}
		}
	}

	// Добавляем новые записи
	for i, entry := range aer.Entries {
		idx := aer.PrevLogIndex + 1 + i
		if idx >= len(node.Log) {
			node.Log = append(node.Log, entry)
		} else if node.Log[idx].Term != entry.Term {
			// Заменяем конфликтные записи
			node.Log = node.Log[:idx]
			node.Log = append(node.Log, entry)
		}
	}
	// Обновляем CommitIndex
	if len(node.Log) > node.CommitIndex {
		node.CommitIndex = len(node.Log)
	}

	return AppendEntriesResponse{Term: node.Term, Success: true, CommitedIndex: node.CommitIndex}
}

func AppendToLog(node *RaftNode, command string) {
	node.Mutex.Lock()
	idx := len(node.Log)
	term := node.Term
	node.Log = append(node.Log, LogEntry{idx, term, command, nil})
	node.Mutex.Unlock()
}
