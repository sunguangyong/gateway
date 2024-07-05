package logic

import (
	"context"

	"xunjikeji.com.cn/gateway/app/interior/api/internal/protocols"
	"xunjikeji.com.cn/gateway/app/interior/api/internal/svc"
	"xunjikeji.com.cn/gateway/app/interior/api/internal/types"
	"xunjikeji.com.cn/gateway/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProtocolWriteTestLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProtocolWriteTestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProtocolWriteTestLogic {
	return &ProtocolWriteTestLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProtocolWriteTestLogic) ProtocolWriteTest(req *types.ProtocolWriteTestRequest) (resp *types.ProtocolWriteTestResponse, err error) {
	resp = &types.ProtocolWriteTestResponse{}
	proto, ok := protocols.ProtocolMapGet(req.ConfigNid)

	if ok {
		err = proto.Write(req)
	}

	if err != nil {
		err = xerr.NewErr(xerr.ServerUnexpectedlyErrorCode, err)
	}

	return
}
