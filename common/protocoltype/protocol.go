package protocoltype

type ReportDataFormat struct {
	AgwId       int64     `json:"agwId"` // 网关id
	DeviceId    int64     `json:"deviceId"`
	Timestamp   string    `json:"timestamp"`
	ContentList []Content `json:"contentList"`
}

type Content struct {
	Pid       string      `json:"pid"`
	Type      string      `json:"type"`
	Addr      interface{} `json:"addr"`
	AddrValue interface{} `json:"addrValue"`
}
