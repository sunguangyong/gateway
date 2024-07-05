package convert

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"

	"github.com/duke-git/lancet/v2/convertor"
)

// ByteToBitSlice 将byte 转为 8位 []
func ByteToBitSlice(value byte) (result []byte) {
	for i := 0; i < 8; i++ {
		bit := (value >> i) & 0x01
		result = append(result, bit)
	}
	return
}

type BytesToNum func(b []byte) (result interface{})

type AnyToBytes func(method, value string) (result []byte, err error)

func ReadBool(b []byte) (result interface{}) {
	if len(b) < 1 {
		return
	}

	if b[0] > 0 {
		return true
	} else {
		return false
	}
}

func ReadUShort(b []byte) (result interface{}) {
	if len(b) < 2 {
		return
	}
	return binary.BigEndian.Uint16(b)
}

func ReadShort(b []byte) (result interface{}) {
	if len(b) < 2 {
		return
	}
	return int16(binary.BigEndian.Uint16(b))
}

func ReadUInt(b []byte) (result interface{}) {
	if len(b) < 4 {
		return
	}
	return binary.BigEndian.Uint32(b)
}

func ReadInt(b []byte) (result interface{}) {
	if len(b) < 4 {
		return
	}
	return int32(binary.BigEndian.Uint32(b))
}

func ReadULong(b []byte) (result interface{}) {
	if len(b) < 8 {
		return
	}
	result = binary.BigEndian.Uint64(b)
	return
}

func ReadLong(b []byte) (result interface{}) {
	if len(b) < 8 {
		return
	}
	result = int64(binary.BigEndian.Uint64(b))
	return
}

func ReadFloat(b []byte) (result interface{}) {
	if len(b) < 4 {
		return
	}
	bits := binary.BigEndian.Uint32(b)
	result = math.Float32frombits(bits)
	return result
}

func ReadDouble(b []byte) (result interface{}) {
	if len(b) < 8 {
		return
	}
	bits := binary.BigEndian.Uint64(b)
	result = math.Float64frombits(bits)
	return result
}

func StringToBoolBytes(method, value string) (result []byte, err error) {
	result = make([]byte, 1)
	if value == "0" {
		result[0] = 0
		return
	}

	if value == "1" {
		result[0] = 1
		return
	}

	err = NewStringToBytesErr(method, value)

	return
}

func StringToShortBytes(method, value string) (result []byte, err error) {
	result = make([]byte, 2)
	val, err := strconv.ParseInt(value, 10, 16)
	if err != nil {
		err = NewStringToBytesErr(method, value)
		return
	}
	binary.BigEndian.PutUint16(result, uint16(val))
	return
}

func StringToUShortBytes(method, value string) (result []byte, err error) {
	result = make([]byte, 2)
	val, err := strconv.ParseUint(value, 10, 16)
	if err != nil {
		err = NewStringToBytesErr(method, value)
		return
	}
	binary.BigEndian.PutUint16(result, uint16(val))
	return
}

func StringToUIntBytes(method, value string) (result []byte, err error) {
	result = make([]byte, 4)
	val, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		err = NewStringToBytesErr(method, value)
		return
	}
	binary.BigEndian.PutUint32(result, uint32(val))
	return
}

func StringToIntBytes(method, value string) (result []byte, err error) {
	result = make([]byte, 4)
	val, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		err = NewStringToBytesErr(method, value)
		return
	}
	binary.BigEndian.PutUint32(result, uint32(val))
	return
}

func StringToULongBytes(method, value string) (result []byte, err error) {
	result = make([]byte, 8)
	val, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		err = NewStringToBytesErr(method, value)
		return
	}
	binary.BigEndian.PutUint64(result, val)
	return
}

func StringToLongBytes(method, value string) (result []byte, err error) {
	result = make([]byte, 8)
	val, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		err = NewStringToBytesErr(method, value)
		return
	}
	binary.BigEndian.PutUint64(result, uint64(val))
	return
}

func StringToFloatBytes(method, value string) (result []byte, err error) {
	result = make([]byte, 8)
	val, err := strconv.ParseFloat(value, 32)
	if err != nil {
		err = NewStringToBytesErr(method, value)
		return
	}
	bits := math.Float32bits(float32(val))
	binary.BigEndian.PutUint32(result, bits)
	return
}

func StringToDoubleBytes(method, value string) (result []byte, err error) {
	result = make([]byte, 8)
	val, err := strconv.ParseFloat(value, 64)
	if err != nil {
		err = NewStringToBytesErr(method, value)
		return
	}
	bits := math.Float64bits(val)
	binary.BigEndian.PutUint64(result, bits)
	return
}

func AnyConvertToBytes(num any) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, num)
	// 如果数字是负数,需要额外添加一个字节来表示负数
	ok := isGreaterThanZero(num)
	if !ok {
		binary.Write(buf, binary.BigEndian, -1)
	}
	return buf.Bytes()
}

func isGreaterThanZero(value any) bool {
	switch reflect.TypeOf(value).Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return reflect.ValueOf(value).Int() > 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return reflect.ValueOf(value).Uint() > 0
	case reflect.Float32, reflect.Float64:
		return reflect.ValueOf(value).Float() > 0
	default:
		return false
	}
}

func CommonToBytes(quantity uint16, value interface{}) (bytesData []byte, err error) {
	bytesData, err = convertor.ToBytes(value)
	if err != nil {
		return
	}

	if len(bytesData) != 8 {
		return
	}

	return bytesData[8-quantity:], nil
}

func NewStringToBytesErr(method, value string) error {
	return errors.New(fmt.Sprintf("mehtod = %s 类型错误 value = %s", method, value))
}
