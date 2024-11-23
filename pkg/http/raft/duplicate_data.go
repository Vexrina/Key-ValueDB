package raft

// import "fmt"

// func applyCommittedEntries(raftNode *RaftNode) {
//     raftNode.Mutex.Lock()
//     defer raftNode.Mutex.Unlock()

//     for raftNode.CommitIndex > raftNode.LastApplied {
//         raftNode.LastApplied++
//         entry := raftNode.Log[raftNode.LastApplied]
//         applyEntry(raftNode, entry)
//     }
// }

// // Пример применения записи к локальному состоянию базы данных
// func applyEntry(raftNode *RaftNode, entry LogEntry) {
//     fmt.Printf("Applying entry at index %d: %+v\n", entry.Index, entry)
//     // Здесь добавляется логика обработки данных в базе `raftNode.AllDbs`
// }
