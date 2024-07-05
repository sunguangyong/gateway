// Code generated by goctl. DO NOT EDIT!

package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	feEdgeDeviceFieldNames          = builder.RawFieldNames(&FeEdgeDevice{})
	feEdgeDeviceRows                = strings.Join(feEdgeDeviceFieldNames, ",")
	feEdgeDeviceRowsExpectAutoSet   = strings.Join(stringx.Remove(feEdgeDeviceFieldNames, "`edge_device_id`", "`create_time`", "`update_time`"), ",")
	feEdgeDeviceRowsWithPlaceHolder = strings.Join(stringx.Remove(feEdgeDeviceFieldNames, "`edge_device_id`", "`create_time`", "`update_time`"), "=?,") + "=?"
)

type (
	feEdgeDeviceModel interface {
		Insert(ctx context.Context, data *FeEdgeDevice) (sql.Result, error)
		FindOne(ctx context.Context, edgeDeviceId int64) (*FeEdgeDevice, error)
		FindOneByMac(ctx context.Context, mac string) (*FeEdgeDevice, error)
		Update(ctx context.Context, data *FeEdgeDevice) error
		Delete(ctx context.Context, edgeDeviceId int64) error
	}

	defaultFeEdgeDeviceModel struct {
		conn  sqlx.SqlConn
		table string
	}

	FeEdgeDevice struct {
		EdgeDeviceId      int64     `db:"edge_device_id"`   // 网关设备ID
		EdgeDeviceName    string    `db:"edge_device_name"` // 网关设备名称
		DeviceDesc        string    `db:"device_desc"`      // 设备描述
		DevicePosition    string    `db:"device_position"`  // 设备位置
		IpAddress         string    `db:"ip_address"`       // 设备IP地址
		EdgeGroupId       int64     `db:"edge_group_id"`    // 设备组ID
		CreateTime        time.Time `db:"create_time"`      // 创建时间
		CreateBy          int64     `db:"create_by"`        // 创建者ID
		UpdateTime        time.Time `db:"update_time"`      // 更新时间
		UpdateBy          int64     `db:"update_by"`        // 更新者ID
		EdgeCloudId       int64     `db:"edge_cloud_id"`    // 所属边缘云id
		SerialNum         string    `db:"serial_num"`       // 网关设备序列号
		TenantId          int64     `db:"tenant_id"`        // 设备所属租户ID
		Area              string    `db:"area"`             // 所属区域
		Mac               string    `db:"mac"`              // mac地址
		Status            int64     `db:"status"`           // 在线/离线
		ReportTime        time.Time `db:"report_time"`      // 上报时间
		IsDeleted         int64     `db:"is_deleted"`       // 是否待删除 0----否 ，1-----是
		SimCard           string    `db:"sim_card"`
		NodePort          string    `db:"node_port"`          // 节点端口
		Type              int64     `db:"type"`               // 1----网关  2-----节点
		AuthorizationCode string    `db:"authorization_code"` // 授权码
		IsEffective       int64     `db:"is_effective"`       // 授权码是否生效 0----否 1----是
		IsKong            int64     `db:"is_kong"`            // 是否调用Kong  0----否 1----是
		GatewayId         int64     `db:"gateway_id"`         // 网关型号id
		ProjectId         int64     `db:"project_id"`         // 所属项目id
		DeviceSn          string    `db:"device_sn"`          // 设备SN
		DeviceLocation    string    `db:"device_location"`    // 空间标识
		GatewayType       string    `db:"gateway_type"`       // 网关型号
		GatewayCategory   string    `db:"gateway_category"`   // 网关类别
		SimIpAddress      string    `db:"sim_ip_address"`     // sim卡ip地址
	}
)

func newFeEdgeDeviceModel(conn sqlx.SqlConn) *defaultFeEdgeDeviceModel {
	return &defaultFeEdgeDeviceModel{
		conn:  conn,
		table: "`fe_edge_device`",
	}
}

func (m *defaultFeEdgeDeviceModel) Insert(ctx context.Context, data *FeEdgeDevice) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, feEdgeDeviceRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.EdgeDeviceName, data.DeviceDesc, data.DevicePosition, data.IpAddress, data.EdgeGroupId, data.CreateBy, data.UpdateBy, data.EdgeCloudId, data.SerialNum, data.TenantId, data.Area, data.Mac, data.Status, data.ReportTime, data.IsDeleted, data.SimCard, data.NodePort, data.Type, data.AuthorizationCode, data.IsEffective, data.IsKong, data.GatewayId, data.ProjectId, data.DeviceSn, data.DeviceLocation, data.GatewayType, data.GatewayCategory, data.SimIpAddress)
	return ret, err
}

func (m *defaultFeEdgeDeviceModel) FindOne(ctx context.Context, edgeDeviceId int64) (*FeEdgeDevice, error) {
	query := fmt.Sprintf("select %s from %s where `edge_device_id` = ? limit 1", feEdgeDeviceRows, m.table)
	var resp FeEdgeDevice
	err := m.conn.QueryRowCtx(ctx, &resp, query, edgeDeviceId)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultFeEdgeDeviceModel) FindOneByMac(ctx context.Context, mac string) (*FeEdgeDevice, error) {
	var resp FeEdgeDevice
	query := fmt.Sprintf("select %s from %s where `mac` = ? limit 1", feEdgeDeviceRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, mac)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultFeEdgeDeviceModel) Update(ctx context.Context, data *FeEdgeDevice) error {
	query := fmt.Sprintf("update %s set %s where `edge_device_id` = ?", m.table, feEdgeDeviceRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.EdgeDeviceName, data.DeviceDesc, data.DevicePosition, data.IpAddress, data.EdgeGroupId, data.CreateBy, data.UpdateBy, data.EdgeCloudId, data.SerialNum, data.TenantId, data.Area, data.Mac, data.Status, data.ReportTime, data.IsDeleted, data.SimCard, data.NodePort, data.Type, data.AuthorizationCode, data.IsEffective, data.IsKong, data.GatewayId, data.ProjectId, data.DeviceSn, data.DeviceLocation, data.GatewayType, data.GatewayCategory, data.SimIpAddress, data.EdgeDeviceId)
	return err
}

func (m *defaultFeEdgeDeviceModel) Delete(ctx context.Context, edgeDeviceId int64) error {
	query := fmt.Sprintf("delete from %s where `edge_device_id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, edgeDeviceId)
	return err
}

func (m *defaultFeEdgeDeviceModel) tableName() string {
	return m.table
}