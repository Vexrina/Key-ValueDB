package interfaces

import "time"

type ITable interface {
	add(key interface{}, value interface{}, dayTime time.Time)
	delete(key interface{})
	put(key interface{}, value interface{}, dayTime time.Time)
	get(key interface{}) interface{}
	size() int
}
