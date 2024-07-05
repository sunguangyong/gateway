package util

import "reflect"

// 判断是否是该类型零值

func IsZeroValue(value interface{}) bool {
	rv := reflect.ValueOf(value)
	zeroValue := reflect.Zero(rv.Type()).Interface()
	return reflect.DeepEqual(value, zeroValue)
}
