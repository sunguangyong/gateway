package model

import (
	"database/sql"

	"gorm.io/gorm"
)

type (
	SqliteNeDeviceDataAccessItemModel interface {
		Insert(n *NeDeviceDataAccessItem) (tx *gorm.DB)
		Find(filter interface{}) (result []NeDeviceDataAccessItem, err error)
		Updates(filter interface{}, data NeDeviceDataAccessItem) (err error)
		Delete(data NeDeviceDataAccessItem) (err error)
		Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) (err error)
	}

	NeDeviceDataAccessItem struct {
		Nid        int64  `gorm:"column:nid;type:bigint(20);primary_key;AUTO_INCREMENT;comment:主键" json:"nid"`
		ConfigNid  int64  `gorm:"column:config_nid;type:bigint(20);comment:ne_device_data_access_config.id;NOT NULL" json:"config_nid"`
		DeviceNid  int64  `gorm:"column:device_nid;type:bigint(20);comment:设备ID;NOT NULL" json:"device_nid"`
		AgwId      int64  `gorm:"column:agw_id;type:bigint(20);comment:所属网关ID" json:"agw_id"`
		TenantId   int64  `gorm:"column:tenant_id;type:bigint(20);comment:租户ID" json:"tenant_id"`
		ConfigType int    `gorm:"column:config_type;type:int(11);default:0;comment:0-读/写,1-只读,2-只写;NOT NULL" json:"config_type"`
		AccessData string `gorm:"column:access_data;type:varchar(4000);comment:读/写配置项(JSON对象);NOT NULL" json:"access_data"`
		CreateTime string `gorm:"column:create_time;type:datetime;comment:创建时间;NOT NULL" json:"create_time"`
	}

	defaultNeDeviceDataAccessItemModel struct {
		conn  *gorm.DB
		table string
	}
)

func (m *NeDeviceDataAccessItem) TableName() string {
	return "ne_device_data_access_item"
}

func newNeDeviceDataAccessItemModel(conn *gorm.DB) *defaultNeDeviceDataAccessItemModel {
	return &defaultNeDeviceDataAccessItemModel{
		conn:  conn,
		table: "`ne_device_data_access_item`",
	}
}

func NewNeDeviceDataAccessItemModel(conn *gorm.DB) SqliteNeDeviceDataAccessItemModel {
	return newNeDeviceDataAccessItemModel(conn)
}

func (m *defaultNeDeviceDataAccessItemModel) Insert(n *NeDeviceDataAccessItem) (tx *gorm.DB) {
	result := m.conn.Create(&n)
	return result
}

func (m *defaultNeDeviceDataAccessItemModel) Find(filter interface{}) (result []NeDeviceDataAccessItem, err error) {
	tx := m.conn.Where(filter).Find(&result)
	if tx.Error != nil {
		return result, tx.Error
	}
	return result, nil
}

func (m *defaultNeDeviceDataAccessItemModel) Updates(filter interface{}, data NeDeviceDataAccessItem) (err error) {
	tx := m.conn.Model(&filter).Updates(data)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (m *defaultNeDeviceDataAccessItemModel) Delete(data NeDeviceDataAccessItem) (err error) {
	tx := m.conn.Delete(data)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (m *defaultNeDeviceDataAccessItemModel) Transaction(fc func(tx *gorm.DB) error,
	opts ...*sql.TxOptions) (err error) {
	return m.conn.Transaction(fc)
}
