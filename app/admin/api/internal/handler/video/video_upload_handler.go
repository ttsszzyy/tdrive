package video

import (
	"T-driver/app/admin/api/internal/logic/video"
	"T-driver/app/admin/api/internal/svc"
	"T-driver/app/admin/api/internal/types"
	"T-driver/common/errors"
	"T-driver/common/response"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"io"
	"net/http"
)

func VideoUploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.VideoUploadReq
		if err := httpx.Parse(r, &req); err != nil {
			logx.Error(err)
			response.Response(w, nil, errors.ParamsError(r.Header.Get("Language")))
			return
		}
		// 从表单中提取文件
		file, handler, err := r.FormFile("file")
		if err != nil {
			logx.Error(err)
			response.Response(w, nil, errors.ParamsError(r.Header.Get("Language")))
			return
		}
		defer file.Close()

		req.File, err = io.ReadAll(file)
		if err != nil {
			logx.Error(err)
			response.Response(w, nil, errors.CustomError("保存文件错误"))
			return
		}
		req.Filename = handler.Filename
		req.Desc = r.FormValue("desc")
		req.Title = r.FormValue("title")

		l := video.NewVideoUploadLogic(r.Context(), svcCtx)
		resp, err := l.VideoUpload(&req)
		response.Response(w, resp, err)
	}
}
