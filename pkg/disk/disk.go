package disk

import (
	db "BD/pkg/database"
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"sort"
	"strconv"
)

type DiskImpl struct {
	DataBases db.DataBaseImpl
	Tables    db.TableImpl
}

type keyValueTypeKey struct {
	Key     string
	Value   db.TableImpl
	KeyType reflect.Type
}

func (d *DiskImpl) writeData() (bool, error) {
	data, err := d.getSortedMassiveByKey()
	if err != nil {
		return false, errors.New("failed to get sorted data: " + err.Error())
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return false, errors.New("failed to marshal data: " + err.Error())
	}

	err = os.WriteFile("data.json", jsonData, 0644)
	if err != nil {
		return false, errors.New("failed to write data: " + err.Error())
	}

	return true, nil
}

func (d *DiskImpl) readData() (bool, error) {
	return true, nil
}

func (d *DiskImpl) getSortedMassiveByKey() ([]keyValueTypeKey, error) {
	allData, err := d.DataBases.SelectAll()
	if err != nil {
		return nil, errors.New("unable to sort keys: " + err.Error())
	}

	var kvtSlice []keyValueTypeKey

	for k, v := range allData {
		keyStr, typekey, err := convertKeyToString(k)
		if err != nil {
			return nil, err
		}
		kvtSlice = append(kvtSlice, keyValueTypeKey{Key: keyStr, Value: v, KeyType: typekey})
	}

	sort.Slice(kvtSlice, func(i, j int) bool {
		return kvtSlice[i].Key < kvtSlice[j].Key
	})

	return kvtSlice, nil
}

func convertKeyToString(key interface{}) (string, reflect.Type, error) {
	switch key := key.(type) {
	case int:
		return strconv.Itoa(key), reflect.TypeOf(key), nil
	case string:
		return key, reflect.TypeOf(key), nil
	default:
		return "", nil, errors.New("unsupported key type")
	}
}
