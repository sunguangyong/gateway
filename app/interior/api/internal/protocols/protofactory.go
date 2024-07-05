package protocols

import (
	"context"
	"sync"

	"xunjikeji.com.cn/gateway/common/protocoltype"

	"github.com/zeromicro/go-zero/core/threading"

	"xunjikeji.com.cn/gateway/core/mqtt"

	"xunjikeji.com.cn/gateway/app/interior/api/internal/types"
	"xunjikeji.com.cn/gateway/common/constant"
	"xunjikeji.com.cn/gateway/common/xerr"
)

var protocolMap map[int64]Protocol
var protocolLock sync.Mutex

func init() {
	protocolMap = make(map[int64]Protocol)
}

//go:generate mockgen -destination=../mock_protocols/mock_protofactory.go -package=mock -source=protofactory.go
type Protocol interface {
	Connect() (err error)                                      // 连接
	Start(ctx context.Context)                                 // 开始采集
	End()                                                      // 结束采集
	ReadTest() (data protocoltype.ReportDataFormat)            // 采集测试
	Write(request *types.ProtocolWriteTestRequest) (err error) // 写plc
}

type ProtoFactory interface {
	CreateProto()
}

// 获取协议
func getFactory(configData types.ConfigData, cancel context.CancelFunc) (protocol Protocol, err error) {
	var mqttConfig mqtt.MqttConfig
	mqttConfig.Host = "tcp://112.111.92.167:1883"
	mqttConfig.UserName = "admin"
	mqttConfig.Password = "123456@2023"
	mqttConfig.Topic = "iot/test"

	switch configData.Config.Protocol {
	case constant.ModbusTcp:
		modbusTcp := new(ModbusTcp)
		protocol, err = modbusTcp.CreateProto(configData, &mqttConfig, cancel)

	case constant.ModbusRtu:
		modbusRtu := new(ModbusRtu)
		protocol, err = modbusRtu.CreateProto(configData, &mqttConfig, cancel)

	case constant.OpcUa:
		opcUaFac := new(OpcUa)
		protocol, err = opcUaFac.CreateProto(configData, &mqttConfig, cancel)

	default:
		return nil, xerr.NotSupportedProtocolError
	}
	return
}

func ProtocolFlow(configData types.ConfigData) (ok bool, err error) {
	ctx, cancel := context.WithCancel(context.Background())

	// 使用 defer 在函数退出时判断是否需要调用 cancel
	defer func() {
		if err != nil {
			cancel()
		}
	}()

	if err != nil {
		return false, err
	}

	protocol, err := getFactory(configData, cancel)

	if err != nil {
		return false, err
	}

	if protocol == nil {
		return false, xerr.NotSupportedProtocolError
	}

	err = protocol.Connect()

	if err != nil {
		return false, err
	}

	// 如果此协议正在采集先停掉
	ProtocolStop(configData.Config.Nid)

	// 启动采集
	threading.GoSafe(func() {
		protocol.Start(ctx)
	})

	return true, nil
}

func ProtocolStop(nid int64) bool {
	proto, ok := ProtocolMapGet(nid)
	if ok {
		proto.End()
	}
	return true
}

func ProtocolReadTest(nid int64) (data protocoltype.ReportDataFormat) {
	proto, ok := ProtocolMapGet(nid)
	if ok {
		data = proto.ReadTest()
	}
	return
}

func ProtocolMapGet(nid int64) (proto Protocol, ok bool) {
	protocolLock.Lock()
	defer protocolLock.Unlock()
	proto, ok = protocolMap[nid]
	return
}

func ProtocolMapDel(nid int64) {
	protocolLock.Lock()
	defer protocolLock.Unlock()
	delete(protocolMap, nid)
}

func ProtocolMapStored(nid int64, proto Protocol) {
	protocolLock.Lock()
	defer protocolLock.Unlock()
	protocolMap[nid] = proto
}
