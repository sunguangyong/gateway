package logic

import (
	"context"
	"encoding/json"

	"github.com/zeromicro/go-zero/core/logx"
	"xunjikeji.com.cn/gateway/app/external/api/internal/svc"
	"xunjikeji.com.cn/gateway/app/external/api/internal/types"
	"xunjikeji.com.cn/gateway/app/external/model"
)

type PointWriteListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPointWriteListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PointWriteListLogic {
	return &PointWriteListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PointWriteListLogic) PointWriteList(req *types.WritePointListRequest) (resp *types.WritepointListResponse, err error) {

	resp = &types.WritepointListResponse{
		Data: make([]types.WritepointList, 0),
	}

	page := model.NewPageInfo(1, 0)
	var filter model.NeDeviceDataAccessItem
	filter.DeviceNid = req.NeDeviceId
	filter.ConfigType = req.ConfigType

	itemArray, err := l.svcCtx.AccessItem.CommonFilterFind(l.ctx, filter, page)
	if err != nil {
		return
	}

	if len(itemArray) == 0 {
		return
	}

	for _, item := range itemArray {
		var jsonOptions AccessDataCommon
		err = json.Unmarshal([]byte(item.AccessData), &jsonOptions)
		if err != nil {
			return
		}

		resp.Data = append(resp.Data, types.WritepointList{
			Nid:       item.Nid,
			Address:   jsonOptions.Address,
			DataName:  jsonOptions.DataName,
			Method:    jsonOptions.Method,
			ConfigNid: item.ConfigNid,
		})
	}
	return
}

type AccessDataCommon struct {
	DataName string `json:"dataName"` // 数据名称
	Method   string `json:"method"`   // 读取方法
	Address  uint16 `json:"address"`  // 地址位
}
