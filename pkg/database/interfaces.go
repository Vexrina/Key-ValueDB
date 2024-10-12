package database

type IDataBase interface {
	add(keyDB any, table Table) (bool, error)
	get(keyDB any) (Table, error)
	remove(keyDB any) (bool, error)
	put(keyDB any, table Table) (bool, error)
}

type Table interface {
	add(key any, value Value) (bool, error)
	delete(key any) (bool, error)
	put(key any, value Value) (bool, error)
	get(key any) (Value, error)
	size() int
}
