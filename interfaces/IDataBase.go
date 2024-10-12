package interfaces

type IDataBase interface {
	add(keyDB interface{}, table ITable) (bool, error)
	get(keyDB interface{}) (ITable, error)
	remove(keyDB interface{}) (bool, error)
	put(keyDB interface{}, table ITable) (bool, error)
}
