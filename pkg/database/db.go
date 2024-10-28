package database

import (
	"errors"
)

type DataBaseImpl struct {
	dataBase map[any]TableImpl
}

func NewDataBaseImpl() *DataBaseImpl {
	return &DataBaseImpl{
		dataBase: make(map[any]TableImpl),
	}
}

func (db *DataBaseImpl) Create(keyDB any, table TableImpl) (bool, error) {
	if _, exists := db.dataBase[keyDB]; exists {
		return false, errors.New("Таблица с таким ключом уже существует")
	}

	db.dataBase[keyDB] = table
	return true, nil
}

func (db *DataBaseImpl) Select(keyDB any) (TableImpl, error) {
	table, exists := db.dataBase[keyDB]
	if !exists {
		return TableImpl{}, errors.New("Таблица не найдена")
	}

	return table, nil
}

func (db *DataBaseImpl) Delete(keyDB any) (bool, error) {
	if _, exists := db.dataBase[keyDB]; !exists {
		return false, errors.New("Такой таблицы не существует")
	}

	delete(db.dataBase, keyDB)
	return true, nil
}

func (db *DataBaseImpl) Rename(keyOld, keyNew any) (bool, error) {
	if _, exists := db.dataBase[keyOld]; !exists {
		return false, errors.New("Такой таблицы не существует")
	}

	if keyOld == keyNew {
		return false, errors.New("Старый и новый ключи совпадают")
	}

	if _, exists := db.dataBase[keyNew]; exists {
		return false, errors.New("Новый ключ уже существует")
	}

	db.dataBase[keyNew] = db.dataBase[keyOld]
	delete(db.dataBase, keyOld)

	return true, nil
}

func (db *DataBaseImpl) SelectAll() (map[any]TableImpl, error) {
	if len(db.dataBase) == 0 {
		return nil, errors.New("нет данных")
	}

	return db.dataBase, nil
}
