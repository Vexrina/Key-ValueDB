package raft

import "log"

func becomeFollower(raftNode *RaftNode, newTerm int) {
	log.Println("im a follower now")
	raftNode.State = "Follower"
	raftNode.Term = newTerm
	raftNode.VotedFor = ""
}

func becomeCandidate(raftNode *RaftNode) {
	log.Println("im a Candidate now")
	raftNode.State = "Candidate"
	raftNode.Term += 1
	raftNode.VotedFor = raftNode.ID
	sendRequestVotes(raftNode)
}

func becomeLeader(raftNode *RaftNode) {
	log.Println("im a Leader now")
	raftNode.State = "Leader"
	initializeLeaderState(raftNode)
	go sendHeartbeats(raftNode)
	saveStateToDisk(raftNode)
}
