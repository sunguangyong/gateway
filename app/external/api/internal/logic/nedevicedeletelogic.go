package logic

import (
	"context"

	"xunjikeji.com.cn/gateway/app/external/api/internal/svc"
	"xunjikeji.com.cn/gateway/app/external/api/internal/types"
	"xunjikeji.com.cn/gateway/app/external/model"
	"xunjikeji.com.cn/gateway/common/client"
	"xunjikeji.com.cn/gateway/common/constant"
	"xunjikeji.com.cn/gateway/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type NeDeviceDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewNeDeviceDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NeDeviceDeleteLogic {
	return &NeDeviceDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *NeDeviceDeleteLogic) NeDeviceDelete(req *types.NeDeviceDeleteRequest) (resp *types.NeDeviceDeleteResponse, err error) {
	resp = &types.NeDeviceDeleteResponse{}

	var nids []int64

	host, _ := l.svcCtx.EdgeDevice.GetHost(l.ctx, req.AgwId)

	for _, i := range req.NeDeviceIds {

		filter := model.NeDeviceDataAccessConfig{
			DeviceNid: i,
		}

		page := model.NewPageInfo(0, 0)

		configArray, _ := l.svcCtx.AccessConfig.CommonFilterFind(l.ctx, filter, page)

		for _, n := range configArray {
			nids = append(nids, n.Nid)
		}

		// 删除设备
		l.svcCtx.NeDevice.Delete(l.ctx, i)
	}

	if len(nids) == 0 {
		return
	}

	// 删除服务器协议
	webRequest := types.ProtocolDelRequest{
		Nids: nids,
	}

	protocolDel := NewProtoDeleteLogic(l.ctx, l.svcCtx)
	protocolDel.ProtoDelete(&webRequest)

	// 删除网关协议
	request := types.ProtoDeleteRequest{
		Nids: nids,
	}

	body := types.ProtoDeleteResponse{}

	if host == "" {
		return
	}
	host = util.GetInteriorHost(host, constant.InteriorPort)
	err = client.PostJson(request, &body, host, constant.ProtoDeployDelete)
	return
}
