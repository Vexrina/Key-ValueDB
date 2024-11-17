package http

import (
	"BD/pkg/database"
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

	db, ok := allDbs[dbName]
	if !ok {
		http.Error(w, "database not found", http.StatusNotFound)
		return
	}

	impl, err := db.Select(tableName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = impl.Insert(kB.KeyName, kB.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
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

	db, ok := allDbs[dbName]
	if !ok {
		http.Error(w, "db not exist", http.StatusNotFound)
	}
	ok, err = db.Delete(kB.KeyName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if !ok {
		fmt.Fprint(w, "something went wrong, the key was not deleted")
	}

	//_, exist := allDbs[kB.KeyName]
	//if !exist {
	//	http.Error(w, "you provide non existing db name for deleting", http.StatusBadRequest)
	//	return
	//}
	//
	//delete(allDbs, kB.KeyName)

	w.WriteHeader(http.StatusOK)
}

func keyGet(w http.ResponseWriter, r *http.Request, allDbs map[string]database.DataBaseImpl) {
	dbName := r.URL.Query().Get("db_name")
	if dbName == "" {
		http.Error(w, "db_name parameter is required", http.StatusBadRequest)
		return
	}

	tableName := r.URL.Query().Get("table_name")
	if tableName == "" {
		http.Error(w, "table_name parameter is required", http.StatusBadRequest)
		return
	}

	var kB keyBody
	err := json.NewDecoder(r.Body).Decode(&kB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if kB.KeyName == "" {
		http.Error(w, "KeyName is required", http.StatusBadRequest)
		return
	}

	db, ok := allDbs[dbName]
	if !ok {
		http.Error(w, "database not found", http.StatusNotFound)
		return
	}

	impl, err := db.Select(tableName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := impl.Get(kB.KeyName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func keyUpdate(w http.ResponseWriter, r *http.Request, allDbs map[string]database.DataBaseImpl) {
	dbName := r.URL.Query().Get("db_name")
	if dbName == "" {
		http.Error(w, "db_name parameter is required", http.StatusBadRequest)
		return
	}

	tableName := r.URL.Query().Get("table_name")
	if tableName == "" {
		http.Error(w, "table_name parameter is required", http.StatusBadRequest)
		return
	}

	var kB keyBody
	err := json.NewDecoder(r.Body).Decode(&kB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if kB.KeyName == "" {
		http.Error(w, "KeyName is required", http.StatusBadRequest)
		return
	}

	db, ok := allDbs[dbName]
	if !ok {
		http.Error(w, "database not found", http.StatusNotFound)
		return
	}

	impl, err := db.Select(tableName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = impl.Update(kB.KeyName, kB.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
