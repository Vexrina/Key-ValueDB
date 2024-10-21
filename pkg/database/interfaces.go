package database

import "time"

type DataBase interface {
	Create(keyDB any, table Table) (bool, error)
	Select(keyDB any) (Table, error)
	Delete(keyDB any) (bool, error)
	Rename(keyOld, keyNew any) (bool, error)
	SelectAll() (map[any]Table, error)
}

type Table interface {
	Delete(keyTable any) (bool, error)
	Insert(keyTable any, value Value) (bool, error)
	Get(keyTable any) (Value, error)
	Update(keyTable any, value Value) (bool, error)
	Size() int
	parseTime(format string) (time.Time, error)
}
