package parameter

import (
	"net/http"

	"T-driver/app/admin/api/internal/logic/parameter"
	"T-driver/app/admin/api/internal/svc"
	"T-driver/app/admin/api/internal/types"
	"T-driver/common/errors"
	"T-driver/common/response"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func DelBotCommandHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DelBotCommandReq
		if err := httpx.Parse(r, &req); err != nil {
			logx.Error(err)
			response.Response(w, nil, errors.ParamsError(r.Header.Get("Language")))
			return
		}

		l := parameter.NewDelBotCommandLogic(r.Context(), svcCtx)
		resp, err := l.DelBotCommand(&req)
		response.Response(w, resp, err)
	}
}