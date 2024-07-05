package model

import (
	"database/sql"

	"gorm.io/gorm"
)

type (
	SqliteNeDeviceDataAccessConfigModel interface {
		Insert(n *NeDeviceDataAccessConfig) (tx *gorm.DB)
		Find(filter interface{}) (result []NeDeviceDataAccessConfig, err error)
		Updates(filter interface{}, data NeDeviceDataAccessConfig) (err error)
		Delete(data NeDeviceDataAccessConfig) (err error)
		Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) (err error)
	}

	NeDeviceDataAccessConfig struct {
		Nid               int64  `gorm:"column:nid;type:bigint(20);primary_key;AUTO_INCREMENT;comment:主键ID" json:"nid"`
		DeviceNid         int64  `gorm:"column:device_nid;type:bigint(20);comment:设备ID;NOT NULL" json:"device_nid"`
		ConfigType        int64  `gorm:"column:config_type;type:int(11);default:0;comment:0-读/写,1-只读,2-只写" json:"config_type"`
		ConfigId          string `gorm:"column:config_id;type:varchar(50);comment:配置ID" json:"config_id"`
		ConfigName        string `gorm:"column:config_name;type:varchar(100);comment:配置名称;NOT NULL" json:"config_name"`
		Endpoint          string `gorm:"column:endpoint;type:varchar(3000);comment:连接URL;NOT NULL" json:"endpoint"`
		Protocol          string `gorm:"column:protocol;type:varchar(50);comment:访问协议：ModbusTcp, ModbusRtu, OpcUa, HTTP, MTCONNECT, MITSUBISH_MC, SIEMENS_S7, OMRON_SINS;NOT NULL" json:"protocol"`
		JsonAccessOptions string `gorm:"column:json_access_options;type:text;comment:访问连接基本配置(JSON格式);NOT NULL" json:"json_access_options"`
		Timeout           int64  `gorm:"column:timeout;type:int(11);default:1000;comment:连接超时时间" json:"timeout"`
		AgwId             int64  `gorm:"column:agw_id;type:bigint(20);comment:所属网关ID" json:"agw_id"`
		TenantId          int64  `gorm:"column:tenant_id;type:bigint(20);comment:租户ID" json:"tenant_id"`
		ProfileNid        int64  `gorm:"column:profile_nid;type:bigint(20);comment:设备配置文件ID" json:"profile_nid"`
		CreateTime        string `gorm:"column:create_time;type:datetime;comment:创建时间;NOT NULL" json:"create_time"`
		CreateBy          int64  `gorm:"column:create_by;type:bigint(20);comment:创建者ID;NOT NULL" json:"create_by"`
		Issued            int64  `gorm:"column:issued;type:tinyint(1);default:0;comment:是否已下发：1-是,0-否;NOT NULL" json:"issued"`
		IssueTime         string `gorm:"column:issue_time;type:datetime;comment:下发时间" json:"issue_time"`
	}

	defaultNeDeviceDataAccessConfigModel struct {
		conn  *gorm.DB
		table string
	}
)

func (m *NeDeviceDataAccessConfig) TableName() string {
	return "`ne_device_data_access_config`"
}

func newNeDeviceDataAccessConfigModel(conn *gorm.DB) *defaultNeDeviceDataAccessConfigModel {
	return &defaultNeDeviceDataAccessConfigModel{
		conn:  conn,
		table: "`ne_device_data_access_config`",
	}
}

func NewNeDeviceDataAccessConfigModel(conn *gorm.DB) SqliteNeDeviceDataAccessConfigModel {
	return newNeDeviceDataAccessConfigModel(conn)
}

func (m *defaultNeDeviceDataAccessConfigModel) Insert(n *NeDeviceDataAccessConfig) (tx *gorm.DB) {
	result := m.conn.Create(&n)
	return result
}

func (m *defaultNeDeviceDataAccessConfigModel) Find(filter interface{}) (result []NeDeviceDataAccessConfig, err error) {
	tx := m.conn.Where(filter).Find(&result)
	if tx.Error != nil {
		return result, tx.Error
	}
	return result, nil
}

func (m *defaultNeDeviceDataAccessConfigModel) Updates(filter interface{}, data NeDeviceDataAccessConfig) (err error) {
	tx := m.conn.Model(&filter).Updates(data)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (m *defaultNeDeviceDataAccessConfigModel) Delete(data NeDeviceDataAccessConfig) (err error) {
	tx := m.conn.Delete(data)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (m *defaultNeDeviceDataAccessConfigModel) Transaction(fc func(tx *gorm.DB) error,
	opts ...*sql.TxOptions) (err error) {

	return m.conn.Transaction(fc)
}

/*
db.Transaction(func(tx *gorm.DB) error {
		tx.Create(&user1)
		tx.Transaction(func(tx2 *gorm.DB) error {
			tx2.Create(&user2)
			return errors.New("rollback user2") // Rollback user2
		})
		tx.Transaction(func(tx2 *gorm.DB) error {
			tx2.Create(&user3)
			return nil
		})
		return nil
	})
*/
