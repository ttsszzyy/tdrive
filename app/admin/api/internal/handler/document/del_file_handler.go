package document

import (
	"net/http"

	"T-driver/app/admin/api/internal/logic/document"
	"T-driver/app/admin/api/internal/svc"
	"T-driver/app/admin/api/internal/types"
	"T-driver/common/errors"
	"T-driver/common/response"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func DelFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DelFileReq
		if err := httpx.Parse(r, &req); err != nil {
			logx.Error(err)
			response.Response(w, nil, errors.ParamsError(r.Header.Get("Language")))
			return
		}

		l := document.NewDelFileLogic(r.Context(), svcCtx)
		resp, err := l.DelFile(&req)
		response.Response(w, resp, err)
	}
}
