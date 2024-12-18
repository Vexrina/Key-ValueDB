package http

import (
	"BD/pkg/database"
	"BD/pkg/http/raft"
	u "BD/pkg/http/utils"
	"BD/pkg/parser"
	"BD/pkg/xlog"
	"fmt"
	"net/http"
)

func databaseHandler(
	allDbs map[string]database.DataBaseImpl,
	node *raft.RaftNode,
) {
	http.HandleFunc(
		"/api/database",
		func(w http.ResponseWriter, r *http.Request) {
			if raft_check(node, w) {
				return
			}
			switch r.Method {
			case http.MethodPost:
				databaseCreate(w, r, allDbs, node)
			case http.MethodDelete:
				dataBaseDelete(w, r, allDbs, node)
			default:
				u.WriteApiError(w, "not allowed method", http.StatusMethodNotAllowed)
			}
		},
	)
}

func tableHandler(
	allDbs map[string]database.DataBaseImpl,
	node *raft.RaftNode,
) {
	http.HandleFunc(
		"/api/table",
		func(w http.ResponseWriter, r *http.Request) {
			if raft_check(node, w) {
				return
			}
			switch r.Method {
			case http.MethodPost:
				tableCreate(w, r, allDbs, node)
			case http.MethodDelete:
				tableDelete(w, r, allDbs, node)
			case http.MethodPut:
				tableRename(w, r, allDbs, node)
			case http.MethodGet:
				tableSelect(w, r, allDbs)
			default:
				u.WriteApiError(w, "not allowed method", http.StatusMethodNotAllowed)
			}
		},
	)
}

func keyHandler(
	allDbs map[string]database.DataBaseImpl,
	node *raft.RaftNode,
) {
	http.HandleFunc(
		"/api/key",
		func(w http.ResponseWriter, r *http.Request) {
			if raft_check(node, w) {
				return
			}
			switch r.Method {
			case http.MethodPost:
				keyInsert(w, r, allDbs, node)
			case http.MethodDelete:
				keyDelete(w, r, allDbs, node)
			case http.MethodGet:
				keyGet(w, r, allDbs)
			case http.MethodPut:
				keyUpdate(w, r, allDbs, node)
			default:
				u.WriteApiError(w, "not allowed method", http.StatusMethodNotAllowed)
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
func Run(p *parser.ParserImpl, port string) {
	peers := getPeers()
	xlog.Debug("", xlog.Field("find number of peers", len(peers)))
	var raftNode *raft.RaftNode = nil
	if len(peers) != 0 {
		xlog.Info("start with raft node")
		raftNode = raft.NewRaftNode(port, peers, port, p)
		raft.VoteHandler(raftNode)
		raft.AppendEntriesHandler(raftNode)
	}

	databaseHandler(p.Databases, raftNode)
	tableHandler(p.Databases, raftNode)
	keyHandler(p.Databases, raftNode)
	ping()

	xlog.Info(fmt.Sprintf("start listening on:%s", port))
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("could not start server: %s\n", err.Error())
	}
}
