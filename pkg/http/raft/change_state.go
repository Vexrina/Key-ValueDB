package raft

import (
	"BD/pkg/xlog"
)

func becomeFollower(raftNode *RaftNode, newTerm int) {
	xlog.Info("im a follower now")
	raftNode.State = "Follower"
	raftNode.Term = newTerm
	raftNode.VotedFor = ""
}

func becomeCandidate(raftNode *RaftNode) {
	xlog.Info("im a Candidate now")
	raftNode.State = "Candidate"
	raftNode.Term += 1
	raftNode.VotedFor = raftNode.ID
	sendRequestVotes(raftNode)
}

func becomeLeader(raftNode *RaftNode) {
	xlog.Info("im a Leader now")
	raftNode.State = "Leader"
	initializeLeaderState(raftNode)
	go sendHeartbeats(raftNode)
	saveStateToDisk(raftNode)
}
