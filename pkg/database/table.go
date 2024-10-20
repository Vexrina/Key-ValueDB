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
	dataTable map[any]Value
}

func NewTableImpl() *TableImpl {
	return &TableImpl{
		dataTable: make(map[any]Value),
	}
}

func (t *TableImpl) Delete(keyTable any) (bool, error) {
	if _, exists := t.dataTable[keyTable]; !exists {
		return false, errors.New("Такого ключа не существует! Удаление невозможно")
	}
	delete(t.dataTable, keyTable)
	return true, nil
}

func (t *TableImpl) Insert(keyTable any, value Value) (bool, error) {
	if _, exists := t.dataTable[keyTable]; exists {
		return false, errors.New("Такой ключ существует! Добавление невозможно")
	}

	t.dataTable[keyTable] = value
	return true, nil
}
func (t *TableImpl) Update(keyTable any, value Value) (bool, error) {
	if _, exists := t.dataTable[keyTable]; !exists {
		return false, errors.New("Такого ключа не существует! Редактирование невозможно")
	}

	t.dataTable[keyTable] = value
	return true, nil
}

func (t *TableImpl) Get(keyTable any) (Value, error) {
	value, exists := t.dataTable[keyTable]
	if !exists {
		return Value{}, errors.New("Ключа не существует")
	}

	if time.Now().After(value.Ttl) {
		delete(t.dataTable, keyTable)
		return Value{}, errors.New("Ключ был удален")
	}

	return value, nil
}

func (t *TableImpl) Size() int {
	return len(t.dataTable)
}

func (t *TableImpl) ParseTime(dateStr string) (time.Time, error) {
	parsedTime, err := time.Parse("02.01.2006", dateStr)
	if err != nil {
		return time.Time{}, errors.New("Невозможно распарсить дату!")
	}
	return parsedTime, nil
}
