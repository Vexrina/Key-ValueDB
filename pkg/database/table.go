package database

import (
	"errors"
	"sync"
	"time"
)

type Value struct {
	value any
	ttl   time.Time
}

type TableImpl struct {
	data map[any]Value
	Mu   sync.RWMutex
}

func NewTableImpl() *TableImpl {
	return &TableImpl{
		data: make(map[any]Value),
	}
}

func (t *TableImpl) Add(key any, value Value) (bool, error) {
	t.Mu.Lock()
	defer t.Mu.Unlock()

	if t.data == nil {
		t.data = make(map[any]Value)
	}

	if _, exists := t.data[key]; exists {
		return false, errors.New("Такой ключ уже существует")
	}

	t.data[key] = value
	return true, nil
}

func (t *TableImpl) Delete(key any) (bool, error) {
	t.Mu.Lock()
	defer t.Mu.Unlock()
	if _, exists := t.data[key]; !exists {
		return false, errors.New("Такого ключа не существует! Удаление невозможно")
	}
	delete(t.data, key)
	return true, nil
}

func (t *TableImpl) Put(key any, value Value) (bool, error) {
	t.Mu.Lock()
	defer t.Mu.Unlock()
	if _, exists := t.data[key]; !exists {
		return false, errors.New("Такого ключа не существует! Редактирование невозможно")
	}

	t.data[key] = value
	return true, nil
}

func (t *TableImpl) Get(key any) (Value, error) {
	t.Mu.Lock()
	defer t.Mu.Unlock()

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

func (t *TableImpl) Size() int {
	t.Mu.Lock()
	defer t.Mu.Unlock()
	return len(t.data)
}
