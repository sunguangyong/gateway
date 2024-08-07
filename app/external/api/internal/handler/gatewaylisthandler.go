package handler

import (
	"net/http"

	"xunjikeji.com.cn/gateway/common/result"

	"github.com/zeromicro/go-zero/rest/httpx"
	"xunjikeji.com.cn/gateway/app/external/api/internal/logic"
	"xunjikeji.com.cn/gateway/app/external/api/internal/svc"
	"xunjikeji.com.cn/gateway/app/external/api/internal/types"
)

func gatewayListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.EdgeDeviceListRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewGatewayListLogic(r.Context(), svcCtx)
		resp, err := l.GatewayList(&req)
		result.HttpResult(r, w, req, resp, err)
	}
}
