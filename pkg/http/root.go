package http

import (
	"BD/pkg/database"
	"BD/pkg/http/raft"
	"BD/pkg/parser"
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
				http.Error(w, "not allowed method", http.StatusMethodNotAllowed)
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
				http.Error(w, "not allowed method", http.StatusMethodNotAllowed)
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
				http.Error(w, "not allowed method", http.StatusMethodNotAllowed)
			}
		},
	)
}

func ping() {
	http.HandleFunc("/api/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "pong")
	})
}

// Run - функция запуска сервера. Запускать через горутину
func Run(
	p *parser.ParserImpl,
	port string,
	peers []string,
) {
	// первостепенно кладем рафт ноду, чтобы у нас сначала запустился таймер выборов и прочего
	// а только потом бд
	// если будет только один сервер выпадет паника(!!!)
	// пока нет идей, как это запустить на одном сервере :) TODO: разобраться
	raftNode := raft.NewRaftNode(port, peers, p)
	raft.VoteHandler(raftNode)
	raft.AppendVoteHandler(raftNode)

	databaseHandler(p.Databases)
	tableHandler(p.Databases)
	keyHandler(p.Databases)
	// метод проверки доступа к сервису
	ping()

	fmt.Printf("start listening on :%s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("could not start server: %s\n", err.Error())
	}
}
