package http

import (
	"BD/pkg/database"
	"BD/pkg/http/raft"
	u "BD/pkg/http/utils"
	"BD/pkg/xlog"
	"encoding/json"
	"fmt"
	"net/http"
)

type keyBody struct {
	KeyName string         `json:"key_name"`
	Value   database.Value `json:"value"`
}

// post method
func keyInsert(
	w http.ResponseWriter,
	r *http.Request,
	allDbs map[string]database.DataBaseImpl,
	node *raft.RaftNode,
) {
	dbName := r.URL.Query().Get("db_name")
	if dbName == "" {
		xlog.Error("dont provided db_name")
		u.WriteApiError(w, "db_name parameter is required", http.StatusBadRequest)
		return
	}

	tableName := r.URL.Query().Get("table_name")
	if tableName == "" {
		xlog.Error("dont provided table_name")
		u.WriteApiError(w, "table_name parameter is required", http.StatusBadRequest)
		return
	}

	var kB keyBody

	err := json.NewDecoder(r.Body).Decode(&kB)
	if err != nil {
		xlog.Error("Cannot decode request body ot keyBody", xlog.ErrorField(err))
		u.WriteApiError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if kB.KeyName == "" {
		xlog.Error("dont provided keyNmae")
		u.WriteApiError(w, "you dont provide key name for creation", http.StatusBadRequest)
		return
	}

	db, ok := allDbs[dbName]
	if !ok {
		xlog.Error("database not found")
		u.WriteApiError(w, "database not found", http.StatusNotFound)
		return
	}

	impl, err := db.Select(tableName)
	if err != nil {
		xlog.Error("got error during table select", xlog.ErrorField(err))
		u.WriteApiError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = impl.Insert(kB.KeyName, kB.Value)
	if err != nil {
		xlog.Error("got error during inserting", xlog.ErrorField(err))
		u.WriteApiError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if node != nil {
		raft.AppendToLog(node, fmt.Sprintf(
			"Table %s insert %s %s %s %s",
			dbName,
			tableName,
			kB.KeyName,
			kB.Value.Val,
			kB.Value.Ttl,
		))
	}
	u.WriteApiOK(w, nil, http.StatusOK)
}

// delete method
func keyDelete(
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

	tableName := r.URL.Query().Get("table_name")
	if tableName == "" {
		u.WriteApiError(w, "db_name parameter is required", http.StatusBadRequest)
		return
	}

	var kB keyBody

	err := json.NewDecoder(r.Body).Decode(&kB)
	if err != nil {
		u.WriteApiError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if kB.KeyName == "" {
		u.WriteApiError(w, "you dont provide db name for deleting", http.StatusBadRequest)
		return
	}

	db, ok := allDbs[dbName]
	if !ok {
		u.WriteApiError(w, "db not exist", http.StatusNotFound)
		return
	}

	ok, err = db.Delete(kB.KeyName)
	if err != nil {
		u.WriteApiError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !ok {
		u.WriteApiError(w, "something went wrong, the key was not deleted", http.StatusInternalServerError)
		return
	}

	if node != nil {
		raft.AppendToLog(node, fmt.Sprintf(
			"Table %s delete %s %s",
			dbName,
			tableName,
			kB.KeyName,
		))
	}

	u.WriteApiOK(w, nil, http.StatusOK)
}

func keyGet(
	w http.ResponseWriter,
	r *http.Request,
	allDbs map[string]database.DataBaseImpl,
) {
	dbName := r.URL.Query().Get("db_name")
	if dbName == "" {
		u.WriteApiError(w, "db_name parameter is required", http.StatusBadRequest)
		return
	}

	tableName := r.URL.Query().Get("table_name")
	if tableName == "" {
		u.WriteApiError(w, "table_name parameter is required", http.StatusBadRequest)
		return
	}

	var kB keyBody
	err := json.NewDecoder(r.Body).Decode(&kB)
	if err != nil {
		u.WriteApiError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if kB.KeyName == "" {
		u.WriteApiError(w, "KeyName is required", http.StatusBadRequest)
		return
	}

	db, ok := allDbs[dbName]
	if !ok {
		u.WriteApiError(w, "database not found", http.StatusNotFound)
		return
	}

	impl, err := db.Select(tableName)
	if err != nil {
		u.WriteApiError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := impl.Get(kB.KeyName)
	if err != nil {
		u.WriteApiError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	u.WriteApiOK(w, data, http.StatusOK)
}

func keyUpdate(
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

	tableName := r.URL.Query().Get("table_name")
	if tableName == "" {
		u.WriteApiError(w, "table_name parameter is required", http.StatusBadRequest)
		return
	}

	var kB keyBody
	err := json.NewDecoder(r.Body).Decode(&kB)
	if err != nil {
		u.WriteApiError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if kB.KeyName == "" {
		u.WriteApiError(w, "KeyName is required", http.StatusBadRequest)
		return
	}

	db, ok := allDbs[dbName]
	if !ok {
		u.WriteApiError(w, "database not found", http.StatusNotFound)
		return
	}

	impl, err := db.Select(tableName)
	if err != nil {
		u.WriteApiError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = impl.Update(kB.KeyName, kB.Value)
	if err != nil {
		u.WriteApiError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if node != nil {
		raft.AppendToLog(node, fmt.Sprintf(
			"Table %s update %s %s",
			dbName,
			tableName,
			kB.KeyName,
		))
	}

	u.WriteApiOK(w, nil, http.StatusOK)
}
