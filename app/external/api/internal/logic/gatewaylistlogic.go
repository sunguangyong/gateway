package logic

import (
	"context"

	"xunjikeji.com.cn/gateway/app/external/model"
	"xunjikeji.com.cn/gateway/common/convert"

	"xunjikeji.com.cn/gateway/app/external/api/internal/svc"
	"xunjikeji.com.cn/gateway/app/external/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GatewayListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGatewayListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GatewayListLogic {
	return &GatewayListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GatewayListLogic) GatewayList(req *types.EdgeDeviceListRequest) (resp *types.EdgeDeviceListResponse, err error) {
	resp = &types.EdgeDeviceListResponse{
		Data: make([]types.EdgeDeviceData, 0),
	}

	var filter model.FeEdgeDevice

	count, err := l.svcCtx.EdgeDevice.CommonFilterCount(l.ctx, filter)

	if err != nil {
		return
	}

	if count == 0 {
		return
	}

	pageInfo := model.NewPageInfo(req.PageIndex, req.PageSize, model.NewOrder(model.EdgeDeviceId, model.DESC))

	dataArray, err := l.svcCtx.EdgeDevice.CommonFilterFind(l.ctx, filter, pageInfo)
	resp.Count = count

	var device types.EdgeDeviceData

	for _, v := range dataArray {
		convert.CopyProperties(&device, v)
		resp.Data = append(resp.Data, device)
	}

	return
}
