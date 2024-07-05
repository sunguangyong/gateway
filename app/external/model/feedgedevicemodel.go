package model

import (
	"context"
	"fmt"
	"strings"

	"xunjikeji.com.cn/gateway/common/xerr"

	sq "github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlc"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ FeEdgeDeviceModel = (*customFeEdgeDeviceModel)(nil)

type (
	// FeEdgeDeviceModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFeEdgeDeviceModel.
	FeEdgeDeviceModel interface {
		feEdgeDeviceModel
		//CommonFind(ctx context.Context, querySql, orderSql, limitSql string) ([]*FeEdgeDevice, error)
		CommonFilterFind(ctx context.Context, filter any, page PageInfo) ([]*FeEdgeDevice, error)
		CommonFilterCount(ctx context.Context, filter any) (size int64, err error)
		ConvertUniqIdxMacErr(err error) error
		GetHost(ctx context.Context, agwId int64) (host string, err error)
	}

	customFeEdgeDeviceModel struct {
		*defaultFeEdgeDeviceModel
	}
)

// NewFeEdgeDeviceModel returns a model for the database table.
func NewFeEdgeDeviceModel(conn sqlx.SqlConn) FeEdgeDeviceModel {
	return &customFeEdgeDeviceModel{
		defaultFeEdgeDeviceModel: newFeEdgeDeviceModel(conn),
	}
}

func (m *defaultFeEdgeDeviceModel) CommonFind(ctx context.Context, querySql, orderSql, limitSql string) ([]*FeEdgeDevice, error) {
	query := fmt.Sprintf("select %s from %s %s %s %s", feEdgeDeviceRows, m.table, querySql, orderSql, limitSql)
	var resp []*FeEdgeDevice
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

func (m *defaultFeEdgeDeviceModel) CommonFilterFind(ctx context.Context, filter any, page PageInfo) ([]*FeEdgeDevice, error) {

	var resp []*FeEdgeDevice
	var sql sq.SelectBuilder
	sql = sq.Select(feEdgeDeviceRows).From(m.table)

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

func (m *defaultFeEdgeDeviceModel) CommonFilterCount(ctx context.Context, filter any) (size int64,
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

const (
	UniqIdxFeEdgeDeviceModelMac = "uniq_idx_fe_edge_device_mac"
)

func (m *defaultFeEdgeDeviceModel) ConvertUniqIdxMacErr(err error) error {
	if err == nil {
		return nil
	}
	ok := strings.Contains(err.Error(), UniqIdxFeEdgeDeviceModelMac)
	if ok {
		return xerr.MacExistsError
	}
	return err
}

func (m *defaultFeEdgeDeviceModel) GetHost(ctx context.Context, agwId int64) (host string, err error) {
	querySql := fmt.Sprintf("where edge_device_id = %d", agwId)
	limitSql := "limit 1"
	query := fmt.Sprintf("select %s from %s %s %s %s", feEdgeDeviceRows, m.table, querySql, "", limitSql)
	var resp []*FeEdgeDevice
	err = m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		if len(resp) > 0 {
			return resp[0].IpAddress, nil
		}
		return "", nil
	case sqlx.ErrNotFound:
		return "", nil
	default:
		return "", err
	}
}
