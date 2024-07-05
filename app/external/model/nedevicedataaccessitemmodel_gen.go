// Code generated by goctl. DO NOT EDIT!

package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	neDeviceDataAccessItemFieldNames          = builder.RawFieldNames(&NeDeviceDataAccessItem{})
	neDeviceDataAccessItemRows                = strings.Join(neDeviceDataAccessItemFieldNames, ",")
	neDeviceDataAccessItemRowsExpectAutoSet   = strings.Join(stringx.Remove(neDeviceDataAccessItemFieldNames, "`nid`", "`create_time`", "`update_time`"), ",")
	neDeviceDataAccessItemRowsWithPlaceHolder = strings.Join(stringx.Remove(neDeviceDataAccessItemFieldNames, "`nid`", "`create_time`", "`update_time`"), "=?,") + "=?"
)

type (
	neDeviceDataAccessItemModel interface {
		Insert(ctx context.Context, data *NeDeviceDataAccessItem) (sql.Result, error)
		FindOne(ctx context.Context, nid int64) (*NeDeviceDataAccessItem, error)
		Update(ctx context.Context, data *NeDeviceDataAccessItem) error
		Delete(ctx context.Context, nid int64) error
	}

	defaultNeDeviceDataAccessItemModel struct {
		conn  sqlx.SqlConn
		table string
	}

	NeDeviceDataAccessItem struct {
		Nid        int64     `db:"nid"`         // 主键
		ConfigNid  int64     `db:"config_nid"`  // ne_device_data_access_config.id
		DeviceNid  int64     `db:"device_nid"`  // 设备ID
		AgwId      int64     `db:"agw_id"`      // 所属网关ID
		TenantId   int64     `db:"tenant_id"`   // 租户ID
		ConfigType int64     `db:"config_type"` // 0-读/写,1-只读,2-只写
		AccessData string    `db:"access_data"` // 读/写配置项(JSON对象)
		CreateTime time.Time `db:"create_time"` // 创建时间
	}
)

func newNeDeviceDataAccessItemModel(conn sqlx.SqlConn) *defaultNeDeviceDataAccessItemModel {
	return &defaultNeDeviceDataAccessItemModel{
		conn:  conn,
		table: "`ne_device_data_access_item`",
	}
}

func (m *defaultNeDeviceDataAccessItemModel) Insert(ctx context.Context, data *NeDeviceDataAccessItem) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?)", m.table, neDeviceDataAccessItemRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.ConfigNid, data.DeviceNid, data.AgwId, data.TenantId, data.ConfigType, data.AccessData)
	return ret, err
}

func (m *defaultNeDeviceDataAccessItemModel) FindOne(ctx context.Context, nid int64) (*NeDeviceDataAccessItem, error) {
	query := fmt.Sprintf("select %s from %s where `nid` = ? limit 1", neDeviceDataAccessItemRows, m.table)
	var resp NeDeviceDataAccessItem
	err := m.conn.QueryRowCtx(ctx, &resp, query, nid)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultNeDeviceDataAccessItemModel) Update(ctx context.Context, data *NeDeviceDataAccessItem) error {
	query := fmt.Sprintf("update %s set %s where `nid` = ?", m.table, neDeviceDataAccessItemRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.ConfigNid, data.DeviceNid, data.AgwId, data.TenantId, data.ConfigType, data.AccessData, data.Nid)
	return err
}

func (m *defaultNeDeviceDataAccessItemModel) Delete(ctx context.Context, nid int64) error {
	query := fmt.Sprintf("delete from %s where `nid` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, nid)
	return err
}

func (m *defaultNeDeviceDataAccessItemModel) tableName() string {
	return m.table
}
