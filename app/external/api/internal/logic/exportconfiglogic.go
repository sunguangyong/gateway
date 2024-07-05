package logic

import (
	"context"
	"encoding/json"

	"xunjikeji.com.cn/gateway/app/external/model"

	"github.com/zeromicro/go-zero/core/logx"
	"xunjikeji.com.cn/gateway/app/external/api/internal/svc"
	"xunjikeji.com.cn/gateway/app/external/api/internal/types"
	"xunjikeji.com.cn/gateway/common/constant"
)

type ExportConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExportConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExportConfigLogic {
	return &ExportConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExportConfigLogic) ExportConfig(req *types.ExportConfigRequest) (resp *types.ProtocolSaveResponse, err error) {
	resp = &types.ProtocolSaveResponse{}
	var filter model.NeDeviceDataAccessConfig
	filter.AgwId = req.AgwId
	pageInfo := model.NewPageInfo(0, 0)
	configList, err := l.svcCtx.AccessConfig.CommonFilterFind(l.ctx, filter, pageInfo)

	if err != nil {
		return
	}

	dataList := make([]types.ConfigData, 0)
	for _, c := range configList {
		data := types.ConfigData{
			Config: types.NeDeviceDataAccessConfig{
				Nid:               c.Nid,
				DeviceNid:         c.DeviceNid,
				ConfigType:        c.ConfigType,
				ConfigId:          c.ConfigId,
				ConfigName:        c.ConfigName,
				Endpoint:          c.Endpoint,
				Protocol:          c.Protocol,
				JsonAccessOptions: c.JsonAccessOptions,
				Timeout:           c.Timeout,
				AgwId:             c.AgwId,
				TenantId:          c.TenantId,
				ProfileNid:        c.ProfileNid,
				CreateTime:        c.CreateTime.Format(constant.LayOut),
				CreateBy:          c.CreateBy,
				Issued:            c.Issued,
				IssueTime:         c.IssueTime.Format(constant.LayOut),
			},
			Point: make([]types.NeDeviceDataAccessItem, 0),
		}
		itemList, err := l.svcCtx.AccessItem.CommonFilterFind(l.ctx, model.NeDeviceDataAccessItem{ConfigNid: c.Nid},
			pageInfo)

		if err != nil {
			return resp, err
		}

		for _, i := range itemList {
			data.Point = append(data.Point, types.NeDeviceDataAccessItem{
				Nid:        i.Nid,
				ConfigNid:  i.ConfigNid,
				ConfigType: i.ConfigType,
				AccessData: i.AccessData,
				CreateTime: i.CreateTime.Format(constant.LayOut),
			})
		}

		dataList = append(dataList, data)
	}

	info, err := json.Marshal(dataList)

	if err != nil {
		return resp, err
	}

	resp.Data = string(info)

	return resp, err
}
