package xerr

const (
	OK               uint32 = 200 //成功返回
	RepeatVerifyCode uint32 = 201 //认证
	//UnSupportedProtocolErrorCode  uint32 = 202    //不支持此协议
	ProtocolExistsErrorCode       uint32 = 202    //协议已存在
	NotSupportedProtocolErrorCode uint32 = 203    // 不支持此协议
	NotFoundNidErrCode            uint32 = 204    // 找不到此协议
	NotFoundConfigErrCode         uint32 = 205    // 找不到此协议
	NotFoundNeDeviceErrCode       uint32 = 206    // 找不到设备
	NotFoundFeEdgeDeviceErrCode   uint32 = 207    // 找不到网关
	MacExistsErrCode              uint32 = 208    // mac 地址已存在
	NotFoundFeEdgeDeviceIpErrCode uint32 = 209    // 网关ip地址不存在
	NotFoundMethodErrCode         uint32 = 210    // 找不到 method
	ServerCommonErrorCode         uint32 = 100001 //全局通用错误
	RequestParamErrorCode         uint32 = 100002 //请求参数错误
	RequestRateLimitErrorCode     uint32 = 100003 //触发限流
	RequestSignErrorCode          uint32 = 100004 //签名错误
	ServerUnexpectedlyErrorCode   uint32 = 100005 //签名错误
)
