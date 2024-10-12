package interfaces

import (
	"errors"
	"sync"
	"time"
)

type ITable interface {
	add(key interface{}, value Value) (bool, error)
	delete(key interface{}) (bool, error)
	put(key interface{}, value Value) (bool, error)
	get(key interface{}) (Value, error)
	size() int
}

type Value struct {
	value interface{}
	ttl   time.Time
}

type Table struct {
	data map[interface{}]Value
	Mu   sync.RWMutex
}

func (t *Table) add(key interface{}, value Value) (bool, error) {
	t.Mu.Lock()
	defer t.Mu.Unlock()

	if t.data == nil {
		t.data = make(map[interface{}]Value)
	}

	if _, exists := t.data[key]; exists {
		return false, errors.New("Такой ключ уже существует")
	}

	t.data[key] = value
	return true, nil
}

func (t *Table) delete(key interface{}) (bool, error) {
	t.Mu.Lock()
	defer t.Mu.Unlock()
	if _, exists := t.data[key]; !exists {
		return false, errors.New("Такого ключа не существует! Удаление невозможно")
	}
	delete(t.data, key)
	return true, nil
}

func (t *Table) put(key interface{}, value Value) (bool, error) {
	t.Mu.Lock()
	defer t.Mu.Unlock()
	if _, exists := t.data[key]; !exists {
		return false, errors.New("Такого ключа не существует! Редактирование невозможно")
	}

	t.data[key] = value
	return true, nil
}

func (t *Table) get(key interface{}) (Value, error) {
	t.Mu.RLock()
	defer t.Mu.RUnlock()

	value, exists := t.data[key]
	if !exists {
		return Value{}, errors.New("Ключа не существует!")
	}

	if time.Now().After(value.ttl) {
		t.Mu.RUnlock()
		t.Mu.Lock()
		delete(t.data, key)
		t.Mu.Unlock()
		t.Mu.RLock()
		return Value{}, errors.New("Ключ был удален")
	}

	return value, nil
}

func (t *Table) size() int {
	t.Mu.RLock()
	defer t.Mu.RUnlock()
	return len(t.data)
}
