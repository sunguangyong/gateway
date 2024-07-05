package logic

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"xunjikeji.com.cn/gateway/common/xerr"

	"xunjikeji.com.cn/gateway/app/external/api/internal/svc"
	"xunjikeji.com.cn/gateway/app/external/api/internal/types"
	"xunjikeji.com.cn/gateway/app/external/model"
	"xunjikeji.com.cn/gateway/common/convert"

	"github.com/zeromicro/go-zero/core/logx"
)

type NeDeviceSaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewNeDeviceSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NeDeviceSaveLogic {
	return &NeDeviceSaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *NeDeviceSaveLogic) NeDeviceSave(req *types.NeDeviceSaveRequest) (resp *types.NeDeviceSaveResponse, err error) {

	resp = &types.NeDeviceSaveResponse{}

	var neDevice model.NeDevice
	convert.CopyProperties(&neDevice, req)

	_, err = l.svcCtx.EdgeDevice.FindOne(l.ctx, req.AgwId)

	if err != nil {
		if err.Error() == sqlx.ErrNotFound.Error() {
			return resp, xerr.NotFoundFeEdgeDeviceErr
		}
		return resp, err
	}

	if req.Nid > 0 { // 更新
		_, err = l.svcCtx.NeDevice.FindOne(l.ctx, req.Nid)

		if err != nil {
			if err.Error() == sqlx.ErrNotFound.Error() {
				return resp, xerr.NotFoundNeDeviceErr
			}
			return resp, err
		}

		l.svcCtx.NeDevice.Update(l.ctx, &neDevice)

	} else { // 新增
		neDevice.CreateTime = time.Now()
		l.svcCtx.NeDevice.Insert(l.ctx, &neDevice)
	}

	return resp, err
}
