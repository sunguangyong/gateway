package mqtt

import (
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"github.com/zeromicro/go-zero/core/threading"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

const (
	qos0 = iota
	qos1
	qos2
)

type MqttConfig struct {
	Host     string
	UserName string
	Password string
	Topic    string
}

type MqttServer struct {
	Host     string
	Name     string
	Topic    string
	Qos      int
	ClientId string
	UserName string
	Password string
	Client   MQTT.Client
}

func NewMqtt(mqttConfig *MqttConfig, clientId string) (svr *MqttServer, err error) {
	svr = new(MqttServer)
	mqttName := "MQTT"
	svr.Host = mqttConfig.Host
	// ClientId必须是唯一的，并且不同的MQTT客户端必须使用不同的ClientId
	svr.ClientId = clientId
	//svr.clientId = common.GetUuid()
	svr.UserName = mqttConfig.UserName
	svr.Password = mqttConfig.Password
	svr.Name = mqttName
	svr.Topic = mqttConfig.Topic

	client, err := svr.NewClient()

	if err != nil {
		log.Println("mqtt client err ==== ", err)
		return nil, err
	}

	log.Println("mqtt success")
	svr.Client = client
	return svr, nil
}

func (svr *MqttServer) NewClient() (client MQTT.Client, err error) {
	connOpts := MQTT.NewClientOptions().AddBroker(svr.Host).SetClientID(svr.ClientId).SetCleanSession(true)
	if svr.UserName != "" {
		connOpts.SetUsername(svr.UserName)
		if svr.Password != "" {
			connOpts.SetPassword(svr.Password)
		}
	}
	connOpts.SetAutoReconnect(true)                    // 启用自动重连
	connOpts.SetConnectRetry(true)                     // 允许在建立连接之前发送数据
	connOpts.SetMaxReconnectInterval(10 * time.Second) // 设置最大重连间隔

	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	connOpts.SetTLSConfig(tlsConfig)
	client = MQTT.NewClient(connOpts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		err = token.Error()
		log.Fatalln(err)
	}
	return client, err
}

func (svr *MqttServer) Write(topic, payload string) {
	threading.GoSafe(func() {
		// 发布MQTT消息
		token := svr.Client.Publish(topic, byte(svr.Qos), true, payload)
		ok := token.Wait()
		if ok {
			log.Printf("write success ===%s \n", payload)
		} else {
			log.Printf("write fail ===%s \n", payload)
		}
	})
}

func (svr *MqttServer) Read(topic string, callback MQTT.MessageHandler) {
	if token := svr.Client.Subscribe(topic, byte(qos0), callback); token.Wait() && token.Error() != nil {
		log.Fatalln(token.Error())
	}
}

// ReadMessage 处理MQTT消息的回调函数
func ReadMessage(client MQTT.Client, message MQTT.Message) {
	fmt.Printf("Received message: %s\n", message.Payload())
}
