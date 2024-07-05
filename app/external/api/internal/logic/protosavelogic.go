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

type ProtoSaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProtoSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProtoSaveLogic {
	return &ProtoSaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProtoSaveLogic) ProtoSave(req *types.ProtocolSaveRequest) (resp *types.ProtocolSaveResponse, err error) {
	var (
		jsonOptions       interface{}
		jsonAccessOptions string
	)

	resp = &types.ProtocolSaveResponse{}

	jsonOptions, ok := types.GetJsonOptionsType(req.Protocol)

	if !ok {
		return resp, xerr.NotSupportedProtocolError
	}

	err = mapstructure.Decode(req.AccessOptions, &jsonOptions)

	if err != nil {
		err = xerr.ConvertAccessOptionsParamsErr(err)
		return
	}

	data, err := json.Marshal(&jsonOptions)
	if err != nil {
		return
	}

	jsonAccessOptions = string(data)
	now := time.Now()

	neDeviceDataAccessConfig := &model.NeDeviceDataAccessConfig{
		DeviceNid:         req.DeviceNid,
		ConfigType:        req.ConfigType,
		ConfigName:        req.ConfigName,
		Endpoint:          req.Endpoint,
		Protocol:          req.Protocol,
		JsonAccessOptions: jsonAccessOptions,
		AgwId:             req.AgwId,
		Issued:            0,
	}

	_, err = l.svcCtx.EdgeDevice.FindOne(l.ctx, req.AgwId)
	if err != nil {
		return resp, xerr.NotFoundFeEdgeDeviceErr
	}

	_, err = l.svcCtx.NeDevice.FindOne(l.ctx, req.DeviceNid)

	if err != nil {
		return resp, xerr.NotFoundNeDeviceErr
	}

	if req.Nid != 0 {

		_, err = l.svcCtx.AccessConfig.FindOne(l.ctx, req.Nid)

		if err != nil {
			if err.Error() == sqlx.ErrNotFound.Error() {
				return resp, xerr.NotFoundNidErr
			}
			return resp, err
		}

		neDeviceDataAccessConfig.Nid = req.Nid
		err = l.svcCtx.AccessConfig.Update(l.ctx, neDeviceDataAccessConfig)
	} else {
		neDeviceDataAccessConfig.CreateTime = now
		_, err = l.svcCtx.AccessConfig.Insert(l.ctx, neDeviceDataAccessConfig)
	}

	if err != nil {
		return
	}

	return
}
