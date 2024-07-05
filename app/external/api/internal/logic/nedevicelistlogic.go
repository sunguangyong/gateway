package logic

import (
	"context"

	"xunjikeji.com.cn/gateway/app/external/model"
	"xunjikeji.com.cn/gateway/common/constant"
	"xunjikeji.com.cn/gateway/common/convert"

	"xunjikeji.com.cn/gateway/app/external/api/internal/svc"
	"xunjikeji.com.cn/gateway/app/external/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type NeDeviceListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewNeDeviceListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NeDeviceListLogic {
	return &NeDeviceListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *NeDeviceListLogic) NeDeviceList(req *types.NeDeviceListRequest) (resp *types.NeDeviceListResponse, err error) {
	resp = &types.NeDeviceListResponse{
		Data:  make([]types.NeDeviceList, 0),
		Count: 0,
	}
	var filter model.NeDevice
	filter.AgwId = req.EdgeDeviceId
	count, err := l.svcCtx.NeDevice.CommonFilterCount(l.ctx, filter)

	if err != nil {
		return
	}

	if count == 0 {
		return
	}

	page := model.NewPageInfo(req.PageIndex, req.PageSize, model.NewOrder(model.Nid, model.DESC))
	dataList, err := l.svcCtx.NeDevice.CommonFilterFind(l.ctx, filter, page)

	if err != nil {
		return
	}

	for _, v := range dataList {
		var data types.NeDeviceList
		convert.CopyProperties(&data, v)
		data.CreateTime = v.CreateTime.Format(constant.LayOut)
		resp.Data = append(resp.Data, data)
	}

	resp.Count = count

	return
}
