package {{.PkgName}}

import (
	"net/http"

    "T-driver/common/errors"
	"T-driver/common/response"
	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/go-zero/core/logx"
	{{.ImportPackages}}
)

func {{.HandlerName}}(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		{{if .HasRequest}}var req types.{{.RequestType}}
		if err := httpx.Parse(r, &req); err != nil {
			logx.Error(err)
			response.Response(w, nil, errors.ParamsError(r.Header.Get("Language")))
			return
		}

		{{end}}l := {{.LogicName}}.New{{.LogicType}}(r.Context(), svcCtx)
		{{if .HasResp}}resp, {{end}}err := l.{{.Call}}({{if .HasRequest}}&req{{end}})
		{{if .HasResp}}response.Response(w, resp, err){{else}}response.Response(w, nil, err){{end}}
	}
}
