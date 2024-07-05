package handler

import (
	"encoding/json"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"xunjikeji.com.cn/gateway/common/result"

	"xunjikeji.com.cn/gateway/app/interior/api/internal/logic"
	"xunjikeji.com.cn/gateway/app/interior/api/internal/svc"
	"xunjikeji.com.cn/gateway/app/interior/api/internal/types"
)

func protocolWriteTestHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ProtocolWriteTestRequest
		//if err := httpx.Parse(r, &req); err != nil {
		//	httpx.Error(w, err)
		//	return
		//}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewProtocolWriteTestLogic(r.Context(), svcCtx)
		resp, err := l.ProtocolWriteTest(&req)
		result.HttpResult(r, w, req, resp, err)
	}
}
