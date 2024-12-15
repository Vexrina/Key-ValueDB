package http

import (
	"BD/pkg/database"
	"BD/pkg/http/raft"
	u "BD/pkg/http/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

type tableBody struct {
	TableName    string `json:"table_name"`
	NewTableName string `json:"new_table_name"`
}

// post method
func tableCreate(
	w http.ResponseWriter,
	r *http.Request,
	allDbs map[string]database.DataBaseImpl,
	node *raft.RaftNode,
) {
	dbName := r.URL.Query().Get("db_name")
	if dbName == "" {
		u.WriteApiError(w, "db_name parameter is required", http.StatusBadRequest)
		return
	}

	var tB tableBody

	err := json.NewDecoder(r.Body).Decode(&tB)
	if err != nil {
		u.WriteApiError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if tB.TableName == "" {
		u.WriteApiError(w, "you dont provide db name for creation", http.StatusBadRequest)
		return
	}

	db, exist := allDbs[dbName]
	if !exist {
		u.WriteApiError(w, "you provide non-existing db name for creation", http.StatusBadRequest)
		return
	}

	newTable := database.NewTableImpl()
	_, err = db.Create(tB.TableName, *newTable)
	if err != nil {
		u.WriteApiError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if node != nil {
		raft.AppendToLog(node, fmt.Sprintf(
			"DB %s create %s",
			dbName,
			tB.TableName,
		))
	}
	u.WriteApiOK(w, nil, http.StatusCreated)
}

// delete method
func tableDelete(
	w http.ResponseWriter,
	r *http.Request,
	allDbs map[string]database.DataBaseImpl,
	node *raft.RaftNode,
) {
	dbName := r.URL.Query().Get("db_name")
	if dbName == "" {
		u.WriteApiError(w, "db_name parameter is required", http.StatusBadRequest)
		return
	}

	var tB tableBody

	err := json.NewDecoder(r.Body).Decode(&tB)
	if err != nil {
		u.WriteApiError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if tB.TableName == "" {
		u.WriteApiError(w, "you dont provide db name for deleting", http.StatusBadRequest)
		return
	}

	_, exist := allDbs[tB.TableName]
	if !exist {
		u.WriteApiError(w, "you provide non existing db name for deleting", http.StatusBadRequest)
		return
	}

	delete(allDbs, tB.TableName)

	if node != nil {
		raft.AppendToLog(node, fmt.Sprintf(
			"DB %s update %s",
			dbName,
			tB.TableName,
		))
	}

	u.WriteApiOK(w, nil, http.StatusAccepted)
}

func tableRename(
	w http.ResponseWriter,
	r *http.Request,
	allDbs map[string]database.DataBaseImpl,
	node *raft.RaftNode,
) {
	dbName := r.URL.Query().Get("db_name")
	if dbName == "" {
		u.WriteApiError(w, "db_name parameter is required", http.StatusBadRequest)
		return
	}
	var tB tableBody
	err := json.NewDecoder(r.Body).Decode(&tB)
	if err != nil {
		u.WriteApiError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if tB.TableName == "" {
		u.WriteApiError(w, "you dont provide db name for renaming", http.StatusBadRequest)
		return
	}

	db, ok := allDbs[dbName]
	if !ok {
		u.WriteApiError(w, "database not found", http.StatusNotFound)
		return
	}

	_, err = db.Rename(tB.TableName, tB.NewTableName)
	if err != nil {
		u.WriteApiError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if node != nil {
		raft.AppendToLog(node, fmt.Sprintf(
			"DB %s rename %s %s",
			dbName,
			tB.TableName,
			tB.NewTableName,
		))
	}

	u.WriteApiOK(w, nil, http.StatusOK)
}

func tableSelect(
	w http.ResponseWriter,
	r *http.Request,
	allDbs map[string]database.DataBaseImpl,
) {
	dbName := r.URL.Query().Get("db_name")
	if dbName == "" {
		u.WriteApiError(w, "db_name parameter is required", http.StatusBadRequest)
		return
	}
	var tB tableBody
	err := json.NewDecoder(r.Body).Decode(&tB)
	if err != nil {
		u.WriteApiError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if tB.TableName == "" {
		u.WriteApiError(w, "you dont provide db name for renaming", http.StatusBadRequest)
		return
	}

	db, ok := allDbs[dbName]
	if !ok {
		u.WriteApiError(w, "database not found", http.StatusNotFound)
		return
	}

	data, err := db.Select(tB.TableName)
	if err != nil {
		u.WriteApiError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	u.WriteApiOK(w, convertMapToJSONCompatible(data), http.StatusOK)
}
