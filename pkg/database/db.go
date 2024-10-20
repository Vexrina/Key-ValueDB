package database

import (
	"errors"
)

type DataBaseImpl struct {
	dataBase map[any]Table
}

func NewDataBaseImpl() *DataBaseImpl {
	return &DataBaseImpl{
		dataBase: make(map[any]Table),
	}
}

func (db *DataBaseImpl) Set(keyDB any, table Table) (bool, error) {

	if _, exists := db.dataBase[keyDB]; exists {
		return false, errors.New("Таблица с таким ключом уже существует")
	}

	db.dataBase[keyDB] = table
	return true, nil
}

func (db *DataBaseImpl) Get(keyDB any) (Table, error) {

	table, exists := db.dataBase[keyDB]
	if !exists {
		return nil, errors.New("Таблица не найдена")
	}

	return table, nil
}

func (db *DataBaseImpl) Remove(keyDB any) (bool, error) {
	if _, exists := db.dataBase[keyDB]; !exists {
		return false, errors.New("Такой таблицы не существует")
	}

	delete(db.dataBase, keyDB)
	return true, nil
}

func (db *DataBaseImpl) Put(keyDB any, table Table) (bool, error) {
	if _, exists := db.dataBase[keyDB]; !exists {
		return false, errors.New("Такой таблицы не существует")
	}

	db.dataBase[keyDB] = table
	return true, nil
}
