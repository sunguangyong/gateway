package logic

import (
	"context"

	"xunjikeji.com.cn/gateway/app/external/model"

	"xunjikeji.com.cn/gateway/app/external/api/internal/svc"
	"xunjikeji.com.cn/gateway/app/external/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PointDeleteAllLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPointDeleteAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PointDeleteAllLogic {
	return &PointDeleteAllLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PointDeleteAllLogic) PointDeleteAll(req *types.PointDelAllRequest) (resp *types.PointDelAllResponse, err error) {

	resp = &types.PointDelAllResponse{}
	configNid := req.ConfigNid

	filter := model.NeDeviceDataAccessItem{
		ConfigNid: configNid,
	}

	page := model.NewPageInfo(0, 0)

	protoList, err := l.svcCtx.AccessItem.CommonFilterFind(l.ctx, filter, page)

	for _, v := range protoList {
		l.svcCtx.AccessItem.Delete(l.ctx, v.Nid)
	}

	return
}
