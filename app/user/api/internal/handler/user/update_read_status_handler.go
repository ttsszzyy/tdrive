package user

import (
	"net/http"

	"T-driver/app/user/api/internal/logic/user"
	"T-driver/app/user/api/internal/svc"
	"T-driver/common/response"
)

func UpdateReadStatusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewUpdateReadStatusLogic(r.Context(), svcCtx)
		err := l.UpdateReadStatus()
		response.Response(w, nil, err)
	}
}
