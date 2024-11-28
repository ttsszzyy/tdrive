package resource

import (
	"net/http"

	"T-driver/app/resource/api/internal/logic/resource"
	"T-driver/app/resource/api/internal/svc"
	"T-driver/app/resource/api/internal/types"
	"T-driver/common/errors"
	"T-driver/common/response"
	"T-driver/common/utils"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func RelayTitanHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RelayTitanReq
		if err := httpx.Parse(r, &req); err != nil {
			logx.Error(err)
			response.Response(w, nil, errors.ParamsError(r.Header.Get("Language")))
			return
		}

		// 获取请求ip
		header := utils.GetHeaderFromRequest(r)
		req.IP = header.IP

		l := resource.NewRelayTitanLogic(r.Context(), svcCtx)
		resp, err := l.RelayTitan(&req)
		if err != nil {
			response.Response(w, resp, err)
			return
		}

		http.Redirect(w, r, resp, http.StatusMovedPermanently)
	}
}
