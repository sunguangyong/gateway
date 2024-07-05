package logic

import (
	"context"

	"xunjikeji.com.cn/gateway/app/external/model"

	"xunjikeji.com.cn/gateway/app/external/api/internal/svc"
	"xunjikeji.com.cn/gateway/app/external/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GatewayDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGatewayDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GatewayDeleteLogic {
	return &GatewayDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GatewayDeleteLogic) GatewayDelete(req *types.EdgeDeviceDeleteRequest) (resp *types.EdgeDeviceDeleteResponse, err error) {
	resp = &types.EdgeDeviceDeleteResponse{}
	for _, i := range req.EdgeDeviceIds {
		var request types.NeDeviceDeleteRequest
		filter := model.NeDevice{
			AgwId: i,
		}
		deviceArray, err := l.svcCtx.NeDevice.CommonFilterFind(l.ctx, filter, model.NewPageInfo(0, 0))

		if err != nil {
			continue
		}

		for _, device := range deviceArray {
			request.NeDeviceIds = append(request.NeDeviceIds, device.Nid)
		}

		newNeDeviceDelete := NewNeDeviceDeleteLogic(l.ctx, l.svcCtx)
		newNeDeviceDelete.NeDeviceDelete(&request)
		request.AgwId = i
		l.svcCtx.EdgeDevice.Delete(l.ctx, i)

	}
	return
}
