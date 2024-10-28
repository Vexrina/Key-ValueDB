package http

import (
	"BD/pkg/database"
	"encoding/json"
	"net/http"
)

type keyBody struct {
	KeyName string `json:"key_name"`
}

// post method
func keyCreate(
	w http.ResponseWriter,
	r *http.Request,
	allDbs map[string]database.DataBaseImpl,
) {
	dbName := r.URL.Query().Get("db_name")
    if dbName == "" {
        http.Error(w, "db_name parameter is required", http.StatusBadRequest)
        return
    }

	tableName := r.URL.Query().Get("table_name")
    if tableName == "" {
        http.Error(w, "db_name parameter is required", http.StatusBadRequest)
        return
    }

	var kB keyBody

	err := json.NewDecoder(r.Body).Decode(&kB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if kB.KeyName == "" {
		http.Error(w, "you dont provide db name for creation", http.StatusBadRequest)
		return
	}

	_, exist := allDbs[kB.KeyName]
	if exist {
		http.Error(w, "you provide existing db name for creation", http.StatusBadRequest)
		return
	}
	
	allDbs[kB.KeyName] = *database.NewDataBaseImpl()

	w.WriteHeader(http.StatusCreated)
}

// delete method
func keyDelete(
	w http.ResponseWriter,
	r *http.Request,
	allDbs map[string]database.DataBaseImpl,
) {
	dbName := r.URL.Query().Get("db_name")
    if dbName == "" {
        http.Error(w, "db_name parameter is required", http.StatusBadRequest)
        return
    }

	tableName := r.URL.Query().Get("table_name")
    if tableName == "" {
        http.Error(w, "db_name parameter is required", http.StatusBadRequest)
        return
    }

	var kB keyBody

	err := json.NewDecoder(r.Body).Decode(&kB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if kB.KeyName == "" {
		http.Error(w, "you dont provide db name for deleting", http.StatusBadRequest)
		return
	}

	_, exist := allDbs[kB.KeyName]
	if !exist {
		http.Error(w, "you provide non existing db name for deleting", http.StatusBadRequest)
		return
	}
	
	delete(allDbs, kB.KeyName)

	w.WriteHeader(http.StatusAccepted)
}