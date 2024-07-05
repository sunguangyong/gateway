package handler

import (
	"net/http"

	"xunjikeji.com.cn/gateway/common/result"

	"github.com/zeromicro/go-zero/rest/httpx"
	"xunjikeji.com.cn/gateway/app/interior/api/internal/logic"
	"xunjikeji.com.cn/gateway/app/interior/api/internal/svc"
	"xunjikeji.com.cn/gateway/app/interior/api/internal/types"
)

func protocolReadTestHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ProtocolTestRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewProtocolReadTestLogic(r.Context(), svcCtx)
		resp, err := l.ProtocolReadTest(&req)
		result.HttpResult(r, w, req, resp, err)
	}
}
