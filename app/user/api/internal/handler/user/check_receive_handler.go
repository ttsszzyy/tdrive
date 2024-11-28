package user

import (
	"net/http"

	"T-driver/app/user/api/internal/logic/user"
	"T-driver/app/user/api/internal/svc"
	"T-driver/app/user/api/internal/types"
	"T-driver/common/errors"
	"T-driver/common/response"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func CheckReceiveHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			req types.ReceivePointsReq
			lan = r.Header.Get("Language")
		)
		if err := httpx.Parse(r, &req); err != nil {
			logx.Error(err)
			response.Response(w, nil, errors.ParamsError(lan))
			return
		}

		l := user.NewCheckReceiveLogic(r.Context(), svcCtx)
		resp, err := l.CheckReceive(&req)
		response.Response(w, resp, err)
	}
}