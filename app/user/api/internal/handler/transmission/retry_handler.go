package transmission

import (
	"net/http"

	"T-driver/app/user/api/internal/logic/transmission"
	"T-driver/app/user/api/internal/svc"
	"T-driver/app/user/api/internal/types"
	"T-driver/common/errors"
	"T-driver/common/response"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func RetryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RetryReq
		if err := httpx.Parse(r, &req); err != nil {
			logx.Error(err)
			response.Response(w, nil, errors.ParamsError(r.Header.Get("Language")))
			return
		}

		l := transmission.NewRetryLogic(r.Context(), svcCtx)
		resp, err := l.Retry(&req)
		response.Response(w, resp, err)
	}
}
