package logic

import (
	"context"
	"time"

	"xunjikeji.com.cn/gateway/common/xerr"

	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"xunjikeji.com.cn/gateway/app/external/model"

	"xunjikeji.com.cn/gateway/app/external/api/internal/svc"
	"xunjikeji.com.cn/gateway/app/external/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GatewaySaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGatewaySaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GatewaySaveLogic {
	return &GatewaySaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GatewaySaveLogic) GatewaySave(req *types.EdgeDeviceSaveRequest) (resp *types.EdgeDeviceSaveResponse, err error) {
	resp = &types.EdgeDeviceSaveResponse{}

	edgeDevice := &model.FeEdgeDevice{
		EdgeDeviceName: req.EdgeDeviceName,
		DevicePosition: req.DevicePosition,
		IpAddress:      req.IpAddress,
		Mac:            req.Mac,
		DeviceDesc:     req.DeviceDesc,
		SimCard:        req.SimCard,
		GatewayId:      req.GatewayId,
		UpdateTime:     time.Now(),
		CreateTime:     time.Now(),
	}

	if req.EdgeDeviceId > 0 { // 修改

		_, err = l.svcCtx.EdgeDevice.FindOne(l.ctx, req.EdgeDeviceId)

		if err != nil {
			if err.Error() == sqlx.ErrNotFound.Error() {
				return resp, xerr.NotFoundFeEdgeDeviceErr
			}
			return resp, err
		}

		edgeDevice.EdgeDeviceId = req.EdgeDeviceId
		err = l.svcCtx.EdgeDevice.Update(l.ctx, edgeDevice)
		err = l.svcCtx.EdgeDevice.ConvertUniqIdxMacErr(err)

	} else {
		_, err = l.svcCtx.EdgeDevice.Insert(l.ctx, edgeDevice)
		err = l.svcCtx.EdgeDevice.ConvertUniqIdxMacErr(err)
		return resp, err
	}
	return
}
