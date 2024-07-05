package xerr

import (
	"fmt"

	"github.com/pkg/errors"
)

/*
常用通用固定错误
*/

var (
	ProtocolExistsError       = NewErr(ProtocolExistsErrorCode, errors.New("此设备已经配置协议"))
	NotSupportedProtocolError = NewErr(NotSupportedProtocolErrorCode, errors.New("暂未支持此协议"))
	NotFoundNidErr            = NewErr(NotFoundNidErrCode, errors.New("找不到此协议"))
	NotFoundConfigErr         = NewErr(NotFoundConfigErrCode, errors.New("找不到此点表"))
	NotFoundNeDeviceErr       = NewErr(NotFoundNeDeviceErrCode, errors.New("找不到此设备"))
	MacExistsError            = NewErr(MacExistsErrCode, errors.New("mac地址 已存在"))
	NotFoundFeEdgeDeviceErr   = NewErr(NotFoundFeEdgeDeviceErrCode, errors.New("找不到此网关"))
	NotFoundFeEdgeDeviceIpErr = NewErr(NotFoundFeEdgeDeviceIpErrCode, errors.New("找不到网关Ip"))

	NotFoundMethodErrErr = NewErr(NotFoundMethodErrCode, errors.New("找不到method方法"))
)

type CodeError struct {
	errCode uint32
	err     error
}

func ConvertAccessOptionsParamsErr(err error) *CodeError {
	if err == nil {
		return nil
	}
	return NewErr(RequestParamErrorCode, errors.New(err.Error()))
}

// GetErrCode 返回给前端的错误码
func (e *CodeError) GetErrCode() uint32 {
	return e.errCode
}

// GetErrMsg 返回给前端显示端错误信息
func (e *CodeError) GetErrMsg() string {
	return e.err.Error()
}

func (e *CodeError) Error() string {

	return fmt.Sprintf("errCode:%d，err:%+v", e.errCode, e.err)
	//return fmt.Sprintf("ErrCode:%d，ErrMsg:%s", e.errCode, e.errMsg)
}

func NewErr(code uint32, err error) *CodeError {
	//err = errors.New(err.Error())
	return &CodeError{errCode: code, err: err}
}

func ErrAdapter(err error) error {
	err = errors.New(err.Error())
	return err
}
