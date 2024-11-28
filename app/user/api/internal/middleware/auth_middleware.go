package middleware

import (
	"T-driver/common/db"
	"T-driver/common/errors"
	"T-driver/common/response"
	"T-driver/common/utils"
	"context"
	"net/http"
	"strings"
	"time"

	initdata "github.com/telegram-mini-apps/init-data-golang"
	"github.com/zeromicro/go-zero/core/logx"
)

type AuthMiddleware struct {
	token string
}

func NewAuthMiddleware(token string) *AuthMiddleware {
	return &AuthMiddleware{
		token: token,
	}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取token
		token := r.Header.Get("authorization")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		authParts := strings.Split(r.Header.Get("authorization"), " ")
		if len(authParts) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		authType := authParts[0]
		authData := authParts[1]
		switch authType {
		case "tma":
			// Validate init data. We consider init data sign valid for 1 hour from their
			// creation moment.
			if err := initdata.Validate(authData, m.token, 12*time.Hour); err != nil {
				logx.Errorf("init data validation failed: %v", err)
				//response.Response(w, nil, errors.NewErrCodeMsg(401, err.Error()))
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// Parse init data. We will surely need it in the future.
			initData, err := initdata.Parse(authData)
			if err != nil {
				response.Response(w, nil, errors.SystemError(r.Header.Get("Language")))
				return
			}
			header := utils.GetHeaderFromRequest(r)
			//logx.Errorw("initData:", logx.Field("initData", initData))
			r = r.WithContext(db.WithInitData(r.Context(), initData))
			r = r.WithContext(context.WithValue(r.Context(), "ip", header.IP))
			r = r.WithContext(context.WithValue(r.Context(), "language", r.Header.Get("Language")))
		}
		// Passthrough to next handler if need
		next(w, r)
	}
}
