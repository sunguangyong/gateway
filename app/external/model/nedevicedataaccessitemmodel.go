package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"

	"github.com/zeromicro/go-zero/core/stringx"

	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ NeDeviceDataAccessItemModel = (*customNeDeviceDataAccessItemModel)(nil)

type (
	// NeDeviceDataAccessItemModel is an interface to be customized, add more methods here,
	// and implement the added methods in customNeDeviceDataAccessItemModel.
	NeDeviceDataAccessItemModel interface {
		neDeviceDataAccessItemModel
		//InsertTime(ctx context.Context, data *NeDeviceDataAccessItem) (sql.Result, error)
		//CommonFind(ctx context.Context, querySql, orderSql, limitSql string) ([]*NeDeviceDataAccessItem, error)
		FindNewsCount(ctx context.Context, paramSql string) (count int64, err error)
		CommonFilterFind(ctx context.Context, f any,
			page PageInfo) ([]*NeDeviceDataAccessItem, error)
		CommonFilterCount(ctx context.Context, f any) (size int64, err error)
	}

	customNeDeviceDataAccessItemModel struct {
		*defaultNeDeviceDataAccessItemModel
	}
)

// NewNeDeviceDataAccessItemModel returns a model for the database table.
func NewNeDeviceDataAccessItemModel(conn sqlx.SqlConn) NeDeviceDataAccessItemModel {
	return &customNeDeviceDataAccessItemModel{
		defaultNeDeviceDataAccessItemModel: newNeDeviceDataAccessItemModel(conn),
	}
}

func (m *defaultNeDeviceDataAccessItemModel) InsertTime(ctx context.Context, data *NeDeviceDataAccessItem) (sql.Result,
	error) {
	neDeviceDataAccessItemInsertRows := strings.Join(stringx.Remove(neDeviceDataAccessItemFieldNames, "`nid`"), ",")
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?)", m.table, neDeviceDataAccessItemInsertRows)
	ret, err := m.conn.ExecCtx(ctx, query, data.ConfigNid, data.ConfigType, data.AccessData, data.CreateTime)
	return ret, err
}

func (m *defaultNeDeviceDataAccessItemModel) CommonFind(ctx context.Context, querySql, orderSql, limitSql string) ([]*NeDeviceDataAccessItem, error) {
	query := fmt.Sprintf("select %s from %s %s %s %s", neDeviceDataAccessItemRows, m.table, querySql, orderSql, limitSql)
	var resp []*NeDeviceDataAccessItem
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

func (m *defaultNeDeviceDataAccessItemModel) FindNewsCount(ctx context.Context, paramSql string) (count int64, err error) {
	query := fmt.Sprintf(`select count(1) from %s %s`, m.table, paramSql)
	err = m.conn.QueryRowCtx(ctx, &count, query)

	switch err {
	case nil:
		return count, nil
	case sqlc.ErrNotFound:
		return 0, ErrNotFound
	default:
		return 0, err
	}
}

func (m *defaultNeDeviceDataAccessItemModel) CommonFilterFind(ctx context.Context, f any,
	page PageInfo) ([]*NeDeviceDataAccessItem, error) {

	var resp []*NeDeviceDataAccessItem
	var sql sq.SelectBuilder
	sql = sq.Select(neDeviceDataAccessItemRows).From(m.table)

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

func (m *defaultNeDeviceDataAccessItemModel) CommonFilterCount(ctx context.Context,
	f any) (size int64, err error) {

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
