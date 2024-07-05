package logic

import (
	"context"
	"encoding/json"
	"fmt"

	"xunjikeji.com.cn/gateway/app/external/model"

	"xunjikeji.com.cn/gateway/app/external/api/internal/svc"
	"xunjikeji.com.cn/gateway/app/external/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PointListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPointListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PointListLogic {
	return &PointListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PointListLogic) PointList(req *types.PointListRequest) (resp *types.PointListResponse, err error) {
	resp = &types.PointListResponse{
		DataList: make([]types.PointData, 0),
	}

	configNid := req.ConfigNid
	querySql := fmt.Sprintf("where config_nid = %d", configNid)

	count, err := l.svcCtx.AccessItem.FindNewsCount(l.ctx, querySql)

	if err != nil {
		return

	}

	if count == 0 {
		return

	}

	resp.Count = count

	filter := model.NeDeviceDataAccessItem{
		ConfigNid: configNid,
	}

	page := model.NewPageInfo(req.PageIndex, req.PageSize, model.NewOrder(model.Nid, model.ASC))

	protoList, err := l.svcCtx.AccessItem.CommonFilterFind(l.ctx, filter, page)

	for _, v := range protoList {

		var accessItems map[string]interface{}

		err = json.Unmarshal([]byte(v.AccessData), &accessItems)

		if err != nil {
			return
		}

		resp.DataList = append(resp.DataList, types.PointData{
			Nid:           v.Nid,
			ConfigType:    v.ConfigType,
			ConfigNid:     v.ConfigNid,
			AccessOptions: accessItems,
		})
	}
	return
}
