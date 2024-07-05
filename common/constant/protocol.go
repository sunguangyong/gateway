package constant

import (
	"github.com/gopcua/opcua/ua"
	"xunjikeji.com.cn/gateway/common/convert"
)

const (
	IssuedOn        = 1
	IssuedOff       = 0
	ConfigTypeRead  = 1
	ConfigTypeWrite = 2
	//ConfigTypeReadAndWrite = 0
)

const (
	ModbusRtu = "ModbusRtu"
	ModbusTcp = "ModbusTcp"
	OpcUa     = "OpcUa"
)

const (
	ReadBoolMethod = "ReadBool"
	//ReadByte   = "ReadByte"
	ReadShortMethod  = "ReadShort"
	ReadUShortMethod = "ReadUShort"
	ReadIntMethod    = "ReadInt"
	ReadUIntMethod   = "ReadUInt"
	ReadLongMethod   = "ReadLong"
	ReadULongMethod  = "ReadULong"
	ReadFloatMethod  = "ReadFloat"
	ReadDoubleMethod = "ReadDouble"
	//ReadString = "ReadString"

	WriteBoolMethod   = "WriteBool"
	WriteUShortMethod = "WriteUShort"
	WriteShortMethod  = "WriteShort"
	WriteUIntMethod   = "WriteUInt"
	WriteIntMethod    = "WriteInt"
	WriteULongMethod  = "WriteULong"
	WriteLongMethod   = "WriteLong"
	WriteFloatMethod  = "WriteFloat"
	WriteDoubleMethod = "WriteDouble"
)

var ProtocolNameMap map[string]string
var OpcUaPolicyMap map[string]string
var OpcUaModeMap map[string]ua.MessageSecurityMode

var ModbusRtuBaudRateMap map[int]int
var ModbusRtuWordMap map[int]int

var ModbusRtuBaudRateMapRtuStopBitsMap map[int]int
var ModbusRtuBaudRateMapRtuParityMap map[string]string
var ModbusRtuBaudRateMapRtuComMap map[string]string

var ReadValueLen map[string]uint16
var WriteValueLen map[string]uint16
var ReadValueFunc map[string]convert.BytesToNum

var WriteValueFunc map[string]convert.AnyToBytes

func init() {

	ReadValueFunc = map[string]convert.BytesToNum{
		ReadBoolMethod:   convert.ReadBool,
		ReadShortMethod:  convert.ReadShort,
		ReadUShortMethod: convert.ReadUShort,
		ReadIntMethod:    convert.ReadInt,
		ReadUIntMethod:   convert.ReadUInt,
		ReadLongMethod:   convert.ReadLong,
		ReadULongMethod:  convert.ReadULong,
		ReadFloatMethod:  convert.ReadFloat,
		ReadDoubleMethod: convert.ReadDouble,
	}

	WriteValueFunc = map[string]convert.AnyToBytes{
		WriteBoolMethod:   convert.StringToBoolBytes,
		WriteUShortMethod: convert.StringToUShortBytes,
		WriteShortMethod:  convert.StringToShortBytes,
		WriteUIntMethod:   convert.StringToUIntBytes,
		WriteIntMethod:    convert.StringToIntBytes,
		WriteULongMethod:  convert.StringToULongBytes,
		WriteLongMethod:   convert.StringToLongBytes,
		WriteFloatMethod:  convert.StringToFloatBytes,
		WriteDoubleMethod: convert.StringToDoubleBytes,
	}

	ReadValueLen = map[string]uint16{
		ReadBoolMethod:   1,
		ReadShortMethod:  1,
		ReadUShortMethod: 1,
		ReadIntMethod:    2,
		ReadUIntMethod:   2,
		ReadLongMethod:   4,
		ReadULongMethod:  4,
		ReadFloatMethod:  2,
		ReadDoubleMethod: 4,
	}

	WriteValueLen = map[string]uint16{
		WriteBoolMethod:   1,
		WriteShortMethod:  1,
		WriteUShortMethod: 1,
		WriteIntMethod:    2,
		WriteUIntMethod:   2,
		WriteLongMethod:   4,
		WriteULongMethod:  4,
		WriteFloatMethod:  2,
		WriteDoubleMethod: 4,
	}

	ProtocolNameMap = map[string]string{
		ModbusRtu: ModbusRtu,
		ModbusTcp: ModbusTcp,
		//OpcUa:      OpcUa,
	}

	OpcUaPolicyMap = map[string]string{
		"SecurityPolicyURIPrefix":              "http://opcfoundation.org/UA/SecurityPolicy#",
		"SecurityPolicyURINone":                "http://opcfoundation.org/UA/SecurityPolicy#None",
		"SecurityPolicyURIBasic128Rsa15":       "http://opcfoundation.org/UA/SecurityPolicy#Basic128Rsa15",
		"SecurityPolicyURIBasic256":            "http://opcfoundation.org/UA/SecurityPolicy#Basic256",
		"SecurityPolicyURIBasic256Sha256":      "http://opcfoundation.org/UA/SecurityPolicy#Basic256Sha256",
		"SecurityPolicyURIAes128Sha256RsaOaep": "http://opcfoundation.org/UA/SecurityPolicy#Aes128_Sha256_RsaOaep",
		"SecurityPolicyURIAes256Sha256RsaPss":  "http://opcfoundation.org/UA/SecurityPolicy#Aes256_Sha256_RsaPss",
	}

	OpcUaModeMap = map[string]ua.MessageSecurityMode{
		"MessageSecurityModeInvalid":        ua.MessageSecurityMode(0),
		"MessageSecurityModeNone":           ua.MessageSecurityMode(1),
		"MessageSecurityModeSign":           ua.MessageSecurityMode(2),
		"MessageSecurityModeSignAndEncrypt": ua.MessageSecurityMode(3),
	}

	ModbusRtuBaudRateMap = map[int]int{
		300:    300,
		1200:   1200,
		2400:   2400,
		4800:   4800,
		9600:   9600,
		14400:  14400,
		19200:  19200,
		28800:  28800,
		38400:  38400,
		57600:  57600,
		74880:  74880,
		115200: 115200,
		230400: 230400,
	}

	ModbusRtuBaudRateMapRtuParityMap = map[string]string{
		"None": "N",
		"Even": "E",
		"Odd":  "O",
	}

	ModbusRtuBaudRateMapRtuStopBitsMap = map[int]int{
		1: 1,
		2: 2,
	}

	ModbusRtuWordMap = map[int]int{
		5: 5,
		6: 6,
		7: 7,
		8: 8,
	}

	ModbusRtuBaudRateMapRtuComMap = map[string]string{
		"COM1": "/dev/ttyS3",
		"COM2": "/dev/ttyS4",
		"COM3": "/dev/ttyS5",
		"COM4": "/dev/ttyS6",
	}
}

//func PlcWriteBytes(method string, value string) (result []byte, err error) {
//	funcValue, ok := WriteValueFunc[method]
//	if !ok {
//		err = errors.New(fmt.Sprintf("找不到 method = %s 方法", method))
//	}
//	result, err = funcValue(method, value)
//	return
//}
