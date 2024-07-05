package logic

import (
	"context"
	"log"
	"time"

	"xunjikeji.com.cn/gateway/common/client"
	"xunjikeji.com.cn/gateway/common/util"

	"xunjikeji.com.cn/gateway/common/xerr"

	"xunjikeji.com.cn/gateway/app/external/api/internal/svc"
	"xunjikeji.com.cn/gateway/app/external/model"
	"xunjikeji.com.cn/gateway/common/constant"

	"xunjikeji.com.cn/gateway/app/external/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProtocolDeployLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProtocolDeployLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProtocolDeployLogic {
	return &ProtocolDeployLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ProtocolDeploy handles the deployment of protocols
func (l *ProtocolDeployLogic) ProtocolDeploy(req *types.ProtocolDeployRequest) (*types.ProtocolDeployResponse, error) {
	resp := &types.ProtocolDeployResponse{
		SuccessArray: make([]types.DeployArray, 0),
		FailArray:    make([]types.DeployArray, 0),
	}

	host, err := l.svcCtx.EdgeDevice.GetHost(l.ctx, req.AgwId)
	if err != nil || host == "" {
		return resp, xerr.NotFoundFeEdgeDeviceIpErr
	}

	host = util.GetInteriorHost(host, constant.InteriorPort)

	request, deviceMap, configMap, err := l.buildDeployRequest(req.NeDeviceIds)

	if err != nil {
		return resp, err
	}

	body, err := l.sendDeployRequest(request, host)

	if err != nil {
		return resp, err
	}

	l.updateConfigArrays(resp, body, deviceMap, configMap)

	return resp, nil
}

// 构建请求参数
func (l *ProtocolDeployLogic) buildDeployRequest(neDeviceIds []int64) (*types.ProtoDeployRequest, map[int64]*model.NeDevice, map[int64]*model.NeDeviceDataAccessConfig, error) {
	request := &types.ProtoDeployRequest{
		DataList: make([]types.DeployDevice, 0),
	}
	deviceMap := make(map[int64]*model.NeDevice)
	configMap := make(map[int64]*model.NeDeviceDataAccessConfig)

	for _, deviceId := range neDeviceIds {
		device, err := l.svcCtx.NeDevice.FindOne(l.ctx, deviceId)
		if err != nil || device == nil {
			continue
		}

		deviceMap[deviceId] = device
		deployDevice := types.DeployDevice{DeviceNid: deviceId}

		filter := model.NeDeviceDataAccessConfig{DeviceNid: deviceId}
		configArray, err := l.svcCtx.AccessConfig.CommonFilterFind(l.ctx, filter, model.NewPageInfo(0, 0))
		if err != nil {
			log.Println(err)
			continue
		}

		for _, config := range configArray {
			configMap[config.Nid] = config
			data := l.DeploySingle(config)
			deployDevice.ConfigData = append(deployDevice.ConfigData, data)
		}

		request.DataList = append(request.DataList, deployDevice)
	}

	return request, deviceMap, configMap, nil
}

// 发送请求
func (l *ProtocolDeployLogic) sendDeployRequest(request *types.ProtoDeployRequest, host string) (*types.InteriorConfigOnResponse, error) {
	var body types.InteriorConfigOnResponse
	err := client.PostJson(request, &body, host, constant.ProtoDeployPath)
	return &body, err
}

// 更新协议部署状态
func (l *ProtocolDeployLogic) updateConfigArrays(resp *types.ProtocolDeployResponse,
	body *types.InteriorConfigOnResponse, deviceMap map[int64]*model.NeDevice, configMap map[int64]*model.NeDeviceDataAccessConfig) {
	for _, v := range body.Data.SuccessArray {

		l.updateDeviceAndConfig(v, deviceMap, configMap, resp.SuccessArray, constant.IssuedOn)
	}

	for _, v := range body.Data.FailArray {
		l.updateDeviceAndConfig(v, deviceMap, configMap, resp.FailArray, constant.IssuedOff)
	}
}

// 更新 设备部署状态
func (l *ProtocolDeployLogic) updateDeviceAndConfig(data types.DeployMsg, deviceMap map[int64]*model.NeDevice,
	configMap map[int64]*model.NeDeviceDataAccessConfig, deployArray []types.DeployArray, issuedStatus int64) {
	deploy := types.DeployArray{}
	configNid := data.ConfigData.Nid
	deviceNid := data.ConfigData.DeviceNid
	device, deviceOk := deviceMap[deviceNid]
	config, configOk := configMap[configNid]

	if deviceOk {
		deploy.DeviceNid = deviceNid
		deploy.DeviceName = device.DeviceName
		device.Issued = issuedStatus
		device.IssueTime = time.Now()
		l.svcCtx.NeDevice.Update(l.ctx, device)
	}

	if configOk {
		deploy.Nid = configNid
		deploy.ConfigName = config.ConfigName
		deploy.DeviceName = device.DeviceName
		config.IssueTime = time.Now()
		config.Issued = issuedStatus
		l.svcCtx.AccessConfig.Update(l.ctx, config)
	}

	deployArray = append(deployArray, deploy)
}

func (l *ProtocolDeployLogic) DeploySingle(configData *model.NeDeviceDataAccessConfig) types.ConfigData {
	data := types.ConfigData{
		Config: types.NeDeviceDataAccessConfig{
			Nid:               configData.Nid,
			DeviceNid:         configData.DeviceNid,
			ConfigType:        configData.ConfigType,
			ConfigId:          configData.ConfigId,
			ConfigName:        configData.ConfigName,
			Endpoint:          configData.Endpoint,
			Protocol:          configData.Protocol,
			JsonAccessOptions: configData.JsonAccessOptions,
			Timeout:           configData.Timeout,
			AgwId:             configData.AgwId,
			TenantId:          configData.TenantId,
			ProfileNid:        configData.ProfileNid,
			CreateTime:        configData.CreateTime.Format(constant.LayOut),
			CreateBy:          configData.CreateBy,
			Issued:            configData.Issued,
			IssueTime:         configData.IssueTime.Format(constant.LayOut),
		},
		Point: make([]types.NeDeviceDataAccessItem, 0),
	}

	filter := model.NeDeviceDataAccessItem{ConfigNid: configData.Nid}
	page := model.NewPageInfo(0, 0)
	itemList, err := l.svcCtx.AccessItem.CommonFilterFind(l.ctx, filter, page)
	if err != nil {
		log.Println(err)
		return data
	}

	for _, i := range itemList {
		data.Point = append(data.Point, types.NeDeviceDataAccessItem{
			Nid:        i.Nid,
			ConfigNid:  i.ConfigNid,
			DeviceNid:  i.DeviceNid,
			AgwId:      i.AgwId,
			TenantId:   i.TenantId,
			ConfigType: i.ConfigType,
			AccessData: i.AccessData,
			CreateTime: i.CreateTime.Format(constant.LayOut),
		})
	}

	return data
}
