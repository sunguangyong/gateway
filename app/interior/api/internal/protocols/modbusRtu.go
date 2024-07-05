package protocols

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"xunjikeji.com.cn/gateway/common/util"

	"xunjikeji.com.cn/gateway/common/constant"
	"xunjikeji.com.cn/gateway/common/protocoltype"

	"xunjikeji.com.cn/gateway/app/interior/api/internal/types"
	//"xunjikeji.com.cn/gateway/app/interior/api/internal/transmit/mqtt"
	"xunjikeji.com.cn/gateway/core/mqtt"

	"github.com/goburrow/serial"
	"github.com/things-go/go-modbus"
)

type ModbusRtu struct {
	Endpoint    string
	JsonOptions types.ModbusRtuJsonOptions
	ConfigData  types.ConfigData
	Client      modbus.Client
	Point       []types.ModbusRtuJsonItems
	cancel      context.CancelFunc
	testChan    chan struct{}
	mqttServer  *mqtt.MqttServer
	mqttConn    *mqtt.MqttConfig
}

//type ModbusRtuJsonOptions struct {
//	Protocol string `json:"protocol"` // 协议名称
//	BaudRate int    `json:"baudRate"` // 波特率
//	DataBits int    `json:"dataBits"` // 数据为
//	StopBits int    `json:"stopBits"` // 停止位
//	Parity   string `json:"parity"`   // 校验方式
//}
//
//type ModbusRtuJsonItems struct {
//	Protocol string `json:"protocol"` // 协议名称
//	SlaveId  byte   `json:"slaveId"`  // 从站id
//	Address  uint16 `json:"address"`  // 起始地址
//	Quantity uint16 `json:"quantity"` // 读取数量
//}

func NewModbusRtu(config types.ConfigData, mqttConfig *mqtt.MqttConfig, cancel context.CancelFunc) (r *ModbusRtu,
	err error) {
	r = &ModbusRtu{}
	r.cancel = cancel
	r.Endpoint = config.Config.Endpoint
	r.ConfigData = config
	configJson := config.Config.JsonAccessOptions
	r.testChan = make(chan struct{})

	var jsonOptions types.ModbusRtuJsonOptions
	err = json.Unmarshal([]byte(configJson), &jsonOptions)
	if err != nil {
		return nil, err
	}
	r.JsonOptions = jsonOptions

	for _, v := range config.Point {
		var item types.ModbusRtuJsonItems
		data := v.AccessData

		err = json.Unmarshal([]byte(data), &item)
		if err != nil {
			return nil, err
		}
		r.Point = append(r.Point, item)
	}
	uuid := util.GetUuid()
	mqttServer, err := mqtt.NewMqtt(mqttConfig, uuid)

	if err != nil {
		return r, err
	}

	r.mqttServer = mqttServer

	return r, err
}

func (r *ModbusRtu) CreateProto(config types.ConfigData, mqttConfig *mqtt.MqttConfig,
	cancel context.CancelFunc) (protocol Protocol, err error) {
	protocol, err = NewModbusRtu(config, mqttConfig, cancel)
	return
}

// Connect 创建连接
func (r *ModbusRtu) Connect() (err error) {
	p := modbus.NewRTUClientProvider(modbus.WithEnableLogger(),
		modbus.WithSerialConfig(serial.Config{
			Address:  r.Endpoint,
			BaudRate: r.JsonOptions.BaudRate,
			DataBits: r.JsonOptions.DataBits,
			StopBits: r.JsonOptions.StopBits,
			Parity:   r.JsonOptions.Parity,
			Timeout:  modbus.SerialDefaultTimeout,
		}))
	client := modbus.NewClient(p)
	err = client.Connect()
	if err != nil {
		return
	}
	r.Client = client
	return
}

// ReConnect 如果断开尝试重连
func (r *ModbusRtu) reConnect() {
	r.Client.Close()
	r.Connect()
}

// Start 开始
func (r *ModbusRtu) Start(ctx context.Context) {
	ProtocolMapStored(r.ConfigData.Config.Nid, Protocol(r))

	if r.ConfigData.Config.ConfigType == constant.ConfigTypeRead {
		r.Read(ctx)
	}

}

func (r *ModbusRtu) Read(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-ticker.C:
			r.gather()
		case <-ctx.Done():
			return
		}
	}
}

func (r *ModbusRtu) Write(request *types.ProtocolWriteTestRequest) (err error) {
	quantity, ok := constant.WriteValueLen[request.Method]

	if !ok {
		return
	}

	dataBytes, err := PlcWriteBytes(request.Method, request.Value)

	fmt.Println("dataBytes=======", dataBytes)

	if err != nil {
		return
	}

	if request.Method == constant.WriteBoolMethod {
		if len(dataBytes) == 0 {
			return
		}

		var flag bool
		if dataBytes[0] > 0 {
			flag = true
		}
		err = r.Client.WriteSingleCoil(r.JsonOptions.SlaveId, request.Address, flag)
	} else {
		err = r.Client.WriteMultipleRegistersBytes(r.JsonOptions.SlaveId, request.Address, quantity, dataBytes)
	}
	return
}

// End 结束
func (r *ModbusRtu) End() {
	r.cancel()
	ProtocolMapDel(r.ConfigData.Config.Nid)
	r.Client.Close()
}

func (r *ModbusRtu) ReadTest() (data protocoltype.ReportDataFormat) {
	data = r.packageData()
	return
}

// 采集
func (r *ModbusRtu) gather() {
	data := r.packageData()
	r.sendData(data)
}

func (r *ModbusRtu) packageData() protocoltype.ReportDataFormat {
	//uTime := time.Now().Format(constant.LayOut)
	data := protocoltype.ReportDataFormat{
		AgwId:       r.ConfigData.Config.AgwId,
		DeviceId:    r.ConfigData.Config.DeviceNid,
		Timestamp:   time.Now().Format(time.DateTime),
		ContentList: make([]protocoltype.Content, 0),
	}

	for _, v := range r.Point {
		content := r.getData(v)
		data.ContentList = append(data.ContentList, content)
	}
	return data
}

// sendData 发送消息
func (r *ModbusRtu) sendData(data protocoltype.ReportDataFormat) {
	strData, err := json.Marshal(&data)
	if err != nil {
		return
	}
	r.mqttServer.Write(r.mqttServer.Topic, string(strData))
}

func (r *ModbusRtu) getData(v types.ModbusRtuJsonItems) (content protocoltype.Content) {

	var (
		value   interface{}
		results []byte
		err     error
	)

	content = protocoltype.Content{}
	content.Addr = v.DataName

	defer func() {
		if err != nil {
			log.Println("modbus tcp err ", err)
		}
	}()

	quantity, ok := constant.ReadValueLen[v.Method]

	if !ok {
		return
	}

	// ReadCoils 01 COIL STATUS 读取线圈状态
	// ReadDiscreteInputs 02 INPTUT STATUS
	// ReadHoldingRegisters 03 HOLDING REGISTER
	//ReadInputRegisters 04 INPTUT REGISTER

	fmt.Println("addr ====", v.Address, "quantity ==== ", quantity)

	if v.Method == constant.ReadBoolMethod {
		results, err = r.Client.ReadCoils(r.JsonOptions.SlaveId, v.Address, quantity*2)
	} else {
		results, err = r.Client.ReadHoldingRegistersBytes(r.JsonOptions.SlaveId, v.Address, quantity)
	}

	// TODO 判断错误类型是否重连
	if err != nil {
		// 超时重连
		//if err.Error() == serial.ErrTimeout.Error() {
		//	r.reConnect()
		//}
		fmt.Println("modbus rtu err", err.Error())
		r.reConnect()
		return
	}

	value, err = PlcReadeBytes(v.Method, results)

	if err != nil {
		log.Println(err)
	}

	content.AddrValue = value

	// 防止串口超时
	time.Sleep(50 * time.Millisecond)

	return
}
