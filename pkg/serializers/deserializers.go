package serializers

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

func (s *SerializerImpl) Deserialize(data []byte, target interface{}) error {
	v := reflect.ValueOf(target)

	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("target must be pointer, but target is: %v", v.Kind())
	}
	v = v.Elem()
	switch v.Kind() {
	case reflect.Uint8:
		return s.deserializeUint8(data, v)
	case reflect.Int:
		return s.deserializeInt(data, v)
	case reflect.String:
		return s.deserializeString(data, v)
	case reflect.Float64:
		return s.deserializeFloat64(data, v)
	case reflect.Bool:
		return s.deserializeBool(data, v)
	case reflect.Slice, reflect.Array:
		return s.deserializeArray(data, v)
	case reflect.Map:
		return s.deserializeMap(data, v)
	case reflect.Struct:
		return s.deserializeStruct(data, v)
	case reflect.Ptr:
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		return s.Deserialize(data, v.Addr().Interface())
	default:
		return fmt.Errorf("s.deserialize got unsupported types: %v", v.Kind())
	}
}

func (s *SerializerImpl) deserializeUint8(data []byte, v reflect.Value) error {
	if len(data) < 1 {
		return fmt.Errorf("deserializeUint8: data too short")
	}

	var result uint8
	buf := bytes.NewReader(data)

	if err := binary.Read(buf, binary.LittleEndian, &result); err != nil {
		return fmt.Errorf("deserializeUint8 got error: %w", err)
	}

	if v.CanSet() {
		v.SetUint(uint64(result))
	} else {
		return fmt.Errorf("deserializeUint8: cannot set value")
	}

	return nil
}

func zigZagDecode(n uint64) int64 {
	return int64((n >> 1) ^ -(n & 1))
}

func (s *SerializerImpl) deserializeInt(data []byte, v reflect.Value) error {
	var result uint64
	var shift uint
	var byteValue byte

	for i := 0; i < len(data); i++ {
		byteValue = data[i]
		result |= uint64(byteValue&0x7F) << shift
		if byteValue&0x80 == 0 {
			break
		}
		shift += 7
	}

	decodedValue := zigZagDecode(result)
	v.SetInt(decodedValue)
	return nil
}

func (s *SerializerImpl) deserializeString(data []byte, v reflect.Value) error {
	var length int64
	lengthValue := reflect.New(reflect.TypeOf(length)).Elem()

	err := s.deserializeInt(data, lengthValue)
	if err != nil {
		return fmt.Errorf("deserializeString getting length got error: %w", err)
	}
	length = lengthValue.Int()

	lengthBytes, err := s.serializeInt(length)
	if err != nil {
		return fmt.Errorf("deserializeString calculating length got error: %w", err)
	}
	readBytes := len(lengthBytes)

	if length == 0 {
		v.SetString("")
		return nil
	}

	if int64(len(data)-readBytes) < length {
		return fmt.Errorf("deserializeString: insufficient data for reading string")
	}

	result := data[readBytes : readBytes+int(length)]
	v.SetString(string(result))
	return nil
}

func (s *SerializerImpl) deserializeFloat64(data []byte, v reflect.Value) error {
	var result float64
	buf := bytes.NewReader(data)
	if err := binary.Read(buf, binary.LittleEndian, &result); err != nil {
		return fmt.Errorf("deserializeFloat64 got error: %w", err)
	}
	v.SetFloat(result)
	return nil
}

func (s *SerializerImpl) deserializeBool(data []byte, v reflect.Value) error {
	var boolByte byte
	buf := bytes.NewReader(data)
	if err := binary.Read(buf, binary.LittleEndian, &boolByte); err != nil {
		return fmt.Errorf("deserializeBool got error: %w", err)
	}
	v.SetBool(boolByte == 1)
	return nil
}

func (s *SerializerImpl) deserializeArray(data []byte, v reflect.Value) error {
	// Декодируем длину массива с использованием deserializeInt
	var length int64
	lengthValue := reflect.New(reflect.TypeOf(length)).Elem()

	err := s.deserializeInt(data, lengthValue)
	if err != nil {
		return fmt.Errorf("deserializeArray_length got error: %w", err)
	}
	length = lengthValue.Int()

	// Создаем новый слайс с декодированной длиной
	v.Set(reflect.MakeSlice(v.Type(), int(length), int(length)))

	lengthBytes, _ := s.serializeInt(length)
	readBytes := len(lengthBytes) // Число байт, использованных для длины массива

	for i := 0; i < int(length); i++ {
		// Декодируем длину каждого элемента с использованием deserializeInt
		var elementLength int64
		elementLengthValue := reflect.New(reflect.TypeOf(elementLength)).Elem()

		err := s.deserializeInt(data[readBytes:], elementLengthValue)
		if err != nil {
			return fmt.Errorf("deserializeArray_elem #%d length reading got error: %w", i, err)
		}
		elementLength = elementLengthValue.Int()
		elementLengthBytes, _ := s.serializeInt(elementLength)

		// Проверяем, достаточно ли данных для чтения элемента
		if readBytes+len(elementLengthBytes)+int(elementLength) > len(data) {
			return fmt.Errorf("deserializeArray_elem #%d insufficient data for reading element", i)
		}

		// Читаем элемент
		elementBytes := data[readBytes+len(elementLengthBytes) : readBytes+len(elementLengthBytes)+int(elementLength)]

		// Обновляем readBytes
		readBytes += len(elementLengthBytes) + int(elementLength)

		// Создаем указатель на тип элемента массива
		elementPtr := reflect.New(v.Type().Elem()).Interface()

		// Десериализуем элемент
		if err := s.Deserialize(elementBytes, elementPtr); err != nil {
			return fmt.Errorf("deserializeArray_elem #%d deserializing got error: %w", i, err)
		}

		// Устанавливаем элемент в массив
		v.Index(i).Set(reflect.ValueOf(elementPtr).Elem())
	}

	return nil
}

func (s *SerializerImpl) deserializeStruct(data []byte, v reflect.Value) error {
	buf := bytes.NewReader(data)

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		var fieldLength int64
		if err := binary.Read(buf, binary.LittleEndian, &fieldLength); err != nil {
			return fmt.Errorf("deserializeStruct_field #%d reading length got error: %w", i, err)
		}

		if fieldLength == -1 {
			// Устанавливаем поле как nil для указателей
			if field.Kind() == reflect.Ptr {
				field.Set(reflect.Zero(field.Type()))
			}
			continue
		}

		if fieldLength < 0 {
			return fmt.Errorf("deserializeStruct_field #%d invalid field length: %d", i, fieldLength)
		}

		fieldBytes := make([]byte, fieldLength)
		if _, err := buf.Read(fieldBytes); err != nil {
			return fmt.Errorf("deserializeStruct_field #%d reading bytes got error: %w", i, err)
		}

		// Десериализуем поле
		fieldPtr := reflect.New(field.Type()).Interface()
		if err := s.Deserialize(fieldBytes, fieldPtr); err != nil {
			return fmt.Errorf("deserializeStruct_field #%d deserializing got error: %w", i, err)
		}

		// Устанавливаем значение поля
		field.Set(reflect.ValueOf(fieldPtr).Elem())
	}

	return nil
}
func (s *SerializerImpl) deserializeMap(data []byte, v reflect.Value) error {
	buf := bytes.NewReader(data)

	var length int64
	if err := binary.Read(buf, binary.LittleEndian, &length); err != nil {
		return fmt.Errorf("deserializeMap_length got error: %w", err)
	}

	v.Set(reflect.MakeMap(v.Type()))

	for i := 0; i < int(length); i++ {
		keyType, err := s.deserializeType(buf)
		if err != nil {
			return fmt.Errorf("deserializeMap_key #%d deserializing type got error: %w", i, err)
		}

		keyValue := reflect.New(keyType).Elem()

		var keyLength int64
		if err := binary.Read(buf, binary.LittleEndian, &keyLength); err != nil {
			return fmt.Errorf("deserializeMap_key #%d reading length got error: %w", i, err)
		}

		keyBytes := make([]byte, keyLength)
		if _, err := buf.Read(keyBytes); err != nil {
			return fmt.Errorf("deserializeMap_key #%d reading bytes got error: %w", i, err)
		}

		if err := s.Deserialize(keyBytes, keyValue.Addr().Interface()); err != nil {
			return fmt.Errorf("deserializeMap_key #%d deserializing got error: %w", i, err)
		}

		valueType, err := s.deserializeType(buf)
		if err != nil {
			return fmt.Errorf("deserializeMap_value #%d deserializing type got error: %w", i, err)
		}

		valueValue := reflect.New(valueType).Elem()

		var valueLength int64
		if err := binary.Read(buf, binary.LittleEndian, &valueLength); err != nil {
			return fmt.Errorf("deserializeMap_value #%d reading length got error: %w", i, err)
		}

		valueBytes := make([]byte, valueLength)
		if _, err := buf.Read(valueBytes); err != nil {
			return fmt.Errorf("deserializeMap_value #%d reading bytes got error: %w", i, err)
		}

		if err := s.Deserialize(valueBytes, valueValue.Addr().Interface()); err != nil {
			return fmt.Errorf("deserializeMap_value #%d deserializing got error: %w", i, err)
		}

		v.SetMapIndex(keyValue, valueValue)
	}

	return nil
}

func (s *SerializerImpl) deserializeType(buf *bytes.Reader) (reflect.Type, error) {
	typeID, err := buf.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("deserializeType got error: %w", err)
	}

	kind := reflect.Kind(typeID)
	return s.kindToType(kind), nil
}

func (s *SerializerImpl) kindToType(kind reflect.Kind) reflect.Type {
	switch kind {
	case reflect.Int:
		return reflect.TypeOf(int(0))
	case reflect.Int64:
		return reflect.TypeOf(int64(0))
	case reflect.Float64:
		return reflect.TypeOf(float64(0))
	case reflect.String:
		return reflect.TypeOf("")
	case reflect.Interface:
		return reflect.TypeOf((*interface{})(nil)).Elem()
	default:
		panic(fmt.Sprintf("unsupported kind: %v", kind))
	}
}
