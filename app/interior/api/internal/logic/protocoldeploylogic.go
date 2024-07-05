package logic

import (
	"context"
	"log"
	"time"

	"xunjikeji.com.cn/gateway/app/interior/api/internal/protocols"

	"xunjikeji.com.cn/gateway/common/constant"

	"gorm.io/gorm"
	"xunjikeji.com.cn/gateway/app/interior/model"

	"xunjikeji.com.cn/gateway/app/interior/api/internal/svc"
	"xunjikeji.com.cn/gateway/app/interior/api/internal/types"

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

func (l *ProtocolDeployLogic) ProtocolDeploy(req *types.ProtoDeployRequest) (*types.ProtoDeployResponse, error) {
	resp := &types.ProtoDeployResponse{
		FailArray:    make([]types.DeployMsg, 0),
		SuccessArray: make([]types.DeployMsg, 0),
	}

	for _, device := range req.DataList {
		if len(device.ConfigData) == 0 {
			l.handleEmptyConfigData(device)
		} else {
			l.handleConfigData(device, resp)
		}
	}

	return resp, nil
}

// 只有设备信息没有协议(服务器已删除协议)
func (l *ProtocolDeployLogic) handleEmptyConfigData(device types.DeployDevice) {
	configFilter := model.NeDeviceDataAccessConfig{DeviceNid: device.DeviceNid}

	configArray, err := l.svcCtx.SqliteAccessConfig.Find(configFilter)
	if err != nil {
		log.Println(err)
	}

	for _, config := range configArray {
		protocols.ProtocolStop(config.Nid)
	}

	l.svcCtx.SqliteAccessConfig.Delete(configFilter)
	pointFilter := model.NeDeviceDataAccessItem{DeviceNid: device.DeviceNid}
	l.svcCtx.SqliteAccessItem.Delete(pointFilter)
}

// 正常部署
func (l *ProtocolDeployLogic) handleConfigData(device types.DeployDevice, resp *types.ProtoDeployResponse) {
	for _, c := range device.ConfigData {
		err := l.svcCtx.SqliteAccessConfig.Transaction(func(tx *gorm.DB) error {
			if err := l.deleteAndInsertConfig(tx, c); err != nil {
				return err
			}
			return l.deleteAndInsertPoints(tx, c.Point)
		})

		if err != nil {
			resp.FailArray = append(resp.FailArray, types.DeployMsg{
				Msg:        err.Error(),
				ConfigData: c.Config,
			})
			continue
		}

		ok, err := protocols.ProtocolFlow(c)
		if ok && err == nil {
			resp.SuccessArray = append(resp.SuccessArray, types.DeployMsg{ConfigData: c.Config})
		} else {
			resp.FailArray = append(resp.FailArray, types.DeployMsg{
				Msg:        err.Error(),
				ConfigData: c.Config,
			})
		}
	}
}

// 删除旧协议信息,存入新协议信息
func (l *ProtocolDeployLogic) deleteAndInsertConfig(tx *gorm.DB, c types.ConfigData) error {
	configs := model.NeDeviceDataAccessConfig{Nid: c.Config.Nid}
	// 删除旧配置信息
	if err := tx.Where(configs).Delete(configs).Error; err != nil {
		return err
	}

	accessConfig := l.convertToAccessConfig(c.Config)
	// 插入新的配置信息
	return tx.Create(&accessConfig).Error
}

// 删除旧点位信息,存入新点位信息
func (l *ProtocolDeployLogic) deleteAndInsertPoints(tx *gorm.DB, points []types.NeDeviceDataAccessItem) error {
	for _, p := range points {
		point := model.NeDeviceDataAccessItem{Nid: p.Nid}
		// 删除旧点位
		if err := tx.Where(point).Delete(point).Error; err != nil {
			return err
		}

		// 插入新点位
		itemConfig := l.convertToItemConfig(p)
		if err := tx.Create(&itemConfig).Error; err != nil {
			return err
		}
	}
	return nil
}

func (l *ProtocolDeployLogic) convertToAccessConfig(config types.NeDeviceDataAccessConfig) (accessConfig *model.
	NeDeviceDataAccessConfig) {

	accessConfig = &model.NeDeviceDataAccessConfig{
		Nid:               config.Nid,
		DeviceNid:         config.DeviceNid,
		ConfigType:        config.ConfigType,
		ConfigId:          config.ConfigId,
		ConfigName:        config.ConfigName,
		Endpoint:          config.Endpoint,
		Protocol:          config.Protocol,
		JsonAccessOptions: config.JsonAccessOptions,
		Timeout:           config.Timeout,
		AgwId:             config.AgwId,
		TenantId:          config.TenantId,
		ProfileNid:        config.ProfileNid,
		CreateTime:        config.CreateTime,
		CreateBy:          config.CreateBy,
		Issued:            constant.IssuedOn,
		IssueTime:         time.Now().Format(constant.LayOut),
	}
	return accessConfig
}

func (l *ProtocolDeployLogic) convertToItemConfig(item types.NeDeviceDataAccessItem) (accessItem *model.
	NeDeviceDataAccessItem) {
	accessItem = &model.NeDeviceDataAccessItem{
		Nid:        item.Nid,
		ConfigType: int(item.ConfigType),
		AgwId:      item.AgwId,
		TenantId:   item.TenantId,
		DeviceNid:  item.DeviceNid,
		CreateTime: item.CreateTime,
		AccessData: item.AccessData,
		ConfigNid:  item.ConfigNid,
	}
	return accessItem
}
