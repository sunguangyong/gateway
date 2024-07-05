package types

type InteriorConfigOnResponse struct {
	Code int                 `json:"code"`
	Msg  string              `json:"msg"`
	Data ProtoDeployResponse `json:"data"`
}

type ProtoDeployResponse struct {
	FailArray    []DeployMsg `json:"failArray"`    // 失败列表
	SuccessArray []DeployMsg `json:"successArray"` // 成功列表
}

type DeployMsg struct {
	ConfigData NeDeviceDataAccessConfig `json:"configData"`
	Msg        string                   `json:"msg"`
}

type ProtoDeployRequest struct {
	DataList []DeployDevice `json:"list"`
}

type DeployDevice struct {
	DeviceNid  int64        `json:"deviceNid"` // 设备ID
	ConfigData []ConfigData `json:"configData"`
}

type ProtoDeleteRequest struct {
	Nids []int64 `json:"nids"` // 协议id
}

type ProtoDeleteResponse struct {
	SucceedNids []int64 `json:"succeedNids"`
	FailNids    []int64 `json:"failNids"`
}

type ConfigData struct {
	Config NeDeviceDataAccessConfig `json:"config"`
	Point  []NeDeviceDataAccessItem `json:"point"`
}

type NeDeviceDataAccessConfig struct {
	Nid               int64  `json:"nid"`                 // 主键ID
	DeviceNid         int64  `json:"device_nid"`          // 设备ID
	ConfigType        int64  `json:"config_type"`         // 0-读/写,1-只读,2-只写
	ConfigId          string `json:"config_id"`           // 配置ID
	ConfigName        string `json:"config_name"`         // 配置名称
	Endpoint          string `json:"endpoint"`            // 连接URL
	Protocol          string `json:"protocol"`            // 访问协议：ModbusTcp, ModbusRtu, OpcUa, HTTP, MTCONNECT, MITSUBISH_MC, SIEMENS_S7, OMRON_SINS
	JsonAccessOptions string `json:"json_access_options"` // 访问连接基本配置(JSON格式)
	Timeout           int64  `json:"timeout"`             // 连接超时时间
	AgwId             int64  `json:"agw_id"`              // 所属网关ID
	TenantId          int64  `json:"tenant_id"`           // 租户ID
	ProfileNid        int64  `json:"profile_nid"`         // 设备配置文件ID
	CreateTime        string `json:"create_time"`         // 创建时间
	CreateBy          int64  `json:"create_by"`           // 创建者ID
	Issued            int64  `json:"issued"`              // 是否已下发：1-是,0-否
	IssueTime         string `json:"issue_time"`          // 下发时间
}

type NeDeviceDataAccessItem struct {
	Nid        int64  `json:"nid"`         // 主键
	ConfigNid  int64  `json:"config_nid"`  // ne_device_data_access_config.id
	DeviceNid  int64  `json:"device_nid"`  // 设备ID
	AgwId      int64  `json:"agw_id"`      // 所属网关ID
	TenantId   int64  `json:"tenant_id"`   // 租户ID
	ConfigType int64  `json:"config_type"` // 0-读/写,1-只读,2-只写
	AccessData string `json:"access_data"` // 读/写配置项(JSON对象)
	CreateTime string `json:"create_time"` // 创建时间
}

type InteriorHeartBeatRequest struct {
}

type InteriorHeartBeatResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
