package logic

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"xunjikeji.com.cn/gateway/app/external/api/internal/svc"
	"xunjikeji.com.cn/gateway/app/external/api/internal/types"
	"xunjikeji.com.cn/gateway/app/external/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type DdaConfigSaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDdaConfigSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DdaConfigSaveLogic {
	return &DdaConfigSaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DdaConfigSaveLogic) DdaConfigSave(req *types.ConfigSaveRequest) (resp *types.ConfigSaveResponse, err error) {
	resp = &types.ConfigSaveResponse{}

	deviceNid, _ := strconv.ParseInt(req.DeviceNid, 10, 64)
	agwId, _ := strconv.ParseInt(req.AgwId, 10, 64)
	configType := int64(req.ConfigType)
	jsonAccessOptions := req.AccessOptions
	jsonString, err := json.Marshal(jsonAccessOptions)
	now := time.Now()

	if err != nil {
		return
	}

	neDeviceDataAccessConfig := &model.NeDeviceDataAccessConfig{
		DeviceNid:         deviceNid,
		ConfigType:        configType,
		ConfigName:        req.ConfigName,
		Endpoint:          req.Endpoint,
		Protocol:          req.Protocol,
		JsonAccessOptions: string(jsonString),
		AgwId:             agwId,
		Issued:            0,
	}

	if req.Nid != "" && req.Nid != "0" { // 修改
		nid, _ := strconv.ParseInt(req.Nid, 10, 64)
		neDeviceDataAccessConfig.Nid = nid
		err = l.svcCtx.AccessConfig.Update(l.ctx, neDeviceDataAccessConfig)
	} else {
		neDeviceDataAccessConfig.CreateTime = now
		_, err = l.svcCtx.AccessConfig.Insert(l.ctx, neDeviceDataAccessConfig)
	}
	return
}
