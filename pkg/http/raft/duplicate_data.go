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
	appendEntries := AppendEntriesRequest{
		Term:         raftNode.Term,
		LeaderID:     raftNode.ID,
		PrevLogIndex: len(raftNode.Log) - 1,
		PrevLogTerm:  raftNode.Log[len(raftNode.Log)-1].Term,
		Entries:      []LogEntry{}, // Пустое тело для heartbeat
		LeaderCommit: raftNode.CommitIndex,
	}

	// Отправка запроса с помощью HTTP
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

	raftNode.Mutex.Lock()
	defer raftNode.Mutex.Unlock()

	// Обработка ответа
	if response.Success {
		// Успешно синхронизировано
	} else if response.Term > raftNode.Term {
		// Узел обнаружил, что его Term устарел
		raftNode.Term = response.Term
		raftNode.State = "Follower"
	}
}

func handleAppendLog(node *RaftNode, aer AppendEntriesRequest) {
	for idx := node.CommitIndex; idx < len(aer.Entries); idx++ {
		node.Log = append(node.Log, aer.Entries[idx])
		node.CommitIndex = idx
	}
}
