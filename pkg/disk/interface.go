package disk

import db "BD/pkg/database"

type Disk interface {
	writeData() (bool, error)
	readData() (bool, error)
	getSortedMassiveByKey() (map[string]db.TableImpl, error)
}
