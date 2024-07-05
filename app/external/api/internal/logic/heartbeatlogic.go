package logic

import (
	"context"

	"xunjikeji.com.cn/gateway/app/external/model"
	"xunjikeji.com.cn/gateway/common/client"
	"xunjikeji.com.cn/gateway/common/constant"
	"xunjikeji.com.cn/gateway/common/util"

	"xunjikeji.com.cn/gateway/app/external/api/internal/svc"
	"xunjikeji.com.cn/gateway/app/external/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HeartbeatLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHeartbeatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HeartbeatLogic {
	return &HeartbeatLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HeartbeatLogic) Heartbeat(req *types.HeartBeatRequest) (resp *types.HeartBeatRequestResponse, err error) {
	resp = &types.HeartBeatRequestResponse{}
	var host string

	var filter model.FeEdgeDevice
	filter.EdgeDeviceId = req.EdgeDeviceId
	pageInfo := model.NewPageInfo(1, 1)
	feEdgeDeviceList, err := l.svcCtx.EdgeDevice.CommonFilterFind(l.ctx, filter, pageInfo)

	if err != nil || len(feEdgeDeviceList) == 0 {
		return resp, err
	}

	host = feEdgeDeviceList[0].IpAddress

	if host == "" {
		return resp, err
	} else {
		host = util.GetInteriorHost(host, constant.InteriorPort)
	}

	request := types.InteriorHeartBeatRequest{}
	body := types.InteriorHeartBeatResponse{}

	err = client.PostJson(request, &body, host, constant.HeartBeatPath)

	if err != nil {
		return resp, err
	}

	return
}
