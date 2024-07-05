package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stringx"

	sq "github.com/Masterminds/squirrel"

	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ NeDeviceDataAccessConfigModel = (*customNeDeviceDataAccessConfigModel)(nil)

type (
	// NeDeviceDataAccessConfigModel is an interface to be customized, add more methods here,
	// and implement the added methods in customNeDeviceDataAccessConfigModel.
	NeDeviceDataAccessConfigModel interface {
		neDeviceDataAccessConfigModel
		//CommonFind(ctx context.Context, querySql, orderSql, limitSql string) ([]*NeDeviceDataAccessConfig, error)
		FindNewsCount(ctx context.Context, paramSql string) (count int64, err error)
		//InsertTime(ctx context.Context, data *NeDeviceDataAccessConfig) (sql.Result, error)
		CommonFilterFind(ctx context.Context, f any,
			page PageInfo) ([]*NeDeviceDataAccessConfig, error)
		CommonFilterCount(ctx context.Context, f any) (count int64, err error)
	}

	customNeDeviceDataAccessConfigModel struct {
		*defaultNeDeviceDataAccessConfigModel
	}
)

// NewNeDeviceDataAccessConfigModel returns a model for the database table.
func NewNeDeviceDataAccessConfigModel(conn sqlx.SqlConn) NeDeviceDataAccessConfigModel {
	return &customNeDeviceDataAccessConfigModel{
		defaultNeDeviceDataAccessConfigModel: newNeDeviceDataAccessConfigModel(conn),
	}
}

func (m *defaultNeDeviceDataAccessConfigModel) CommonFind(ctx context.Context, querySql, orderSql, limitSql string) ([]*NeDeviceDataAccessConfig, error) {
	query := fmt.Sprintf("select %s from %s %s %s %s", neDeviceDataAccessConfigRows, m.table, querySql, orderSql, limitSql)
	var resp []*NeDeviceDataAccessConfig
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultNeDeviceDataAccessConfigModel) FindNewsCount(ctx context.Context, paramSql string) (size int64,
	err error) {
	query := fmt.Sprintf(`select count(1) from %s %s`, m.table, paramSql)
	err = m.conn.QueryRowCtx(ctx, &size, query)

	switch err {
	case nil:
		return size, nil
	case sqlc.ErrNotFound:
		return size, ErrNotFound
	default:
		return size, err
	}
}

func (m *defaultNeDeviceDataAccessConfigModel) InsertTime(ctx context.Context,
	data *NeDeviceDataAccessConfig) (sql.Result, error) {

	neDeviceDataAccessConfigTimeRows := strings.Join(stringx.Remove(neDeviceDataAccessConfigFieldNames, "`nid`"), ",")

	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table,
		neDeviceDataAccessConfigTimeRows)
	ret, err := m.conn.ExecCtx(ctx, query, data.DeviceNid, data.ConfigType, data.ConfigId, data.ConfigName,
		data.Endpoint, data.Protocol, data.JsonAccessOptions, data.Timeout, data.AgwId, data.TenantId,
		data.ProfileNid, data.CreateTime, data.CreateBy, data.Issued, data.IssueTime)
	return ret, err
}

func (m *defaultNeDeviceDataAccessConfigModel) CommonFilterFind(ctx context.Context, f any,
	page PageInfo) ([]*NeDeviceDataAccessConfig, error) {

	var resp []*NeDeviceDataAccessConfig
	var sql sq.SelectBuilder
	sql = sq.Select(neDeviceDataAccessConfigRows).From(m.table)

	if page.GetLimit() > 0 {
		sql = sql.Limit(uint64(page.GetLimit())).Offset(uint64(page.GetOffset()))
	}

	if page.GetOrder() != "" {
		sql = sql.OrderBy(page.GetOrder())
	}

	sql = CommonFmtSql(sql, f)

	query, arg, err := sql.ToSql()

	if err != nil {
		return nil, err
	}

	err = m.conn.QueryRowsCtx(ctx, &resp, query, arg...)

	switch err {
	case nil:
		return resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}

}

func (m *defaultNeDeviceDataAccessConfigModel) CommonFilterCount(ctx context.Context, f any) (size int64, err error) {

	sql := sq.Select("count(1)").From(m.table)

	sql = CommonFmtSql(sql, f)

	query, arg, err := sql.ToSql()

	if err != nil {
		return 0, err
	}

	err = m.conn.QueryRowCtx(ctx, &size, query, arg...)

	switch err {
	case nil:
		return size, nil
	case sqlc.ErrNotFound:
		return size, ErrNotFound
	default:
		return size, err
	}

}
