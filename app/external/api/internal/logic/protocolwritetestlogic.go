package logic

import (
	"context"
	"errors"

	"xunjikeji.com.cn/gateway/common/client"

	"xunjikeji.com.cn/gateway/common/constant"
	"xunjikeji.com.cn/gateway/common/util"
	"xunjikeji.com.cn/gateway/common/xerr"

	"xunjikeji.com.cn/gateway/app/external/api/internal/svc"
	"xunjikeji.com.cn/gateway/app/external/api/internal/types"

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

func (l *ProtocolWriteTestLogic) ProtocolWriteTest(req *types.ProtocolWriteTestRequest) (resp *types.ProtocolTestResponse,
	err error) {
	host, err := l.svcCtx.EdgeDevice.GetHost(l.ctx, req.AgwId)

	if host == "" {
		return resp, xerr.NotFoundFeEdgeDeviceIpErr
	} else {
		host = util.GetInteriorHost(host, constant.InteriorPort)
	}

	body := WriteResponse{}

	err = client.PostJson(&req, &body, host, constant.ProtoWriteTest)

	if err != nil {
		return resp, xerr.NewErr(xerr.ServerCommonErrorCode, err)
	}

	if body.Code != 200 {
		return resp, xerr.NewErr(body.Code, errors.New(body.Msg))
	}

	return
}

type WriteResponse struct {
	Code uint32 `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
	} `json:"data"`
}
