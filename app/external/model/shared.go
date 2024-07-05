package model

import (
	sqlnull "database/sql"
	"fmt"
	"reflect"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"xunjikeji.com.cn/gateway/common/util"
)

const (
	DESC         = 1
	ASC          = 0
	Nid          = "nid"
	EdgeDeviceId = "edge_device_id"
	ConfigId     = "config_id"
)

type PageInfo struct {
	Page   int64     `json:"page" form:"page"`         // 页码
	Size   int64     `json:"pageSize" form:"pageSize"` // 每页大小
	Orders []OrderBy `json:"orderBy" form:"orderBy"`   // 排序信息
}

type OrderBy struct {
	Filed string `json:"filed" form:"filed"` //要排序的字段名
	Sort  int    `json:"sort" form:"sort"`   //排序的方式：0 降序、1 升序
}

func NewPageInfo(page, pageSize int64, orders ...OrderBy) (pageInfo PageInfo) {
	pageInfo = PageInfo{
		Page:   page,
		Size:   pageSize,
		Orders: orders,
	}
	return pageInfo
}

func NewOrder(filed string, sort int) OrderBy {
	return OrderBy{
		Filed: filed,
		Sort:  sort,
	}
}

func (p *PageInfo) PostOrder(orders ...OrderBy) {
	p.Orders = append(p.Orders, orders...)
}

func (p *PageInfo) GetLimit() int64 {
	if p.Page <= 0 {
		return 0
	}
	return p.Size
}

func (p *PageInfo) GetOffset() int64 {
	if p.Page == 0 {
		return 0
	}
	return p.Size * (p.Page - 1)
}

func (p *PageInfo) GetOrder() (orderSql string) {
	var fieldArray []string

	if len(p.Orders) == 0 {
		return
	}

	for _, v := range p.Orders {
		if v.Sort == 0 {
			fieldArray = append(fieldArray, fmt.Sprintf("`%s` asc", v.Filed))
		} else if v.Sort == 1 {
			fieldArray = append(fieldArray, fmt.Sprintf("`%s` desc", v.Filed))
		}
	}

	orderSql = fmt.Sprintf("%s", strings.Join(fieldArray, ","))

	return
}

// 将struct 属性的的非零值转换为where 条件

func CommonFmtSql(sql sq.SelectBuilder, data any) sq.SelectBuilder {
	//// 获取对象的反射值
	value := reflect.ValueOf(data)
	//
	// 判断类型
	switch value.Kind() {
	case reflect.Map:
		sql = sql.Where(data)
	case reflect.Struct:
		sql = StructFmtSql(sql, value)
	case reflect.String:
		sql = sql.Where(data)
	default:
		sql = sql.Where(data)
	}
	return sql
}

func MapFmtSql(sql sq.SelectBuilder, value reflect.Value) sq.SelectBuilder {
	for _, key := range value.MapKeys() {
		sql = sql.Where(fmt.Sprintf("%s = ?", key.Interface().(string)), value.MapIndex(key).Interface())
	}
	return sql
}

// StructFmtSql 将对象属性的非零值组装成sql 使用时注意是否支持该属性的类型 主要支持string int 和 float 类型
func StructFmtSql(sql sq.SelectBuilder, value reflect.Value) sq.SelectBuilder {
	typ := value.Type()
	for i := 0; i < value.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := value.Field(i)
		tag := field.Tag.Get("db")
		// 属性位该类型的零值
		if util.IsZeroValue(fieldValue.Interface()) {
			continue
		}

		var condition string
		var args []interface{}
		switch fieldValue.Interface().(type) {
		case sqlnull.NullString:
			s := fieldValue.Interface().(sqlnull.NullString)
			if s.Valid {
				condition = fmt.Sprintf("%s = ?", tag)
				args = []interface{}{s.String}
			}

		case sqlnull.NullInt64:
			n := fieldValue.Interface().(sqlnull.NullInt64)
			if n.Valid {
				condition = fmt.Sprintf("%s = ?", tag)
				args = []interface{}{n.Int64}
			}

		case sqlnull.NullInt32:
			n := fieldValue.Interface().(sqlnull.NullInt32)
			if n.Valid {
				condition = fmt.Sprintf("%s = ?", tag)
				args = []interface{}{n.Int32}
			}

		case sqlnull.NullInt16:
			n := fieldValue.Interface().(sqlnull.NullInt16)
			if n.Valid {
				condition = fmt.Sprintf("%s = ?", tag)
				args = []interface{}{n.Int16}
			}

		case sqlnull.NullFloat64:
			n := fieldValue.Interface().(sqlnull.NullFloat64)
			if n.Valid {
				condition = fmt.Sprintf("%s = ?", tag)
				args = []interface{}{n.Float64}
			}

		default:
			condition = fmt.Sprintf("%s = ?", tag)
			args = []interface{}{fieldValue.Interface()}
		}

		if condition != "" {
			sql = sql.Where(condition, args...)
		}
	}

	return sql
}
