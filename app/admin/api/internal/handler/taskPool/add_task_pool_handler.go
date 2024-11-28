package taskPool

import (
	"net/http"

	"T-driver/app/admin/api/internal/logic/taskPool"
	"T-driver/app/admin/api/internal/svc"
	"T-driver/app/admin/api/internal/types"
	"T-driver/common/errors"
	"T-driver/common/response"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func AddTaskPoolHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddTaskPoolReq
		if err := httpx.Parse(r, &req); err != nil {
			logx.Error(err)
			response.Response(w, nil, errors.ParamsError(r.Header.Get("Language")))
			return
		}

		l := taskPool.NewAddTaskPoolLogic(r.Context(), svcCtx)
		resp, err := l.AddTaskPool(&req)
		response.Response(w, resp, err)
	}
}
