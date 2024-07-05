package protocols

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"xunjikeji.com.cn/gateway/common/protocoltype"

	"xunjikeji.com.cn/gateway/common/constant"
	"xunjikeji.com.cn/gateway/common/util"

	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"

	"context"

	"xunjikeji.com.cn/gateway/app/interior/api/internal/types"
	"xunjikeji.com.cn/gateway/core/mqtt"
)

type OpcUa struct {
	Endpoint    string
	ConfigData  types.ConfigData
	JsonOptions OpcUaJsonOptions
	Client      *opcua.Client
	Point       []OpcUaJsonItems
	cancel      context.CancelFunc
	testChan    chan struct{}
	mqttServer  *mqtt.MqttServer
	mqttConn    *mqtt.MqttConfig
}

type OpcUaJsonOptions struct {
	Protocol string `json:"protocol"`
	UserName string `json:"userName"`
	PassWd   string `json:"passWd"`
	Policy   string `json:"policy"`
	Mode     string `json:"mode"`
	Cert     bool   `json:"cert"`
}

type OpcUaJsonItems struct {
	Protocol string `json:"protocol"`
	NodeId   string `json:"nodeId"`
}

func NewOpcUa(config types.ConfigData, mqttConfig *mqtt.MqttConfig, cancel context.CancelFunc) (o *OpcUa, err error) {
	o = &OpcUa{}
	o.cancel = cancel
	o.ConfigData = config
	o.Endpoint = config.Config.Endpoint
	o.Point = make([]OpcUaJsonItems, 0, len(config.Point))
	configJson := config.Config.JsonAccessOptions
	o.testChan = make(chan struct{})

	var jsonOptions OpcUaJsonOptions
	err = json.Unmarshal([]byte(configJson), &jsonOptions)

	for _, v := range config.Point {
		var item OpcUaJsonItems
		data := v.AccessData

		err = json.Unmarshal([]byte(data), &item)
		if err != nil {
			log.Println("err == ", err.Error())
			continue
		}
		o.Point = append(o.Point, item)
	}

	if err != nil {
		return nil, err
	}

	o.JsonOptions = jsonOptions

	uuid := util.GetUuid()
	mqttServer, err := mqtt.NewMqtt(mqttConfig, uuid)

	if err != nil {
		return o, err
	}

	o.mqttServer = mqttServer

	return o, err
}

func (o *OpcUa) CreateProto(config types.ConfigData, mqttConfig *mqtt.MqttConfig,
	cancel context.CancelFunc) (protocol Protocol, err error) {
	protocol, err = NewOpcUa(config, mqttConfig, cancel)
	return
}

func (o *OpcUa) Connect() (err error) {
	endpoint := o.Endpoint

	endpoints, err := opcua.GetEndpoints(context.Background(), endpoint)

	if err != nil {
		log.Println("OPC GetEndpoints: %w", err.Error())
		return
	}

	c, err := util.GenerateCert() // This is where you generate the certificate

	if err != nil {
		err = fmt.Errorf("generateCert: %w", err)
	}

	opts := []opcua.Option{
		opcua.SessionTimeout(30 * time.Minute),
		opcua.AutoReconnect(true),
		opcua.ReconnectInterval(time.Second * 10),
		opcua.Lifetime(30 * time.Minute),
		opcua.RequestTimeout(3 * time.Second),
	}

	if o.JsonOptions.UserName != "" && o.JsonOptions.PassWd != "" {
		opts = append(opts, opcua.AuthUsername(o.JsonOptions.UserName, o.JsonOptions.PassWd))
	} else {
		opts = append(opts, opcua.SecurityMode(ua.MessageSecurityModeNone))
	}

	if o.JsonOptions.Mode != "" {
		opts = append(opts, opcua.SecurityMode(constant.OpcUaModeMap[o.JsonOptions.Mode]))
	}

	if o.JsonOptions.Policy != "" {
		opts = append(opts, opcua.SecurityPolicy(o.JsonOptions.Policy))
	}

	if o.JsonOptions.Mode != "" && o.JsonOptions.Policy != "" && o.JsonOptions.UserName != "" && o.JsonOptions.PassWd != "" {
		ep := opcua.SelectEndpoint(endpoints, o.JsonOptions.Policy, constant.OpcUaModeMap[o.JsonOptions.Mode])
		opts = append(opts, opcua.SecurityFromEndpoint(ep, ua.UserTokenTypeUserName))
	}

	if o.JsonOptions.Cert {
		pk, ok := c.PrivateKey.(*rsa.PrivateKey) // This is where you set the private key
		if !ok {
			err = fmt.Errorf("invalid private key")
			return
		}
		cert := c.Certificate[0]
		opts = append(opts, opcua.PrivateKey(pk))
		opts = append(opts, opcua.Certificate(cert))
	}

	client := opcua.NewClient(endpoint, opts...)
	o.Client = client

	ctx := context.Background()

	if err = o.Client.Connect(ctx); err != nil {
		log.Println("opcua err === ", err.Error())
		return err
	}

	return nil
}

// ReConnect 如果断开尝试重连
func (o *OpcUa) ReConnect() {
	o.Client.Close()
	o.Connect()
}

func (o *OpcUa) Start(ctx context.Context) {
	ProtocolMapStored(o.ConfigData.Config.Nid, Protocol(o))
	ticker := time.NewTicker(time.Second * 5)
	defer o.Client.Close()
	for {
		select {
		case <-ticker.C:
			o.gather()
		case <-ctx.Done():
			ProtocolMapDel(o.ConfigData.Config.Nid)
			log.Println("proto end nid == ", o.ConfigData.Config.Nid)
			return
		}
	}
}

func (o *OpcUa) gather() {
	data := o.packageData()
	o.sendData(data)
}

func (o *OpcUa) packageData() (data protocoltype.ReportDataFormat) {
	data = protocoltype.ReportDataFormat{
		AgwId:       o.ConfigData.Config.AgwId,
		DeviceId:    o.ConfigData.Config.DeviceNid,
		Timestamp:   time.Now().Format(time.DateTime),
		ContentList: make([]protocoltype.Content, 0),
	}

	req := &ua.ReadRequest{
		MaxAge:             2000,
		NodesToRead:        []*ua.ReadValueID{},
		TimestampsToReturn: ua.TimestampsToReturnBoth,
	}

	for _, v := range o.Point {
		nodeId, err := ua.ParseNodeID(v.NodeId)
		if err != nil {
			log.Printf("invalid node id: %v", v.NodeId)
			continue
		}
		req.NodesToRead = append(req.NodesToRead, &ua.ReadValueID{NodeID: nodeId})
	}

	resp, err := o.Client.Read(req)

	if err != nil {
		log.Println(err.Error())
		o.ReConnect()
		return
	}

	for i, v := range resp.Results {

		if v.Status != ua.StatusOK {
			log.Println(" nodeId =", v.Value.NodeID(), "Status not OK: %v", v.Status)
			continue
		}

		content := protocoltype.Content{}
		content.Addr = req.NodesToRead[i].NodeID
		content.AddrValue = v.Value.Value()
		data.ContentList = append(data.ContentList, content)
	}
	return data
}

func (o *OpcUa) sendData(data protocoltype.ReportDataFormat) {
	strData, err := json.Marshal(&data)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("topic ===== ", o.mqttServer.Topic)
	o.mqttServer.Write(o.mqttServer.Topic, string(strData))
}

func (o *OpcUa) End() {
	o.cancel()
	o.Client.Close()
	ProtocolMapDel(o.ConfigData.Config.Nid)
}

func (o *OpcUa) ReadTest() (data protocoltype.ReportDataFormat) {
	// o.testChan <- struct{}{}
	data = o.packageData()
	return
}

func (o *OpcUa) Write(request *types.ProtocolWriteTestRequest) (err error) {
	return
}
