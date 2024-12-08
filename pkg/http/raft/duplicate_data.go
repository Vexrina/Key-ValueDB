package raft

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func applyCommittedEntries(raftNode *RaftNode) {
	raftNode.Mutex.Lock()
	defer raftNode.Mutex.Unlock()

	for raftNode.CommitIndex > raftNode.LastApplied {
		raftNode.LastApplied++
		entry := raftNode.Log[raftNode.LastApplied]
		applyEntry(raftNode, entry)
	}
}

func applyEntry(raftNode *RaftNode, entry LogEntry) {
	fmt.Printf("Applying entry at index %d: %+v\n", entry.Index, entry)
	_, err := raftNode.parser.Parse(entry.Command)
	if err != nil {
		fmt.Printf("raftNode.parser.Parse: %v", err)
	}
}

func startApplyLoop(raftNode *RaftNode) {
	go func() {
		for {
			time.Sleep(10 * time.Millisecond)
			if raftNode.CommitIndex > raftNode.LastApplied {
				applyCommittedEntries(raftNode)
			}
		}
	}()
}

func sendAppendEntries(peer string, raftNode *RaftNode) {
	raftNode.Mutex.Lock()
	prevLogIndex := raftNode.NextIndex[peer] - 1
	prevLogTerm := 0
	if prevLogIndex >= 0 {
		prevLogTerm = raftNode.Log[prevLogIndex].Term
	}
	entries := raftNode.Log[raftNode.NextIndex[peer]:]
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
		fmt.Printf("Failed to marshal AppendEntries request: %v\n", err)
		return
	}

	resp, err := http.Post("http://"+peer+"/api/raft/append-entries", "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("Failed to send AppendEntries to %s: %v\n", peer, err)
		return
	}
	defer resp.Body.Close()

	var response AppendEntriesResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		fmt.Printf("Failed to decode AppendEntries response: %v\n", err)
		return
	}

	// Обработка ответа
	raftNode.Mutex.Lock()
	defer raftNode.Mutex.Unlock()

	if response.Success {
		// Обновляем NextIndex и MatchIndex
		raftNode.NextIndex[peer] = prevLogIndex + len(entries) + 1
		raftNode.MatchIndex[peer] = raftNode.NextIndex[peer] - 1
	} else if response.Term > raftNode.Term {
		// Узел обнаружил, что его Term устарел
		raftNode.Term = response.Term
		raftNode.State = "Follower"
	} else {
		// Уменьшаем NextIndex для повторной отправки
		raftNode.NextIndex[peer] = max(1, raftNode.NextIndex[peer]-1)
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
	if aer.LeaderCommit > node.CommitIndex {
		node.CommitIndex = min(aer.LeaderCommit, len(node.Log)-1)
	}

	return AppendEntriesResponse{Term: node.Term, Success: true}
}
