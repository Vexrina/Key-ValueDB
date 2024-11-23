package raft

func becomeFollower(raftNode *RaftNode, newTerm int) {
	raftNode.State = "Follower"
	raftNode.Term = newTerm
	raftNode.VotedFor = ""
}

func becomeCandidate(raftNode *RaftNode) {
	raftNode.State = "Candidate"
	raftNode.Term += 1
	raftNode.VotedFor = raftNode.ID
	sendRequestVotes(raftNode)
}

func becomeLeader(raftNode *RaftNode) {
	raftNode.State = "Leader"
	initializeLeaderState(raftNode)
	go sendHeartbeats(raftNode)
}
