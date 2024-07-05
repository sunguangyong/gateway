package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"xunjikeji.com.cn/gateway/app/external/api/internal/config"
	"xunjikeji.com.cn/gateway/app/external/model"
)

type ServiceContext struct {
	Config       config.Config
	AccessConfig model.NeDeviceDataAccessConfigModel
	AccessItem   model.NeDeviceDataAccessItemModel
	EdgeDevice   model.FeEdgeDeviceModel
	NeDevice     model.NeDeviceModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:       c,
		AccessConfig: model.NewNeDeviceDataAccessConfigModel(sqlx.NewMysql(c.IotMysql.DataSource)),
		AccessItem:   model.NewNeDeviceDataAccessItemModel(sqlx.NewMysql(c.IotMysql.DataSource)),
		EdgeDevice:   model.NewFeEdgeDeviceModel(sqlx.NewMysql(c.IotMysql.DataSource)),
		NeDevice:     model.NewNeDeviceModel(sqlx.NewMysql(c.IotMysql.DataSource)),
	}
}
