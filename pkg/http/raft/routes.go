package raft

import (
	"encoding/json"
	"net/http"
)

func VoteHandler(node *RaftNode) {
	http.HandleFunc(
		"/api/internal/raft/vote",
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "not allowed method", http.StatusMethodNotAllowed)
				return
			}
			var voteRequest VoteRequest
			if err := json.NewDecoder(r.Body).Decode(&voteRequest); err != nil {
				http.Error(w, "invalid request body", http.StatusBadRequest)
				return
			}

			response := handleVoteRequest(node, voteRequest)

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(response); err != nil {
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			}
		})
}

func AppendVoteHandler(node *RaftNode) {
	http.HandleFunc(
		"/api/internal/raft/append-entries",
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "not allowed method", http.StatusMethodNotAllowed)
				return
			}
			var appendEntries AppendEntriesRequest
			if err := json.NewDecoder(r.Body).Decode(&appendEntries); err != nil {
				http.Error(w, "invalid request body", http.StatusBadRequest)
				return
			}
			if len(appendEntries.Entries) != 0 {
				handleAppendLog(node, appendEntries)
			}

			response := AppendEntriesResponse{
				Term:    node.Term,
				Success: true,
			}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(response); err != nil {
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			}
		},
	)
}
