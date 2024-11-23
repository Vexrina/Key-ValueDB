package raft

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func handleVoteRequest(raftNode *RaftNode, voteRequest VoteRequest) VoteResponse {
	raftNode.Mutex.Lock()
	defer raftNode.Mutex.Unlock()

	response := VoteResponse{
		Term:        raftNode.Term,
		VoteGranted: false,
	}

	// Отвергаем запрос, если термин кандидата старше
	if voteRequest.Term < raftNode.Term {
		return response
	}

	// Если термин выше, обновляем свой термин и переходим в Follower
	if voteRequest.Term > raftNode.Term {
		raftNode.Term = voteRequest.Term
		raftNode.State = "Follower"
		raftNode.VotedFor = "" // Сбрасываем голос
	}

	// Проверяем, можем ли проголосовать за кандидата
	if (raftNode.VotedFor == "" || raftNode.VotedFor == voteRequest.CandidateID) &&
		isCandidateLogUpToDate(raftNode, voteRequest.LastLogIndex, voteRequest.LastLogTerm) {
		raftNode.VotedFor = voteRequest.CandidateID
		response.VoteGranted = true
		// Сбрасываем таймер выборов
		raftNode.ResetElectionChan <- true
	}

	return response
}

func isCandidateLogUpToDate(raftNode *RaftNode, candidateLastLogIndex, candidateLastLogTerm int) bool {
	if len(raftNode.Log) == 0 {
		return true // Если лог узла пуст, кандидат актуален
	}

	lastLogIndex := len(raftNode.Log) - 1
	lastLogTerm := raftNode.Log[lastLogIndex].Term

	if candidateLastLogTerm > lastLogTerm {
		return true
	}
	if candidateLastLogTerm == lastLogTerm && candidateLastLogIndex >= lastLogIndex {
		return true
	}

	return false
}

func sendRequestVotes(raftNode *RaftNode) {
	raftNode.Mutex.Lock()
	term := raftNode.Term
	candidateID := raftNode.ID
	lastLogIndex := len(raftNode.Log) - 1
	lastLogTerm := 0
	if lastLogIndex >= 0 {
		lastLogTerm = raftNode.Log[lastLogIndex].Term
	}
	raftNode.Mutex.Unlock()

	voteRequest := VoteRequest{
		Term:         term,
		CandidateID:  candidateID,
		LastLogIndex: lastLogIndex,
		LastLogTerm:  lastLogTerm,
	}

	votes := 1 // Голос за самого себя
	for _, peer := range raftNode.Peers {
		go func(peer string) {
			response := sendVoteRequestToPeer(peer, &voteRequest)
			if response.VoteGranted {
				raftNode.Mutex.Lock()
				votes++
				if votes > len(raftNode.Peers)/2 && raftNode.State == "Candidate" {
					becomeLeader(raftNode)
				}
				raftNode.Mutex.Unlock()
			} else if response.Term > raftNode.Term {
				raftNode.Mutex.Lock()
				becomeFollower(raftNode, response.Term)
				raftNode.Mutex.Unlock()
			}
		}(peer)
	}
}

func electionTimeout(raftNode *RaftNode) {
	for {
		select {
		case <-raftNode.ResetElectionChan:
			// Таймер сброшен, ничего не делаем
		case <-time.After(raftNode.ElectionTimeout):
			raftNode.Mutex.Lock()
			if raftNode.State != "Leader" {
				becomeCandidate(raftNode)
			}
			raftNode.Mutex.Unlock()
		}
	}
}

func sendVoteRequestToPeer(peer string, voteRequest *VoteRequest) VoteResponse {
	body, err := json.Marshal(voteRequest)
	if err != nil {
		fmt.Printf("Failed to marshal VoteRequest: %v\n", err)
		return VoteResponse{}
	}

	resp, err := http.Post("http://"+peer+"/api/raft/request-vote", "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("Failed to send VoteRequest to %s: %v\n", peer, err)
		return VoteResponse{}
	}
	defer resp.Body.Close()

	var response VoteResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		fmt.Printf("Failed to decode VoteResponse: %v\n", err)
		return VoteResponse{}
	}

	return response
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
