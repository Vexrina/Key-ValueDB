package database

import (
	"errors"
	"time"
)

type Value struct {
	Val any
	Ttl time.Time
}

type TableImpl struct {
	data map[any]Value
}

func NewTableImpl() *TableImpl {
	return &TableImpl{
		data: make(map[any]Value),
	}
}

func (t *TableImpl) Delete(key any) (bool, error) {
	if _, exists := t.data[key]; !exists {
		return false, errors.New("Такого ключа не существует! Удаление невозможно")
	}
	delete(t.data, key)
	return true, nil
}

func (t *TableImpl) Put(key any, value Value) (bool, error) {
	if _, exists := t.data[key]; !exists {
		return false, errors.New("Такого ключа не существует! Редактирование невозможно")
	}

	t.data[key] = value
	return true, nil
}

func (t *TableImpl) Get(key any) (Value, error) {
	value, exists := t.data[key]
	if !exists {
		return Value{}, errors.New("Ключа не существует!")
	}

	if time.Now().After(value.Ttl) {
		delete(t.data, key)
		return Value{}, errors.New("Ключ был удален")
	}

	return value, nil
}

func (t *TableImpl) Size() int {
	return len(t.data)
}
