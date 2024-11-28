package share

import (
	"net/http"

	"T-driver/app/user/api/internal/logic/share"
	"T-driver/app/user/api/internal/svc"
	"T-driver/app/user/api/internal/types"
	"T-driver/common/errors"
	"T-driver/common/response"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ShareFileUrlHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetShareURLReq
		if err := httpx.Parse(r, &req); err != nil {
			logx.Error(err)
			response.Response(w, nil, errors.ParamsError(r.Header.Get("Language")))
			return
		}

		l := share.NewShareFileUrlLogic(r.Context(), svcCtx)
		resp, err := l.ShareFileUrl(&req)
		if err != nil {
			response.Response(w, resp, err)
			return
		}
		http.Redirect(w, r, resp, http.StatusMovedPermanently)
	}
}
