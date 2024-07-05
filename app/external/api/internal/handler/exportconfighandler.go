package handler

import (
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"xunjikeji.com.cn/gateway/app/external/api/internal/logic"
	"xunjikeji.com.cn/gateway/app/external/api/internal/svc"
	"xunjikeji.com.cn/gateway/app/external/api/internal/types"
)

func exportConfigHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ExportConfigRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewExportConfigLogic(r.Context(), svcCtx)
		resp, err := l.ExportConfig(&req)
		fmt.Println(resp)

		if err != nil {
			return
		}
		fileName := "config.json"
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
		w.Header().Set("Content-Transfer-Encoding", "binary")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")
		fmt.Println("llllllllllll", resp.Data)
		_, err = w.Write([]byte(resp.Data))
		if err != nil {
			httpx.Error(w, err)
			return
		}
		return
	}
}
