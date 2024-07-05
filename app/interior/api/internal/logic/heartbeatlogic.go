package logic

import (
	"context"

	"xunjikeji.com.cn/gateway/app/interior/api/internal/svc"
	"xunjikeji.com.cn/gateway/app/interior/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HeartBeatLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHeartBeatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HeartBeatLogic {
	return &HeartBeatLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HeartBeatLogic) HeartBeat(req *types.HeartBeatRequest) (resp *types.HeartBeatRequestResponse, err error) {
	resp = &types.HeartBeatRequestResponse{}
	return
}
