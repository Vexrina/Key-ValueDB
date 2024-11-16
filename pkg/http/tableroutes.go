package http

import (
	"BD/pkg/database"
	"encoding/json"
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

	dsb, exist := allDbs[dbName]
	if !exist {
		http.Error(w, "you provide non-existing db name for creation", http.StatusBadRequest)
		return
	}
	table := database.NewTableImpl()
	_, err = dsb.Create(tB.TableName, *table)

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
		http.Error(w, "you dont provide table name for deleting", http.StatusBadRequest)
		return
	}

	db, _ := allDbs[dbName]
	_, err = db.Select(tB.TableName)
	if err != nil {
		http.Error(w, "you provide non existing table name for deleting", http.StatusBadRequest)
		return
	}

	_, err = db.Delete(tB.TableName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusAccepted)
}

func tableRename(w http.ResponseWriter, r *http.Request, allDbs map[string]database.DataBaseImpl) {
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
		http.Error(w, "you dont provide table name for renaming", http.StatusBadRequest)
		return
	}

	if tB.NewTableName == "" {
		http.Error(w, "you dont provide new table name for renaming", http.StatusBadRequest)
		return
	}

	db, ok := allDbs[dbName]
	if !ok {
		http.Error(w, "database not found", http.StatusNotFound)
		return
	}

	_, err = db.Rename(tB.TableName, tB.NewTableName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func tableSelect(w http.ResponseWriter, r *http.Request, allDbs map[string]database.DataBaseImpl) {
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
		http.Error(w, "you dont provide db name for renaming", http.StatusBadRequest)
		return
	}

	db, ok := allDbs[dbName]
	if !ok {
		http.Error(w, "database not found", http.StatusNotFound)
		return
	}

	data, err := db.Select(tB.TableName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
