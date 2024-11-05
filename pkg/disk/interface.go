package disk

import (
	db "BD/pkg/database"
	"reflect"
)

type Disk interface {
	writeData() (bool, error)
	readData() (bool, error)
	getSortedMassiveByKey() (map[string]db.TableImpl, error)
	convertStringKeyToType(k any, t reflect.Type) (any, error)
	insertData(kvt []keyValueTypeKey) bool
	convertKeyToString(key any) (string, reflect.Type, error)
}
