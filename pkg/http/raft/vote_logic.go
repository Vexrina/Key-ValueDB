package raft

import (
	"BD/pkg/xlog"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// handleVoteRequest - Основной хендлер для голосования. Валидируем кандидата и (не) голосуем за него
func handleVoteRequest(raftNode *RaftNode, voteRequest VoteRequest) VoteResponse {
	raftNode.Mutex.Lock()
	defer raftNode.Mutex.Unlock()
	xlog.Info("Start handling vote request", xlog.Field("voteRequest", fmt.Sprintf("%+v", voteRequest)))
	response := VoteResponse{
		Term:        raftNode.Term,
		VoteGranted: false,
	}

	// Отвергаем запрос, если термин кандидата старше
	if voteRequest.Term < raftNode.Term {
		xlog.Info("Decline vote request")
		return response
	}

	// Если термин выше, обновляем свой термин и переходим в Follower
	if voteRequest.Term > raftNode.Term {
		becomeFollower(raftNode, voteRequest.Term, voteRequest.LeaderPeer)
	}

	// Проверяем, можем ли проголосовать за кандидата
	if canVote(raftNode, &voteRequest) {
		xlog.Info("Vote for candidate", xlog.Field("candidateID", voteRequest.CandidateID))
		raftNode.VotedFor = voteRequest.CandidateID
		response.VoteGranted = true
		// Сбрасываем таймер выборов
		raftNode.LeaderPeer = voteRequest.LeaderPeer
		raftNode.ResetElectionChan <- true
	}

	return response
}

func canVote(node *RaftNode, request *VoteRequest) bool {
	firstPart := node.VotedFor == "" || node.VotedFor == request.CandidateID
	secondPart := isCandidateLogUpToDate(node, request.LastLogIndex, request.LastLogTerm)
	return firstPart && secondPart
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
	leaderPeer := raftNode.myPeer

	raftNode.Mutex.Unlock()

	voteRequest := VoteRequest{
		Term:         term,
		CandidateID:  candidateID,
		LastLogIndex: lastLogIndex,
		LastLogTerm:  lastLogTerm,
		LeaderPeer:   leaderPeer,
	}
	xlog.Info("start requesting votes")
	votes := 1 // Голос за самого себя
	for _, peer := range raftNode.Peers {
		go func(peer string) {
			xlog.Info("send request to peer", xlog.Field("peer", peer))
			response := sendVoteRequestToPeer(peer, &voteRequest)
			if response.VoteGranted {
				xlog.Info("got voteGranted from peer", xlog.Field("peer", peer))
				raftNode.Mutex.Lock()
				votes++
				if votes > len(raftNode.Peers)/2 && raftNode.State == "Candidate" {
					becomeLeader(raftNode)
				}
				raftNode.Mutex.Unlock()
			} else if response.Term > raftNode.Term {
				xlog.Info("got greater term")
				raftNode.Mutex.Lock()
				becomeFollower(raftNode, response.Term, voteRequest.LeaderPeer)
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
			if raftNode.State != "Leader" {
				becomeCandidate(raftNode)
			}
		}
	}
}

func sendVoteRequestToPeer(peer string, voteRequest *VoteRequest) VoteResponse {
	body, err := json.Marshal(voteRequest)
	if err != nil {
		xlog.Error("Failed to marshal VoteRequest", xlog.ErrorField(err))
		return VoteResponse{}
	}

	resp, err := http.Post("http://"+peer+"/api/internal/raft/vote", "application/json", bytes.NewBuffer(body))
	if err != nil {
		xlog.Error("Failed to send VoteRequest", xlog.Field("peer direction", peer), xlog.ErrorField(err))
		return VoteResponse{}
	}
	defer resp.Body.Close()

	var response VoteResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		xlog.Error("Failed to decode VoteResponse", xlog.ErrorField(err))
		return VoteResponse{}
	}

	return response
}
