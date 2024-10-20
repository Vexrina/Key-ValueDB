package serializers

import (
	"bytes"
	"reflect"
)

type Serializer interface {
	// Serialize позволяет сериализировать любое значение в набор байтов
	Serialize(value any) ([]byte, error)
	// serializeUint8 записывает в буффер Uint8
	serializeUint8(value uint8) ([]byte, error)
	// serializeInt записывает в буфер int64 используя VarInt и zigZag кодирование
	serializeInt(value int64) ([]byte, error)
	// serializeString записывает в буфер string, длина кодируется VarInt'ом
	serializeString(value string) ([]byte, error)
	// serializeFloat записывает в буфер float64
	serializeFloat(value float64) ([]byte, error)
	// serializeBool записывает в буфер бул (байт либо 1, либо 0)
	serializeBool(value bool) ([]byte, error)
	// serializeArray записывает в буфер массив, длина кодируется VarInt'ом
	serializeArray(value reflect.Value) ([]byte, error)
	// serializeStruct записывает в буфер структуры
	serializeStruct(value reflect.Value) ([]byte, error)
	// serializeMap записывает в буфер мапки. Перед ключами и значениями записываются их типы, используя serializeType
	serializeMap(value reflect.Value) ([]byte, error)
	// serializeType записывает в буфер тип значения, полезно для map[any]any
	serializeType(buf *bytes.Buffer, typ reflect.Type)

	// Derialize позволяет десериализировать набор байтов в нужный тип
	Deserialize(data []byte, target any) error
	// deserializeUint8 читает байтовый массив и записывает в Uint8
	deserializeUint8(data []byte, v reflect.Value) error
	// deserializeInt читает байтовый массив и записывает в Int используя VarInt и zigZag декодирование
	deserializeInt(data []byte, v reflect.Value) error
	// deserializeString читает байтовый массив и записывает в String. Длина декодируется с помощью VarInt и zigZag декодирование
	deserializeString(data []byte, v reflect.Value) error
	// deserializeFloat читает байтовый массив и записывает в float64
	deserializeFloat(data []byte, v reflect.Value) error
	// deserializeBool читает байтовый массив и записывает в bool
	deserializeBool(data []byte, v reflect.Value) error
	// deserializeArray читает байтовый массив и записывает в слайс нужного типа. Читает длину самого массива и длину каждого элемента с помощью VarInt и zigZag декодирования
	deserializeArray(data []byte, v reflect.Value) error
	// deserializeStruct читает байтовый массив и записывает в нужную структуру, заранее опреденную. Обходит структуру как граф, без циклов
	deserializeStruct(data []byte, v reflect.Value) error
	// deserializeMap читает байтовый массив и записывает в необходимую мапку, поддерживая map[any]any
	deserializeMap(data []byte, v reflect.Value) error
	// deserializeType читает буфер и отвечает, в какой тип необходимо кастануть байты
	deserializeType(buf *bytes.Reader) (reflect.Type, error)

	// zigZagEncode переводит из знаковых чисел в беззнаковые. Нужно для VarInt с отрицательными числами. Объяснение: https://www.programmerall.com/article/9705468100/
	zigZagEncode(i int64) uint64
	// zigZagDecode переводит из безнаковые в знаковые. Обратная операция zigZagEncode
	zigZagDecode(n uint64) int64
	// kindToType используя switch возвращает тип, в который нужно кастануть байты. Используется в deserializeType
	kindToType(kind reflect.Kind) reflect.Type
}

type SerializerImpl struct{}
