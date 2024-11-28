package document

import (
	"T-driver/app/user/api/internal/logic/document"
	"T-driver/app/user/api/internal/svc"
	"T-driver/app/user/api/internal/types"
	"T-driver/common/errors"
	"T-driver/common/response"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"io"
	"net/http"
	"strconv"
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
		/*req.Pid, _ = strconv.ParseInt(r.FormValue("pid"), 10, 64)
		req.TransitType, _ = strconv.ParseInt(r.FormValue("transit_type"), 10, 64)*/
		if err != nil {
			logx.Error(err)
			response.Response(w, nil, errors.ParamsError(r.Header.Get("Language")))
			return
		}
		defer file.Close()
		req.Id, _ = strconv.ParseInt(r.FormValue("id"), 10, 64)
		req.AssetName = header.Filename
		req.AssetSize = header.Size
		bytes, err := io.ReadAll(file)
		if err != nil {
			logx.Error(err)
			response.Response(w, nil, errors.CustomError("读取文件失败"))
			return
		}
		req.Flie = bytes
		l := document.NewUploadLogic(r.Context(), svcCtx)
		resp, err := l.Upload(&req)
		response.Response(w, resp, err)
	}
}
