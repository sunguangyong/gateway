package handler

import (
	"fmt"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"time"
	"xunjikeji.com.cn/gateway/app/external/api/internal/logic"
	"xunjikeji.com.cn/gateway/app/external/api/internal/svc"
	"xunjikeji.com.cn/gateway/app/external/api/internal/types"
	"xunjikeji.com.cn/gateway/common/excel"
)

func GatewayListExportHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.EdgeDeviceListRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewGatewayListLogic(r.Context(), svcCtx)
		resp, err := l.GatewayList(&req)

		if err != nil {
			return
		}

		var interfaceSlice []interface{}
		for _, item := range resp.Data {
			interfaceSlice = append(interfaceSlice, item)
		}

		disposition := fmt.Sprintf("attachment; filename=%s.xlsx", time.Now().Format("2006-01-02-15-04-05"))
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", disposition)
		w.Header().Set("Content-Transfer-Encoding", "binary")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")
		excel.ExportByStruct(w, []string{"网关id", "网关名称", "网关位置", "ip地址", "mac地址", "设备描述", "sim卡", "网关型号", "在线"}, interfaceSlice,
			"b")
	}
}
