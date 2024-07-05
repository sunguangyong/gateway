package xerr

import "fmt"

type CtxErr struct {
	Err     error
	Content interface{}
}

func NewCtxErr(err error, v interface{}) *CtxErr {
	return &CtxErr{
		Err:     err,
		Content: v,
	}
}

func (e *CtxErr) Error() string {
	if e != nil && e.Err != nil {
		return e.Err.Error()
	}
	return fmt.Sprintf("unkonw error")
}
