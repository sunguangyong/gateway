package logic

import (
	"context"

	"xunjikeji.com.cn/gateway/common/constant"

	"xunjikeji.com.cn/gateway/app/external/api/internal/svc"
	"xunjikeji.com.cn/gateway/app/external/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProtoDropDownLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProtoDropDownLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProtoDropDownLogic {
	return &ProtoDropDownLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProtoDropDownLogic) ProtoDropDown(req *types.ProtoDropDownRequest) (resp *types.ProtoDropDownResponse, err error) {
	resp = &types.ProtoDropDownResponse{
		List: make([]types.Dropdown, 0),
	}

	dataList := make([]types.Dropdown, 0)

	for _, v := range constant.ProtocolNameMap {
		data := types.Dropdown{
			Value: v,
			Label: v,
		}
		dataList = append(dataList, data)
	}

	sortDropdown(dataList)
	resp.List = dataList
	return
}
