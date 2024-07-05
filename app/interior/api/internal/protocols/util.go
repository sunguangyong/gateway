package protocols

import (
	"errors"
	"fmt"

	"xunjikeji.com.cn/gateway/common/constant"
	"xunjikeji.com.cn/gateway/common/xerr"
)

func PlcWriteBytes(method string, value string) (result []byte, err error) {
	funcValue, ok := constant.WriteValueFunc[method]
	if !ok {
		err = errors.New(fmt.Sprintf("找不到 method = %s 方法", method))
	}
	result, err = funcValue(method, value)
	return
}

func PlcReadeBytes(method string, value []byte) (result interface{}, err error) {
	// TODO 非法值问题 float NAN 类型
	funcType, ok := constant.ReadValueFunc[method]

	if !ok {
		err = xerr.NotFoundMethodErrErr
		return
	}
	result = funcType(value)
	return
}
