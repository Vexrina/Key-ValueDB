package http

import (
	"BD/pkg/database"
	"BD/pkg/http/raft"
	"log"
	"BD/pkg/parser"
	"fmt"
	"net/http"
	"os"
)

func databaseHandler(
	allDbs map[string]database.DataBaseImpl,
	opts ...func(w http.ResponseWriter, r *http.Request),
) {
	http.HandleFunc(
		"/api/database",
		func(w http.ResponseWriter, r *http.Request) {
			for _, opt :=range opts{
				opt(w, r)
			}
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
	opts ...func(w http.ResponseWriter, r *http.Request),
) {
	http.HandleFunc(
		"/api/table",
		func(w http.ResponseWriter, r *http.Request) {
			for _, opt :=range opts{
				opt(w, r)
			}
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
	opts ...func(w http.ResponseWriter, r *http.Request),
) {
	http.HandleFunc(
		"/api/key",
		func(w http.ResponseWriter, r *http.Request) {
			for _, opt :=range opts{
				opt(w, r)
			}
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
) {
	
	peers := getPeers()
	log.Println(fmt.Sprintf("find %d peers", len(peers)))

	if len(peers) != 0 {
		log.Println("start with raft node")
		raftNode := raft.NewRaftNode(port, peers, port, p)
		raft.VoteHandler(raftNode)
		raft.AppendVoteHandler(raftNode)
		databaseHandler(p.Databases, raft_check(raftNode))
		tableHandler(p.Databases,raft_check(raftNode))
		keyHandler(p.Databases,  raft_check(raftNode))
	} else {
		log.Println("start without raft node")
		databaseHandler(p.Databases)
		tableHandler(p.Databases)
		keyHandler(p.Databases)
	}
	ping()
	log.Println(fmt.Sprintf("start listening on:%s", port))
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("could not start server: %s\n", err.Error())
	}
}

func raft_check(node *raft.RaftNode) func(w http.ResponseWriter, r *http.Request){
	return func(w http.ResponseWriter, r *http.Request){
		if node.State != "Leader" {
			http.Error(w, "not leader, try another pod", http.StatusForbidden)
		}
	}
}

func getPeers() []string {
	var peers []string
	idx := 1
	for {
		peer := os.Getenv("PEER" + fmt.Sprint(idx))
		if peer == "" {
			break
		}
		log.Printf("find new peer: %s\n", peer)
		peers = append(peers, peer)
		idx += 1
	}
	return peers
}
