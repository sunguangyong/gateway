package model

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ NeDeviceModel = (*customNeDeviceModel)(nil)

type (
	// NeDeviceModel is an interface to be customized, add more methods here,
	// and implement the added methods in customNeDeviceModel.
	NeDeviceModel interface {
		neDeviceModel
		CommonFilterCount(ctx context.Context, filter any) (size int64, err error)
		CommonFilterFind(ctx context.Context, filter any, page PageInfo) ([]*NeDevice, error)
	}

	customNeDeviceModel struct {
		*defaultNeDeviceModel
	}
)

// NewNeDeviceModel returns a model for the database table.
func NewNeDeviceModel(conn sqlx.SqlConn) NeDeviceModel {
	return &customNeDeviceModel{
		defaultNeDeviceModel: newNeDeviceModel(conn),
	}
}

func (m *defaultNeDeviceModel) CommonFilterFind(ctx context.Context, filter any,
	page PageInfo) ([]*NeDevice, error) {

	var resp []*NeDevice
	var sql sq.SelectBuilder
	sql = sq.Select(neDeviceRows).From(m.table)

	if page.GetLimit() > 0 {
		sql = sql.Limit(uint64(page.GetLimit())).Offset(uint64(page.GetOffset()))
	}

	if page.GetOrder() != "" {
		sql = sql.OrderBy(page.GetOrder())
	}

	sql = CommonFmtSql(sql, filter)

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

func (m *defaultNeDeviceModel) CommonFilterCount(ctx context.Context, filter any) (size int64,
	err error) {

	sql := sq.Select("count(1)").From(m.table)

	sql = CommonFmtSql(sql, filter)

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
