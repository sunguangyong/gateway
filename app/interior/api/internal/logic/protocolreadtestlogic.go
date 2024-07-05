package logic

import (
	"context"

	"xunjikeji.com.cn/gateway/common/xerr"

	"xunjikeji.com.cn/gateway/common/constant"

	"xunjikeji.com.cn/gateway/app/interior/api/internal/protocols"

	"xunjikeji.com.cn/gateway/app/interior/api/internal/svc"
	"xunjikeji.com.cn/gateway/app/interior/api/internal/types"
	"xunjikeji.com.cn/gateway/app/interior/model"

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

	if len(req.NeDeviceIds) == 0 {
		return
	}

	filter := model.NeDeviceDataAccessConfig{
		DeviceNid:  req.NeDeviceIds[0],
		ConfigType: constant.ConfigTypeRead,
	}

	configArray, err := l.svcCtx.SqliteAccessConfig.Find(filter)

	if len(configArray) == 0 {
		return
	}

	if err != nil {
		if err != nil {
			err = xerr.NewErr(xerr.ServerUnexpectedlyErrorCode, err)
		}
		return
	}

	data := protocols.ProtocolReadTest(configArray[0].Nid)

	resp.DeviceName = ""
	resp.Timestamp = data.Timestamp

	for _, b := range data.ContentList {
		resp.ContentList = append(resp.ContentList, types.Content{
			Addr:      b.Addr,
			AddrValue: b.AddrValue,
		})
	}
	return
}
