package gotag

import (
	"log"
	"reflect"
)

const (
	ERR_SUCESS        = 0
	ERR_FAILED        = 1
	ERR_INVALID_INPUT = 2
	ERR_NULL_ACCESS   = 3
)

const (
	TAG_VALUE_TYPE_BOOLEAN   = 0
	TAG_VALUE_TYPE_INT8      = 1
	TAG_VALUE_TYPE_INT16     = 2
	TAG_VALUE_TYPE_INT32     = 3
	TAG_VALUE_TYPE_INT64     = 4
	TAG_VALUE_TYPE_INT       = 5
	TAG_VALUE_TYPE_UINT8     = 6
	TAG_VALUE_TYPE_UINT16    = 7
	TAG_VALUE_TYPE_UINT32    = 8
	TAG_VALUE_TYPE_UINT64    = 9
	TAG_VALUE_TYPE_UINT      = 10
	TAG_VALUE_TYPE_FLOAT     = 11
	TAG_VALUE_TYPE_DOUBLE    = 12
	TAG_VALUE_TYPE_STRING    = 13
	TAG_VALUE_TYPE_BYTEARRAY = 14
)

type Value struct {
	i int64
	u uint64
	f float32
	d float64
	s string
	b []byte
}

func (m *Value) GetFloat() float32 {
	if m != nil {
		return m.f
	}
	return 0
}

func (m *Value) GetDouble() float64 {
	if m != nil {
		return m.d
	}
	return 0
}

func (m *Value) GetInt() int64 {
	if m != nil {
		return m.i
	}
	return 0
}

func (m *Value) GetUint() uint64 {
	if m != nil {
		return m.u
	}
	return 0
}

func (m *Value) GetStr() string {
	if m != nil {
		return m.s
	}
	return ""
}

func (m *Value) GetBytes() []byte {
	if m != nil {
		return m.b
	}
	return nil
}

func NewValue(value interface{}) *Value {
	switch reflect.TypeOf(value).Kind() {
	case reflect.Int:
		return &Value{i: int64(value.(int))}
	case reflect.Int16:
		return &Value{i: int64(value.(int16))}
	case reflect.Int32:
		return &Value{i: int64(value.(int32))}
	case reflect.Int64:
		return &Value{i: value.(int64)}
	case reflect.Uint:
		return &Value{u: uint64(value.(uint))}
	case reflect.Uint16:
		return &Value{u: uint64(value.(uint16))}
	case reflect.Uint32:
		return &Value{u: uint64(value.(uint32))}
	case reflect.Uint64:
		return &Value{u: value.(uint64)}
	case reflect.Float32:
		return &Value{f: value.(float32)}
	case reflect.Float64:
		return &Value{d: value.(float64)}
	case reflect.String:
		return &Value{s: value.(string)}
	case reflect.Array, reflect.Slice:
		return &Value{b: value.([]byte)}
	default:
		log.Printf("kind: %v\n", reflect.TypeOf(value).Kind())
	}
	return &Value{}
}

type Tag struct {
	SourceName string
	TagName    string
	Val        *Value
	ValType    int32
	Ts         uint64
	Unit       string
}
