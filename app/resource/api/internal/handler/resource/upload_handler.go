package resource

import (
	"io"
	"net/http"
	"strconv"

	"T-driver/app/resource/api/internal/logic/resource"
	"T-driver/app/resource/api/internal/svc"
	"T-driver/app/resource/api/internal/types"
	"T-driver/common/errors"
	"T-driver/common/response"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UploadReq
		if err := httpx.Parse(r, &req); err != nil {
			logx.Error(err)
			response.Response(w, nil, errors.ParamsError(r.Header.Get("Language")))
			return
		}
		file, header, err := r.FormFile("file")
		req.Source, _ = strconv.ParseInt(r.FormValue("source"), 10, 64)
		req.Pid, _ = strconv.ParseInt(r.FormValue("pid"), 10, 64)
		req.TransitType, _ = strconv.ParseInt(r.FormValue("transit_type"), 10, 64)
		if err != nil {
			logx.Error(err)
			response.Response(w, nil, errors.ParamsError(r.Header.Get("Language")))
			return
		}
		defer file.Close()
		req.AssetName = header.Filename
		req.AssetSize = header.Size
		bytes, err := io.ReadAll(file)
		if err != nil {
			logx.Error(err)
			response.Response(w, nil, errors.CustomError("读取文件失败"))
			return
		}
		req.Flie = bytes
		l := resource.NewUploadLogic(r.Context(), svcCtx)
		resp, err := l.Upload(&req)
		response.Response(w, resp, err)
	}
}
