package user

import (
	"net/http"

	"T-driver/app/user/api/internal/logic/user"
	"T-driver/app/user/api/internal/svc"
	"T-driver/common/response"
)

func GetRsaPublicKeyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewGetRsaPublicKeyLogic(r.Context(), svcCtx)
		resp, err := l.GetRsaPublicKey()
		response.Response(w, resp, err)
	}
}
