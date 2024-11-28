package trade

import (
	"net/http"

	"T-driver/app/user/api/internal/logic/trade"
	"T-driver/app/user/api/internal/svc"
	"T-driver/common/response"
)

func GetAirdropTokenHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := trade.NewGetAirdropTokenLogic(r.Context(), svcCtx)
		resp, err := l.GetAirdropToken()
		response.Response(w, resp, err)
	}
}
