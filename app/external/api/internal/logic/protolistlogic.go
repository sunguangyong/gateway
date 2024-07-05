package logic

import (
	"context"
	"encoding/json"

	"xunjikeji.com.cn/gateway/app/external/model"

	"xunjikeji.com.cn/gateway/common/constant"

	"xunjikeji.com.cn/gateway/app/external/api/internal/svc"
	"xunjikeji.com.cn/gateway/app/external/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProtoGainLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProtoGainLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProtoGainLogic {
	return &ProtoGainLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProtoGainLogic) ProtoGain(req *types.ProtocolListRequest) (resp *types.ProtocolListResponse,
	err error) {
	resp = &types.ProtocolListResponse{}

	filter := model.NeDeviceDataAccessConfig{
		DeviceNid: req.DeviceNid,
	}
	page := model.NewPageInfo(0, 0)
	protocolList, err := l.svcCtx.AccessConfig.CommonFilterFind(l.ctx, filter, page)

	if err != nil {
		return
	}

	for _, proto := range protocolList {

		var accessOptions map[string]interface{}

		err = json.Unmarshal([]byte(proto.JsonAccessOptions), &accessOptions)

		if err != nil {
			return
		}

		var issueTime string

		if proto.Issued > 0 {
			issueTime = proto.IssueTime.Format(constant.LayOut)
		}

		resp.Data = append(resp.Data, types.ProtocolGain{
			Nid:           proto.Nid,
			DeviceNid:     proto.DeviceNid,
			ConfigType:    proto.ConfigType,
			ConfigName:    proto.ConfigName,
			Endpoint:      proto.Endpoint,
			Protocol:      proto.Protocol,
			AccessOptions: accessOptions,
			AgwId:         proto.AgwId,
			IssueTime:     issueTime,
			Issued:        proto.Issued,
		})
	}
	return
}
