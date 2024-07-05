package types

import "xunjikeji.com.cn/gateway/common/constant"

var protoItemsMap map[string]interface{}
var protoJsonOptionsMap map[string]interface{}

func init() {
	// 每次增加新协议需要在此增加类型
	protoItemsMap = map[string]interface{}{
		constant.ModbusRtu: ModbusRtuJsonItems{},
		constant.ModbusTcp: ModbusTcpJsonItems{},
		constant.OpcUa:     OpcUaJsonItems{},
	}

	protoJsonOptionsMap = map[string]interface{}{
		constant.ModbusRtu: ModbusRtuJsonOptions{},
		constant.ModbusTcp: ModbusTcpJsonOptions{},
		constant.OpcUa:     OpcUaJsonOptions{},
	}

}

func GetJsonOptionsType(protoType string) (v interface{}, ok bool) {
	v, ok = protoJsonOptionsMap[protoType]
	return
}

func GetItemType(protoType string) (v interface{}, ok bool) {
	v, ok = protoItemsMap[protoType]
	return
}
