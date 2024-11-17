package http

import (
	"BD/pkg/database"
	"fmt"
	"net/http"
)

func databaseHandler(
	allDbs map[string]database.DataBaseImpl,
) {
	http.HandleFunc(
		"/api/database",
		func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				databaseCreate(w, r, allDbs)
			case http.MethodDelete:
				dataBaseDelete(w, r, allDbs)
			default:
				http.Error(w, "not allowed method", http.StatusBadRequest)
			}
		},
	)
}

func tableHandler(
	allDbs map[string]database.DataBaseImpl,
) {
	http.HandleFunc(
		"/api/table",
		func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				tableCreate(w, r, allDbs)
			case http.MethodDelete:
				tableDelete(w, r, allDbs)
			case http.MethodPut:
				tableRename(w, r, allDbs)
			case http.MethodGet:
				tableSelect(w, r, allDbs)
			default:
				http.Error(w, "not allowed method", http.StatusBadRequest)
			}
		},
	)
}

func keyHandler(
	allDbs map[string]database.DataBaseImpl,
) {
	http.HandleFunc(
		"/api/key",
		func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				keyInsert(w, r, allDbs)
			case http.MethodDelete:
				keyDelete(w, r, allDbs)
			case http.MethodGet:
				keyGet(w, r, allDbs)
			case http.MethodPut:
				keyUpdate(w, r, allDbs)
			default:
				http.Error(w, "not allowed method", http.StatusBadRequest)
			}
		},
	)
}

func ping() {
	http.HandleFunc("/api/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "pong")
	})
}

func Run(allDbs map[string]database.DataBaseImpl, port string) {
	databaseHandler(allDbs)
	tableHandler(allDbs)
	keyHandler(allDbs)
	ping()

	fmt.Printf("start listening on :%s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("could not start server: %s\n", err.Error())
	}
}
