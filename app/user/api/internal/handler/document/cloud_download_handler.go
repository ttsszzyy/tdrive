package document

import (
	"fmt"
	"net/http"
	"net/url"

	"T-driver/app/user/api/internal/logic/document"
	"T-driver/app/user/api/internal/svc"
	"T-driver/app/user/api/internal/types"
	"T-driver/common/errors"
	"T-driver/common/response"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func CloudDownloadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CloudDownloadReq
		if err := httpx.Parse(r, &req); err != nil {
			logx.Error(err)
			response.Response(w, nil, errors.ParamsError(r.Header.Get("Language")))
			return
		}

		l := document.NewCloudDownloadLogic(r.Context(), svcCtx)
		resp, err := l.CloudDownload(&req)
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", url.PathEscape(resp.AssetName)))
		w.Header().Set("Content-Transfer-Encoding", "binary")
		w.Header().Set("Content-Type", "application/octet-stream;charset=utf8")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", resp.AssetSize))
		_, err = w.Write(resp.Flie)
		response.ResponseBlob(w, nil, err)
	}
}
