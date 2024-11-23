package http

import (
	"BD/pkg/database"
	"encoding/json"
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
) {
	var dB dbBody
	err := json.NewDecoder(r.Body).Decode(&dB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if dB.DbName == "" {
		http.Error(w, "you dont provide db name for creation", http.StatusBadRequest)
		return
	}
	_, exist := allDbs[dB.DbName]
	if exist {
		http.Error(w, "you provide existing db name for creation", http.StatusBadRequest)
		return
	}
	allDbs[dB.DbName] = *database.NewDataBaseImpl()

	w.WriteHeader(http.StatusCreated)
}

func dataBaseDelete(
	w http.ResponseWriter,
	r *http.Request,
	allDbs map[string]database.DataBaseImpl,
) {
	var dB dbBody
	err := json.NewDecoder(r.Body).Decode(&dB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if dB.DbName == "" {
		http.Error(w, "you dont provide db name for deleting", http.StatusBadRequest)
		return
	}
	_, exist := allDbs[dB.DbName]
	if !exist {
		http.Error(w, "you provide non existing db name for deleting", http.StatusBadRequest)
		return
	}

	delete(allDbs, dB.DbName)

	w.WriteHeader(http.StatusAccepted)
}
