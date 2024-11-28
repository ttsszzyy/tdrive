package middleware

import (
	"T-driver/common/errors"
	"T-driver/common/lib/jwt"
	"T-driver/common/response"
	"T-driver/common/utils"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminAuthMiddleware struct {
	Jwt jwt.JWT
}

func NewAdminAuthMiddleware(jwt jwt.JWT) *AdminAuthMiddleware {
	return &AdminAuthMiddleware{
		Jwt: jwt,
	}
}

func (m *AdminAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		token = strings.TrimPrefix(token, "Bearer ")
		payload, err := m.Jwt.Parse(token)
		if err != nil {
			logx.Error(err)
			response.Response(w, nil, errors.UnauthError(r.Header.Get("Language")))
			return
		}
		id, ok := payload["id"]
		if !ok {
			response.Response(w, nil, errors.UnauthError(r.Header.Get("Language")))
			return
		}
		account, ok := payload["account"]
		if !ok {
			response.Response(w, nil, errors.UnauthError(r.Header.Get("Language")))
			return
		}
		i, err := strconv.ParseInt(fmt.Sprintf("%v", id), 10, 64)
		if err != nil {
			response.Response(w, nil, errors.UnauthError(r.Header.Get("Language")))
			return
		}
		header := utils.GetHeaderFromRequest(r)
		r = r.WithContext(context.WithValue(r.Context(), "ip", header.IP))
		r = r.WithContext(context.WithValue(r.Context(), "os", header.OS))
		r = r.WithContext(context.WithValue(r.Context(), "broswer", header.Broswer))
		r = r.WithContext(context.WithValue(r.Context(), "id", i))
		r = r.WithContext(context.WithValue(r.Context(), "account", fmt.Sprintf("%s", account)))
		r = r.WithContext(context.WithValue(r.Context(), "language", r.Header.Get("Language")))
		next(w, r)
	}
}
