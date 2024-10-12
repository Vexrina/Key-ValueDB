package interfaces

import "time"

type ITable interface {
	add(key interface{}, value Value) (bool, error)
	delete(key interface{}) (bool, error)
	put(key interface{}, value Value) (bool, error)
	get(key interface{}) (interface{}, error)
	size() int
}

type Value struct {
	value interface{}
	ttl   time.Time
}
