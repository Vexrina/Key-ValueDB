package http

import (
	"BD/pkg/database"
	"BD/pkg/http/raft"
	u "BD/pkg/http/utils"
	"BD/pkg/xlog"
	"fmt"
	"net/http"
	"os"
)

func raft_check(node *raft.RaftNode, w http.ResponseWriter) bool {
	if node != nil && node.State == "Follower" {
		u.WriteApiError(w,
			fmt.Sprintf(
				"now im not a leader or candidate. Leader is: %s. MyState is %s. firstCheck: %s, SecondCheck: %s",
				node.LeaderPeer,
				node.State,
				fmt.Sprint(node != nil),
				fmt.Sprint(node.State != "Follower"),
			),
			http.StatusForbidden,
		)
		return true
	}
	return false
}

func getPeers() []string {
	var peers []string
	idx := 1
	for {
		peer := os.Getenv("PEER" + fmt.Sprint(idx))
		if peer == "" {
			break
		}
		xlog.Debug("find new peer", xlog.Field(fmt.Sprintf("Peer %d", idx), peer))
		peers = append(peers, peer)
		idx += 1
	}
	return peers
}

func convertMapToJSONCompatible(input database.TableImpl) map[string]database.Value {
	result := make(map[string]database.Value)
	for key, value := range input.DataTable {
		strKey := fmt.Sprintf("%v", key)
		result[strKey] = value
	}
	return result
}
