package svc

import (
	"xunjikeji.com.cn/gateway/app/interior/api/internal/config"
	"xunjikeji.com.cn/gateway/app/interior/model"
)

type ServiceContext struct {
	Config             config.Config
	SqliteAccessConfig model.SqliteNeDeviceDataAccessConfigModel
	SqliteAccessItem   model.SqliteNeDeviceDataAccessItemModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:             c,
		SqliteAccessConfig: model.NewNeDeviceDataAccessConfigModel(model.NewSqliteConn(c.Sqlite.SqliteDsn)),
		SqliteAccessItem:   model.NewNeDeviceDataAccessItemModel(model.NewSqliteConn(c.Sqlite.SqliteDsn)),
	}
}
