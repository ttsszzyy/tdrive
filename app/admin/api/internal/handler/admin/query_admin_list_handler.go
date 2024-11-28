package admin

import (
	"net/http"

	"T-driver/app/admin/api/internal/logic/admin"
	"T-driver/app/admin/api/internal/svc"
	"T-driver/app/admin/api/internal/types"
	"T-driver/common/errors"
	"T-driver/common/response"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func QueryAdminListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.QueryAdminListReq
		if err := httpx.Parse(r, &req); err != nil {
			logx.Error(err)
			response.Response(w, nil, errors.ParamsError(r.Header.Get("Language")))
			return
		}

		l := admin.NewQueryAdminListLogic(r.Context(), svcCtx)
		resp, err := l.QueryAdminList(&req)
		response.Response(w, resp, err)
	}
}
