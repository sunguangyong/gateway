package logic

import (
	"context"

	"xunjikeji.com.cn/gateway/app/external/api/internal/svc"
	"xunjikeji.com.cn/gateway/app/external/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GatewayListExportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGatewayListExportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GatewayListExportLogic {
	return &GatewayListExportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GatewayListExportLogic) GatewayListExport(req *types.EdgeDeviceListRequest) (resp *types.EdgeDeviceListRequest, err error) {
	// todo: add your logic here and delete this line

	return
}
