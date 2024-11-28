package video

import (
	"net/http"

	"T-driver/app/video/api/internal/logic/video"
	"T-driver/app/video/api/internal/svc"
	"T-driver/app/video/api/internal/types"
	"T-driver/common/errors"
	"T-driver/common/response"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func VideoListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.VideoListReq
		if err := httpx.Parse(r, &req); err != nil {
			logx.Error(err)
			response.Response(w, nil, errors.ParamsError(r.Header.Get("Language")))
			return
		}

		l := video.NewVideoListLogic(r.Context(), svcCtx)
		resp, err := l.VideoList(&req)
		response.Response(w, resp, err)
	}
}