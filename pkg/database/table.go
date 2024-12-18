package database

import (
	"errors"
	"fmt"
	"time"
)

type Value struct {
	Val    any       `json:"val"`
	Ttl    time.Time `json:"ttl"`
	BadTtl string
}

type TableImpl struct {
	DataTable map[any]Value `json:"table"`
}

func NewTableImpl() *TableImpl {
	return &TableImpl{
		DataTable: make(map[any]Value),
	}
}

func NewValue(val any, dateStr string) (Value, error) {
	ttl, err := parseTime(dateStr)
	if err != nil {
		return Value{}, fmt.Errorf("parseTime error: %w", err)
	}
	return Value{
		Val:    val,
		Ttl:    ttl,
		BadTtl: dateStr,
	}, nil
}

func (t *TableImpl) Delete(keyTable any) (bool, error) {
	if _, exists := t.DataTable[keyTable]; !exists {
		return false, errors.New("Такого ключа не существует! Удаление невозможно")
	}
	delete(t.DataTable, keyTable)
	return true, nil
}

func (t *TableImpl) Insert(keyTable any, value Value) (bool, error) {
	t.checker()
	if _, exists := t.DataTable[keyTable]; exists {
		return false, errors.New("Такой ключ существует! Добавление невозможно")
	}

	t.DataTable[keyTable] = value
	return true, nil
}

func (t *TableImpl) Update(keyTable any, value Value) (bool, error) {
	t.checker()
	if _, exists := t.DataTable[keyTable]; !exists {
		return false, errors.New("Такого ключа не существует! Редактирование невозможно")
	}

	t.DataTable[keyTable] = value
	return true, nil
}

func (t *TableImpl) Get(keyTable any) (Value, error) {
	t.checker()
	value, exists := t.DataTable[keyTable]
	if !exists {
		return Value{}, errors.New("Ключа не существует")
	}

	if time.Now().After(value.Ttl) {
		delete(t.DataTable, keyTable)
		return Value{}, errors.New("Ключ был удален")
	}

	return value, nil
}

func (t *TableImpl) Size() int {
	return len(t.DataTable)
}

func parseTime(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Now().Add(5 * time.Minute).Truncate(time.Second), nil
	}
	parsedTime, err := time.Parse("02.01.2006T15:04:05", dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("time.Parse: %w", err)
	}
	return parsedTime, nil
}

func (t *TableImpl) checker() {
	if t.DataTable == nil {
		t.DataTable = make(map[any]Value)
	}
}
