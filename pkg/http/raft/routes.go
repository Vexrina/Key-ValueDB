package raft

import (
	u "BD/pkg/http/utils"
	"encoding/json"
	"net/http"
)

func VoteHandler(node *RaftNode) {
	go http.HandleFunc(
		"/api/internal/raft/vote",
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				u.WriteApiError(w, "not allowed method", http.StatusMethodNotAllowed)
				return
			}
			var voteRequest VoteRequest
			if err := json.NewDecoder(r.Body).Decode(&voteRequest); err != nil {
				u.WriteApiError(w, "invalid request body", http.StatusBadRequest)
				return
			}

			response := handleVoteRequest(node, voteRequest)

			u.WriteApiOK(w, response, http.StatusAccepted)
		})
}

func AppendEntriesHandler(node *RaftNode) {
	go http.HandleFunc(
		"/api/internal/raft/append-entries",
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				u.WriteApiError(w, "not allowed method", http.StatusMethodNotAllowed)
				return
			}
			var appendEntries AppendEntriesRequest
			if err := json.NewDecoder(r.Body).Decode(&appendEntries); err != nil {
				u.WriteApiError(w, "invalid request body", http.StatusBadRequest)
				return
			}
			if len(appendEntries.Entries) != 0 {
				response := handleAppendLog(node, appendEntries)
				u.WriteApiOK(w, response, http.StatusAccepted)
				return
			}

			response := AppendEntriesResponse{
				Term:    node.Term,
				Success: true,
			}
			u.WriteApiOK(w, response, http.StatusAccepted)
		},
	)
}
