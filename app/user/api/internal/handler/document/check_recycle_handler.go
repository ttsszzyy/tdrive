package document

import (
	"net/http"

	"T-driver/app/user/api/internal/logic/document"
	"T-driver/app/user/api/internal/svc"
	"T-driver/common/response"
)

func CheckRecycleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := document.NewCheckRecycleLogic(r.Context(), svcCtx)
		resp, err := l.CheckRecycle()
		response.Response(w, resp, err)
	}
}
