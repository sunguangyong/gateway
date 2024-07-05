package logic

import (
	"context"

	"xunjikeji.com.cn/gateway/app/interior/api/internal/protocols"

	"xunjikeji.com.cn/gateway/app/interior/api/internal/svc"
	"xunjikeji.com.cn/gateway/app/interior/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProtocolStopLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProtocolStopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProtocolStopLogic {
	return &ProtocolStopLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProtocolStopLogic) ProtocolStop(req *types.ProtoStopRequest) (resp *types.ProtoStopResponse, err error) {
	resp = &types.ProtoStopResponse{}
	for _, i := range req.Nids {
		ok := protocols.ProtocolStop(i)
		if ok {
			resp.SucceedNids = append(resp.SucceedNids, i)
		} else {
			resp.FailNids = append(resp.FailNids, i)
		}
	}
	return
}
