package serializers

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

func (s *SerializerImpl) Serialize(value interface{}) ([]byte, error) {
	v := reflect.ValueOf(value)

	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil, fmt.Errorf("s.serialize got nil pointer")
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Uint8:
		return s.serializeUint8(uint8(v.Uint()))
	case reflect.Int:
		return s.serializeInt(v.Int())
	case reflect.String:
		return s.serializeString(v.String())
	case reflect.Float64, reflect.Float32:
		return s.serializeFloat(v.Float())
	case reflect.Bool:
		return s.serializeBool(v.Bool())
	case reflect.Slice, reflect.Array:
		return s.serializeArray(v)
	case reflect.Map:
		return s.serializeMap(v)
	case reflect.Struct:
		return s.serializeStruct(v)
	default:
		return nil, fmt.Errorf("s.serialize got unsupported types: %v", v.Kind())
	}
}

func (s *SerializerImpl) serializeUint8(value uint8) ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.LittleEndian, value); err != nil {
		return nil, fmt.Errorf("serializeUint8 got error: %w", err)
	}

	return buf.Bytes(), nil
}

func zigZagEncode(i int64) uint64 {
	return uint64((i << 1) ^ (i >> 63))
}

func (s *SerializerImpl) serializeInt(i int64) ([]byte, error) {
	zigzagEncoded := zigZagEncode(i)

	buf := new(bytes.Buffer)
	for zigzagEncoded >= 0x80 {
		buf.WriteByte(byte(zigzagEncoded&0x7F | 0x80))
		zigzagEncoded >>= 7
	}
	buf.WriteByte(byte(zigzagEncoded))
	return buf.Bytes(), nil
}

func (s *SerializerImpl) serializeString(str string) ([]byte, error) {
	buf := new(bytes.Buffer)
	length := int64(len(str))
	lengthBytes, err := s.serializeInt(length)
	if err != nil {
		return nil, fmt.Errorf("serializeString got error: %w", err)
	}

	buf.Write(lengthBytes)
	buf.WriteString(str)
	return buf.Bytes(), nil
}

func (s *SerializerImpl) serializeFloat(f float64) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, f); err != nil {
		return nil, fmt.Errorf("serializeFloat got error: %w", err)
	}
	return buf.Bytes(), nil
}

func (s *SerializerImpl) serializeBool(b bool) ([]byte, error) {
	buf := new(bytes.Buffer)
	var boolByte byte
	if b {
		boolByte = 1
	} else {
		boolByte = 0
	}
	if err := binary.Write(buf, binary.LittleEndian, boolByte); err != nil {
		return nil, fmt.Errorf("serializeBool got error: %w", err)
	}
	return buf.Bytes(), nil
}

func (s *SerializerImpl) serializeArray(arr reflect.Value) ([]byte, error) {
	buf := new(bytes.Buffer)

	length := int64(arr.Len())
	lengthBytes, err := s.serializeInt(length)
	if err != nil {
		return nil, fmt.Errorf("serializeString got error: %w", err)
	}
	buf.Write(lengthBytes)

	for idx := 0; idx < arr.Len(); idx++ {
		element := arr.Index(idx).Interface()

		elementBytes, err := s.Serialize(element)
		if err != nil {
			return nil, fmt.Errorf("serializeArray_elem #%d serializing got error: %w", idx, err)
		}

		elementLength := int64(len(elementBytes))
		lengthBytes, err := s.serializeInt(elementLength)
		if err != nil {
			return nil, fmt.Errorf("serializeString got error: %w", err)
		}
		buf.Write(lengthBytes)

		if _, err := buf.Write(elementBytes); err != nil {
			return nil, fmt.Errorf("serializeArray_elem #%d writing bytes got error: %w", idx, err)
		}
	}

	return buf.Bytes(), nil
}

func (s *SerializerImpl) serializeStruct(v reflect.Value) ([]byte, error) {
	buf := new(bytes.Buffer)

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		if field.Kind() == reflect.Ptr {
			if field.IsNil() {
				// Для nil указателя записываем отрицательную длину или другое специальное значение
				if err := binary.Write(buf, binary.LittleEndian, int64(-1)); err != nil {
					return nil, fmt.Errorf("serializeStruct_field #%d writing nil length got error: %w", i, err)
				}
				continue
			}
		}

		// Сериализуем каждое поле
		fieldBytes, err := s.Serialize(field.Interface())
		if err != nil {
			return nil, fmt.Errorf("serializeStruct_field #%d serializing got error: %w", i, err)
		}

		fieldLength := int64(len(fieldBytes))
		if err := binary.Write(buf, binary.LittleEndian, fieldLength); err != nil {
			return nil, fmt.Errorf("serializeStruct_field #%d writing length got error: %w", i, err)
		}

		if _, err := buf.Write(fieldBytes); err != nil {
			return nil, fmt.Errorf("serializeStruct_field #%d writing bytes got error: %w", i, err)
		}
	}

	return buf.Bytes(), nil
}

func (s *SerializerImpl) serializeMap(m reflect.Value) ([]byte, error) {
	buf := new(bytes.Buffer)

	length := int64(m.Len())
	if err := binary.Write(buf, binary.LittleEndian, length); err != nil {
		return nil, fmt.Errorf("serializeMap_length got error: %w", err)
	}

	for idx, key := range m.MapKeys() {
		keyType := reflect.TypeOf(key.Interface())
		if err := s.serializeType(buf, keyType); err != nil {
			return nil, fmt.Errorf("serializeMap_key #%d serializing type got error: %w", idx, err)
		}

		keyBytes, err := s.Serialize(key.Interface())
		if err != nil {
			return nil, fmt.Errorf("serializeMap_key #%d serializing got error: %w", idx, err)
		}

		keyLength := int64(len(keyBytes))
		if err := binary.Write(buf, binary.LittleEndian, keyLength); err != nil {
			return nil, fmt.Errorf("serializeMap_key #%d writing length got error: %w", idx, err)
		}
		if _, err := buf.Write(keyBytes); err != nil {
			return nil, fmt.Errorf("serializeMap_key #%d writing bytes got error: %w", idx, err)
		}

		value := m.MapIndex(key).Interface()
		valueType := reflect.TypeOf(value)
		if err := s.serializeType(buf, valueType); err != nil {
			return nil, fmt.Errorf("serializeMap_value #%d serializing type got error: %w", idx, err)
		}

		valueBytes, err := s.Serialize(value)
		if err != nil {
			return nil, fmt.Errorf("serializeMap_value #%d serializing got error: %w", idx, err)
		}

		valueLength := int64(len(valueBytes))
		if err := binary.Write(buf, binary.LittleEndian, valueLength); err != nil {
			return nil, fmt.Errorf("serializeMap_value #%d writing length got error: %w", idx, err)
		}
		if _, err := buf.Write(valueBytes); err != nil {
			return nil, fmt.Errorf("serializeMap_value #%d writing bytes got error: %w", idx, err)
		}
	}

	return buf.Bytes(), nil
}

func (s *SerializerImpl) serializeType(buf *bytes.Buffer, typ reflect.Type) error {
	typeID := uint8(typ.Kind())
	if err := buf.WriteByte(typeID); err != nil {
		return fmt.Errorf("serializeType got error: %w", err)
	}
	return nil
}
