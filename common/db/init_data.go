package db

import (
	"context"
	initdata "github.com/telegram-mini-apps/init-data-golang"
)

type contextKey string

const (
	_initDataKey contextKey = "init-data"
)

// Returns new context with specified init data.
func WithInitData(ctx context.Context, initData initdata.InitData) context.Context {
	return context.WithValue(ctx, _initDataKey, initData)
}

// Returns the init data from the specified context.
func CtxInitData(ctx context.Context) (initdata.InitData, bool) {
	initData, ok := ctx.Value(_initDataKey).(initdata.InitData)
	return initData, ok
}
