package logic

import (
	"context"

	"xunjikeji.com.cn/gateway/app/external/api/internal/svc"
	"xunjikeji.com.cn/gateway/app/external/api/internal/types"
	"xunjikeji.com.cn/gateway/app/external/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProtoDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProtoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProtoDeleteLogic {
	return &ProtoDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProtoDeleteLogic) ProtoDelete(req *types.ProtocolDelRequest) (resp *types.ProtocolDelResponse, err error) {
	resp = &types.ProtocolDelResponse{}
	for _, i := range req.Nids {
		var request types.PointDelRequest
		filter := model.NeDeviceDataAccessItem{
			ConfigNid: i,
		}

		items, err := l.svcCtx.AccessItem.CommonFilterFind(l.ctx, filter, model.PageInfo{Size: 0, Page: 0})

		if err != nil {
			continue
		}

		for _, item := range items {
			request.Nids = append(request.Nids, item.Nid)
		}

		pointDelete := NewPointDeleteLogic(l.ctx, l.svcCtx)
		// 删除点表
		pointDelete.PointDelete(&request)
		// 删除协议
		l.svcCtx.AccessConfig.Delete(l.ctx, i)
	}

	return
}
