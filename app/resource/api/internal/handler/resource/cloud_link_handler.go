package resource

import (
	"net/http"

	"T-driver/app/resource/api/internal/logic/resource"
	"T-driver/app/resource/api/internal/svc"
	"T-driver/app/resource/api/internal/types"
	"T-driver/common/errors"
	"T-driver/common/response"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func CloudLinkHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CloudLinkReq
		if err := httpx.Parse(r, &req); err != nil {
			logx.Error(err)
			response.Response(w, nil, errors.ParamsError(r.Header.Get("Language")))
			return
		}

		l := resource.NewCloudLinkLogic(r.Context(), svcCtx)
		resp, err := l.CloudLink(&req)
		response.Response(w, resp, err)
	}
}
