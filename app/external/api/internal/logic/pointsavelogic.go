package logic

import (
	"context"
	"encoding/json"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"xunjikeji.com.cn/gateway/common/xerr"

	"github.com/mitchellh/mapstructure"
	"xunjikeji.com.cn/gateway/app/external/api/internal/svc"
	"xunjikeji.com.cn/gateway/app/external/api/internal/types"
	"xunjikeji.com.cn/gateway/app/external/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type PointSaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPointSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PointSaveLogic {
	return &PointSaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PointSaveLogic) PointSave(req *types.PointSaveRequest) (resp *types.PointSaveResponse, err error) {
	resp = &types.PointSaveResponse{}

	var (
		jsonOptions     interface{}
		protocolType    string
		jsonAccessItems string
	)

	proto, err := l.svcCtx.AccessConfig.FindOne(l.ctx, req.ConfigNid)

	if err != nil {
		if err.Error() == sqlx.ErrNotFound.Error() {
			return resp, xerr.NotFoundNidErr
		}
		return resp, err
	}

	protocolType = proto.Protocol

	jsonOptions, ok := types.GetItemType(protocolType)

	if !ok {
		return
	}

	err = mapstructure.Decode(req.AccessOptions, &jsonOptions)
	if err != nil {
		return
	}

	data, err := json.Marshal(&jsonOptions)

	if err != nil {
		return
	}

	jsonAccessItems = string(data)
	now := time.Now()

	point := &model.NeDeviceDataAccessItem{
		ConfigNid:  req.ConfigNid,
		DeviceNid:  proto.DeviceNid,
		AgwId:      proto.AgwId,
		TenantId:   proto.TenantId,
		ConfigType: req.ConfigType,
		AccessData: jsonAccessItems,
		CreateTime: now,
	}

	if req.Nid != 0 {

		_, err = l.svcCtx.AccessItem.FindOne(l.ctx, req.Nid)

		if err != nil {
			if err.Error() == sqlx.ErrNotFound.Error() {
				return resp, xerr.NotFoundConfigErr
			}
			return resp, err
		}

		point.Nid = req.Nid
		err = l.svcCtx.AccessItem.Update(l.ctx, point)
	} else {
		_, err = l.svcCtx.AccessItem.Insert(l.ctx, point)
	}
	return
}
