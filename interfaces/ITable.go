package interfaces

import "time"

type ITable interface {
	add(key interface{}, value interface{}, dayTime time.Time) (bool, error)
	delete(key interface{}) (bool, error)
	put(key interface{}, value interface{}, dayTime time.Time) (bool, error)
	get(key interface{}) (interface{}, error)
	size() int
}
