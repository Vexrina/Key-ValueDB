package database

type DataBase interface {
	Set(keyDB any, table Table) (bool, error)
	Get(keyDB any) (Table, error)
	Remove(keyDB any) (bool, error)
	Put(keyDB any, table Table) (bool, error)
}

type Table interface {
	Delete(key any) (bool, error)
	Put(key any, value Value) (bool, error)
	Get(key any) (Value, error)
	Size() int
}
