package convert

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

// CopyProperties 将src的字段值赋值到和dst相同的字段
func CopyProperties(dst, src any) error {
	dstType, srcType := reflect.TypeOf(dst), reflect.TypeOf(src)

	if dstType.Kind() != reflect.Ptr || dstType.Elem().Kind() != reflect.Struct {
		return errors.New("CopyProperties: parameter dst should be struct pointer")
	}

	if srcType.Kind() == reflect.Ptr {
		srcType = srcType.Elem()
	}
	if srcType.Kind() != reflect.Struct {
		return errors.New("CopyProperties: parameter src should be a struct or struct pointer")
	}

	bytes, err := json.Marshal(src)
	if err != nil {
		return fmt.Errorf("CopyProperties: unable to marshal src: %s", err)
	}
	err = json.Unmarshal(bytes, dst)
	if err != nil {
		return fmt.Errorf("CopyProperties: unable to unmarshal into dst: %s", err)
	}

	return nil
}
