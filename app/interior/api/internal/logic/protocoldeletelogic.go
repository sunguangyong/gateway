package logic

import (
	"context"

	"gorm.io/gorm"
	"xunjikeji.com.cn/gateway/app/interior/api/internal/protocols"
	"xunjikeji.com.cn/gateway/app/interior/model"

	"xunjikeji.com.cn/gateway/app/interior/api/internal/svc"
	"xunjikeji.com.cn/gateway/app/interior/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProtocolDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProtocolDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProtocolDeleteLogic {
	return &ProtocolDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProtocolDeleteLogic) ProtocolDelete(req *types.ProtoDeleteRequest) (resp *types.ProtoDeleteResponse,
	err error) {
	resp = &types.ProtoDeleteResponse{
		SucceedNids: make([]int64, 0),
		FailNids:    make([]int64, 0),
	}
	for _, i := range req.Nids {
		ok := protocols.ProtocolStop(i)
		if ok {
			// 使用事务删除
			err = l.svcCtx.SqliteAccessConfig.Transaction(func(tx *gorm.DB) error {

				defer func() {
					if err != nil {
						resp.FailNids = append(resp.FailNids, i)
					}
				}()

				accessConfig := model.NeDeviceDataAccessConfig{
					Nid: i,
				}

				err = tx.Where(accessConfig).Delete(accessConfig).Error

				if err != nil {
					return err
				}

				itemConfig := model.NeDeviceDataAccessItem{ConfigNid: i}
				err = tx.Where(itemConfig).Delete(itemConfig).Error

				if err != nil {
					resp.FailNids = append(resp.FailNids, i)
					return err
				}
				resp.SucceedNids = append(resp.SucceedNids, i)
				return nil
			})
		}
	}

	return
}
