package protocols

import (
	"xunjikeji.com.cn/gateway/app/interior/api/internal/svc"
	"xunjikeji.com.cn/gateway/app/interior/api/internal/types"
	"xunjikeji.com.cn/gateway/app/interior/model"
	"xunjikeji.com.cn/gateway/common/constant"
)

// StartProtocols 服务开启时获取本地保存的所有配置信息
func StartProtocols(serverCtx *svc.ServiceContext) {
	var configFilter model.NeDeviceDataAccessConfig
	configFilter.Issued = constant.IssuedOn
	accessConfigs, err := serverCtx.SqliteAccessConfig.Find(configFilter)
	if err != nil {
		return
	}

	for _, accessConfig := range accessConfigs {

		var itemFilter model.NeDeviceDataAccessItem
		itemFilter.ConfigNid = accessConfig.Nid
		data := types.ConfigData{
			Config: types.NeDeviceDataAccessConfig{
				Nid:               accessConfig.Nid,
				DeviceNid:         accessConfig.DeviceNid,
				ConfigType:        accessConfig.ConfigType,
				ConfigId:          accessConfig.ConfigId,
				ConfigName:        accessConfig.ConfigName,
				Endpoint:          accessConfig.Endpoint,
				Protocol:          accessConfig.Protocol,
				JsonAccessOptions: accessConfig.JsonAccessOptions,
				Timeout:           accessConfig.Timeout,
				AgwId:             accessConfig.AgwId,
				TenantId:          accessConfig.TenantId,
				ProfileNid:        accessConfig.ProfileNid,
				CreateTime:        accessConfig.CreateTime,
				CreateBy:          accessConfig.CreateBy,
				Issued:            accessConfig.Issued,
				IssueTime:         accessConfig.IssueTime,
			},
			Point: make([]types.NeDeviceDataAccessItem, 0),
		}
		items, _ := serverCtx.SqliteAccessItem.Find(itemFilter)
		for _, item := range items {
			data.Point = append(data.Point, types.NeDeviceDataAccessItem{
				Nid:        item.Nid,
				ConfigNid:  item.ConfigNid,
				ConfigType: int64(item.ConfigType),
				AccessData: item.AccessData,
				CreateTime: item.CreateTime,
			})
		}
		ProtocolFlow(data)
	}
}
