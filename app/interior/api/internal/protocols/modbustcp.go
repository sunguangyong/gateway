package protocols

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"xunjikeji.com.cn/gateway/common/protocoltype"

	"xunjikeji.com.cn/gateway/common/util"

	"github.com/zeromicro/go-zero/core/logx"

	"context"

	"xunjikeji.com.cn/gateway/app/interior/api/internal/types"
	"xunjikeji.com.cn/gateway/common/constant"
	"xunjikeji.com.cn/gateway/core/mqtt"

	"github.com/things-go/go-modbus"
)

type ModbusTcp struct {
	Endpoint    string
	JsonOptions types.ModbusTcpJsonOptions
	ConfigData  types.ConfigData
	Client      modbus.Client
	Point       []types.ModbusTcpJsonItems
	cancel      context.CancelFunc
	testChan    chan struct{}
	mqttServer  *mqtt.MqttServer
	mqttConn    *mqtt.MqttConfig
}

func NewModbusTcp(config types.ConfigData, mqttConfig *mqtt.MqttConfig, cancel context.CancelFunc) (m *ModbusTcp,
	err error) {
	m = &ModbusTcp{}
	m.cancel = cancel
	m.Endpoint = config.Config.Endpoint
	m.ConfigData = config
	configJson := config.Config.JsonAccessOptions
	m.testChan = make(chan struct{})

	var jsonOptions types.ModbusTcpJsonOptions
	err = json.Unmarshal([]byte(configJson), &jsonOptions)
	if err != nil {
		return nil, err
	}
	m.JsonOptions = jsonOptions

	for _, v := range config.Point {
		var item types.ModbusTcpJsonItems
		data := v.AccessData

		err = json.Unmarshal([]byte(data), &item)
		if err != nil {
			return nil, err
		}
		m.Point = append(m.Point, item)
	}

	uuid := util.GetUuid()
	mqttServer, err := mqtt.NewMqtt(mqttConfig, uuid)

	if err != nil {
		return m, err
	}

	m.mqttServer = mqttServer

	return m, err
}

func (m *ModbusTcp) Connect() (err error) {
	addr := m.Endpoint
	p := modbus.NewTCPClientProvider(addr, modbus.WithEnableLogger())
	client := modbus.NewClient(p)
	err = client.Connect()
	if err != nil {
		log.Println("connect failed, ", err)
		return err
	}
	m.Client = client
	return nil
}

// ReConnect 如果断开尝试重连
func (m *ModbusTcp) reConnect() {
	m.Client.Close()
	m.Connect()
}

func (m *ModbusTcp) Start(ctx context.Context) {
	ProtocolMapStored(m.ConfigData.Config.Nid, Protocol(m))
	if m.ConfigData.Config.ConfigType == constant.ConfigTypeRead {
		m.Read(ctx)
	}
}

func (m *ModbusTcp) Read(ctx context.Context) (err error) {
	ticker := time.NewTicker(time.Second * 30)
	for {
		select {
		case <-ticker.C:
			m.gather()
		case <-ctx.Done():
			log.Println("proto end nid == ", m.ConfigData.Config.Nid)
			return
		}
	}
}

func (m *ModbusTcp) gather() {
	data := m.packageData()
	m.sendData(data)
}

func (m *ModbusTcp) packageData() protocoltype.ReportDataFormat {
	data := protocoltype.ReportDataFormat{
		AgwId:       m.ConfigData.Config.AgwId,
		DeviceId:    m.ConfigData.Config.DeviceNid,
		ContentList: make([]protocoltype.Content, 0),
		Timestamp:   time.Now().Format(time.DateTime),
	}

	for _, v := range m.Point {
		content := m.getData(v)
		data.ContentList = append(data.ContentList, content)
	}
	return data
}

func (m *ModbusTcp) sendData(data protocoltype.ReportDataFormat) {
	strData, err := json.Marshal(&data)
	if err != nil {
		log.Println(err)
		return
	}
	m.mqttServer.Write(m.mqttServer.Topic, string(strData))
}

func (m *ModbusTcp) End() {
	m.cancel()
	m.Client.Close()
	ProtocolMapDel(m.ConfigData.Config.Nid)
}

func (m *ModbusTcp) ReadTest() (data protocoltype.ReportDataFormat) {
	// m.testChan <- struct{}{}
	data = m.packageData()
	return
}

func (m *ModbusTcp) Write(request *types.ProtocolWriteTestRequest) (err error) {
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
		err = m.Client.WriteSingleCoil(m.JsonOptions.SlaveId, request.Address, flag)
	} else {
		err = m.Client.WriteMultipleRegistersBytes(m.JsonOptions.SlaveId, request.Address, quantity, dataBytes)
	}
	return
}

func (m *ModbusTcp) CreateProto(config types.ConfigData, mqttConfig *mqtt.MqttConfig,
	cancel context.CancelFunc) (protocol Protocol, err error) {
	protocol, err = NewModbusTcp(config, mqttConfig, cancel)
	return
}

func (m *ModbusTcp) getData(v types.ModbusTcpJsonItems) (content protocoltype.Content) {

	var (
		results []byte
		err     error
		value   interface{}
	)

	content = protocoltype.Content{}
	content.Addr = v.DataName
	content.Type = v.Method

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

	if v.Method == constant.ReadBoolMethod {
		results, err = m.Client.ReadCoils(m.JsonOptions.SlaveId, v.Address, quantity)
	} else {
		results, err = m.Client.ReadHoldingRegistersBytes(m.JsonOptions.SlaveId, v.Address, quantity)
	}

	fmt.Println("addr ==== ", v.Address, "result === ", results)

	if err != nil {
		logx.Errorf("modbus tcp err %s", err.Error())
		m.reConnect()
		return
	}

	value, err = PlcReadeBytes(v.Method, results)

	if err != nil {
		log.Println(err)
	}
	content.AddrValue = value

	return
}
