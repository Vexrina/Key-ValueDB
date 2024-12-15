package http

import (
	"BD/pkg/database"
	"BD/pkg/http/raft"
	u "BD/pkg/http/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

type dbBody struct {
	DbName string `json:"db_name"`
}

// post method
func databaseCreate(
	w http.ResponseWriter,
	r *http.Request,
	allDbs map[string]database.DataBaseImpl,
	node *raft.RaftNode,
) {
	var dB dbBody
	err := json.NewDecoder(r.Body).Decode(&dB)
	if err != nil {
		u.WriteApiError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if dB.DbName == "" {
		u.WriteApiError(w, "you dont provide db name for creation", http.StatusBadRequest)
		return
	}
	_, exist := allDbs[dB.DbName]
	if exist {
		u.WriteApiError(w, "you provide existing db name for creation", http.StatusBadRequest)
		return
	}
	allDbs[dB.DbName] = *database.NewDataBaseImpl()

	if node != nil {
		raft.AppendToLog(node, fmt.Sprintf("SYSTEM create %s", dB.DbName))
	}
	u.WriteApiOK(w, nil, http.StatusCreated)
}

func dataBaseDelete(
	w http.ResponseWriter,
	r *http.Request,
	allDbs map[string]database.DataBaseImpl,
	node *raft.RaftNode,
) {
	var dB dbBody
	err := json.NewDecoder(r.Body).Decode(&dB)
	if err != nil {
		u.WriteApiError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if dB.DbName == "" {
		u.WriteApiError(w, "you dont provide db name for deleting", http.StatusBadRequest)
		return
	}
	_, exist := allDbs[dB.DbName]
	if !exist {
		u.WriteApiError(w, "you provide non existing db name for deleting", http.StatusBadRequest)
		return
	}

	delete(allDbs, dB.DbName)
	if node != nil {
		raft.AppendToLog(node, fmt.Sprintf("SYSTEM delete %s", dB.DbName))
	}
	u.WriteApiOK(w, nil, http.StatusAccepted)
}
