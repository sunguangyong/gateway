package xerr

var message map[uint32]string

func init() {
	message = make(map[uint32]string)
	message[OK] = "success"
	message[ServerCommonErrorCode] = "server error"
	message[RequestSignErrorCode] = "signature error"
}

func MapErrMsg(errCode uint32) string {
	if msg, ok := message[errCode]; ok {
		return msg
	} else {
		return "server error"
	}
}

func IsCodeErr(errCode uint32) bool {
	if _, ok := message[errCode]; ok {
		return true
	} else {
		return false
	}
}
