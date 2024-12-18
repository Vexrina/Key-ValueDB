package raft

import (
	"BD/pkg/xlog"
)

func becomeFollower(raftNode *RaftNode, newTerm int, leaderPeer string) {
	xlog.Info("im a follower now")
	raftNode.State = "Follower"
	raftNode.Term = newTerm
	raftNode.VotedFor = ""
	raftNode.LeaderPeer = leaderPeer
}

func becomeCandidate(raftNode *RaftNode) {
	xlog.Info("im a Candidate now")
	raftNode.State = "Candidate"
	raftNode.Term += 1
	raftNode.VotedFor = raftNode.ID
	raftNode.LeaderPeer = raftNode.myPeer
	sendRequestVotes(raftNode)
}

func becomeLeader(raftNode *RaftNode) {
	xlog.Info("im a Leader now")
	// Устанавливаем все необходимые параметры для состояния лидера
	raftNode.State = "Leader"

	for _, peer := range raftNode.Peers {
		raftNode.NextIndex[peer] = len(raftNode.Log) + 1
		raftNode.MatchIndex[peer] = 0
	}
	raftNode.LeaderPeer = raftNode.myPeer
	go sendHeartbeats(raftNode)
	if err := saveStateToDisk(raftNode); err != nil {
		xlog.Error("got error during saveStateToDisk", xlog.ErrorField(err))
	}

}
