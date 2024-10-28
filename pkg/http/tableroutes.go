package http

import (
	"BD/pkg/database"
	"encoding/json"
	"net/http"
)

type tableBody struct {
	TableName string `json:"table_name"`
}

// post method
func tableCreate(
	w http.ResponseWriter,
	r *http.Request,
	allDbs map[string]database.DataBaseImpl,
) {
	dbName := r.URL.Query().Get("db_name")
    if dbName == "" {
        http.Error(w, "db_name parameter is required", http.StatusBadRequest)
        return
    }

	var tB tableBody

	err := json.NewDecoder(r.Body).Decode(&tB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if tB.TableName == "" {
		http.Error(w, "you dont provide db name for creation", http.StatusBadRequest)
		return
	}

	_, exist := allDbs[tB.TableName]
	if exist {
		http.Error(w, "you provide existing db name for creation", http.StatusBadRequest)
		return
	}
	
	allDbs[tB.TableName] = *database.NewDataBaseImpl()

	w.WriteHeader(http.StatusCreated)
}

// delete method
func tableDelete(
	w http.ResponseWriter,
	r *http.Request,
	allDbs map[string]database.DataBaseImpl,
) {
	dbName := r.URL.Query().Get("db_name")
    if dbName == "" {
        http.Error(w, "db_name parameter is required", http.StatusBadRequest)
        return
    }

	var tB tableBody

	err := json.NewDecoder(r.Body).Decode(&tB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if tB.TableName == "" {
		http.Error(w, "you dont provide db name for deleting", http.StatusBadRequest)
		return
	}

	_, exist := allDbs[tB.TableName]
	if !exist {
		http.Error(w, "you provide non existing db name for deleting", http.StatusBadRequest)
		return
	}
	
	delete(allDbs, tB.TableName)

	w.WriteHeader(http.StatusAccepted)
}