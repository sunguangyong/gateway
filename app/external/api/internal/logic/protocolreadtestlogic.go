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

type ProtocolReadTestLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProtocolReadTestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProtocolReadTestLogic {
	return &ProtocolReadTestLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProtocolReadTestLogic) ProtocolReadTest(req *types.ProtocolTestRequest) (resp *types.ProtocolTestResponse, err error) {
	resp = &types.ProtocolTestResponse{}
	host, err := l.svcCtx.EdgeDevice.GetHost(l.ctx, req.AgwId)

	if host == "" {
		return resp, xerr.NotFoundFeEdgeDeviceIpErr
	} else {
		host = util.GetInteriorHost(host, constant.InteriorPort)
	}

	body := GateWayProtoTest{}

	err = client.PostJson(&req, &body, host, constant.ProtoReadTest)

	if err != nil {
		return resp, xerr.NewErr(xerr.ServerCommonErrorCode, err)
	}

	if body.Code != 200 {
		return resp, xerr.NewErr(body.Code, errors.New(body.Msg))
	}

	device, _ := l.svcCtx.NeDevice.FindOne(l.ctx, req.NeDeviceIds[0])

	if device != nil {
		resp.DeviceName = device.DeviceName
	}

	resp.Timestamp = body.Data.Timestamp

	for _, b := range body.Data.ContentList {
		resp.ContentList = append(resp.ContentList, types.Content{
			Addr:      b.Addr,
			AddrValue: b.AddrValue,
		})
	}

	return
}

type GateWayProtoTest struct {
	Code uint32 `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		DeviceName  string `json:"deviceName"`
		Timestamp   string `json:"timestamp"`
		ContentList []struct {
			Addr      string `json:"addr"`
			AddrValue int    `json:"addrValue"`
		} `json:"contentList"`
	} `json:"data"`
}
