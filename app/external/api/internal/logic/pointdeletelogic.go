package logic

import (
	"context"

	"xunjikeji.com.cn/gateway/app/external/api/internal/svc"
	"xunjikeji.com.cn/gateway/app/external/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PointDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPointDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PointDeleteLogic {
	return &PointDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PointDeleteLogic) PointDelete(req *types.PointDelRequest) (resp *types.PointDelResponse, err error) {

	resp = &types.PointDelResponse{}

	for _, i := range req.Nids {
		l.svcCtx.AccessItem.Delete(l.ctx, i)
	}
	return
}
