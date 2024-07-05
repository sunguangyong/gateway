package handler

import (
	"encoding/json"
	"net/http"

	"xunjikeji.com.cn/gateway/common/result"

	"github.com/zeromicro/go-zero/rest/httpx"
	"xunjikeji.com.cn/gateway/app/external/api/internal/logic"
	"xunjikeji.com.cn/gateway/app/external/api/internal/svc"
	"xunjikeji.com.cn/gateway/app/external/api/internal/types"
)

func protocolWriteTestHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ProtocolWriteTestRequest
		//if err := httpx.Parse(r, &req); err != nil {
		//	httpx.Error(w, err)
		//	return
		//}
		// 防止参数里的 interface err
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewProtocolWriteTestLogic(r.Context(), svcCtx)
		resp, err := l.ProtocolWriteTest(&req)
		result.HttpResult(r, w, req, resp, err)
	}
}
