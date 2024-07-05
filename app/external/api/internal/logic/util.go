package logic

import (
	"sort"

	"xunjikeji.com.cn/gateway/app/external/api/internal/types"
)

func sortDropdown(data []types.Dropdown) {
	sort.Slice(data, func(i, j int) bool {
		switch data[i].Label.(type) {
		case int:
			return data[i].Label.(int) < data[j].Label.(int)
		case string:
			return data[i].Label.(string) < data[j].Label.(string)
		default:
			return false
		}
	})
}
