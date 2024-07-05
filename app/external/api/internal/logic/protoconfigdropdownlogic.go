package logic

import (
	"context"

	"xunjikeji.com.cn/gateway/common/constant"

	"xunjikeji.com.cn/gateway/app/external/api/internal/svc"
	"xunjikeji.com.cn/gateway/app/external/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProtoConfigDropDownLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProtoConfigDropDownLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProtoConfigDropDownLogic {
	return &ProtoConfigDropDownLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProtoConfigDropDownLogic) ProtoConfigDropDown(req *types.ProtoConfigDropDownRequest) (resp *types.ProtoConfigDropDownResponse, err error) {
	resp = &types.ProtoConfigDropDownResponse{
		DropDown: make(map[string]interface{}),
	}
	protoName := req.Protocol
	switch protoName {
	case constant.ModbusTcp:
		modbusTcpDropDown(resp)

	case constant.ModbusRtu:
		modbusRtuDropDown(resp)

	case constant.OpcUa:
		opcuaDropDown(resp)
	}

	return
}

func modbusRtuDropDown(resp *types.ProtoConfigDropDownResponse) {

	var baudRateList []types.Dropdown
	var dataBitsList []types.Dropdown
	var stopBitsList []types.Dropdown
	var parityList []types.Dropdown
	var comList []types.Dropdown

	for k, v := range constant.ModbusRtuBaudRateMap {
		baudRateList = append(baudRateList, types.Dropdown{
			Value: v,
			Label: k,
		})
	}

	for k, v := range constant.ModbusRtuWordMap {
		dataBitsList = append(dataBitsList, types.Dropdown{
			Value: v,
			Label: k,
		})
	}

	for k, v := range constant.ModbusRtuBaudRateMapRtuStopBitsMap {
		stopBitsList = append(stopBitsList, types.Dropdown{
			Value: v,
			Label: k,
		})
	}

	for k, v := range constant.ModbusRtuBaudRateMapRtuParityMap {
		parityList = append(parityList, types.Dropdown{
			Value: v,
			Label: k,
		})
	}

	for k, v := range constant.ModbusRtuBaudRateMapRtuComMap {
		comList = append(comList, types.Dropdown{
			Value: v,
			Label: k,
		})
	}

	// 排序
	sortDropdown(baudRateList)
	sortDropdown(dataBitsList)
	sortDropdown(stopBitsList)
	sortDropdown(parityList)
	sortDropdown(comList)

	resp.DropDown["baudRate"] = baudRateList
	resp.DropDown["dataBits"] = dataBitsList
	resp.DropDown["stopBits"] = stopBitsList
	resp.DropDown["parity"] = parityList
	resp.DropDown["com"] = comList
}

func modbusTcpDropDown(resp *types.ProtoConfigDropDownResponse) {
}

func opcuaDropDown(resp *types.ProtoConfigDropDownResponse) {
	var policyList []types.Dropdown
	var modeList []types.Dropdown

	for k, v := range constant.OpcUaPolicyMap {
		policyList = append(policyList, types.Dropdown{
			Value: v,
			Label: k,
		})
	}

	for k, _ := range constant.OpcUaModeMap {
		modeList = append(modeList, types.Dropdown{
			Value: k,
			Label: k,
		})
	}

	sortDropdown(policyList)
	sortDropdown(modeList)

	resp.DropDown["policy"] = policyList
	resp.DropDown["mode"] = modeList
}
