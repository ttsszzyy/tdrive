package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ BotCommandModel = (*customBotCommandModel)(nil)

type (
	// BotCommandModel is an interface to be customized, add more methods here,
	// and implement the added methods in customBotCommandModel.
	BotCommandModel interface {
		botCommandModel
	}

	customBotCommandModel struct {
		*defaultBotCommandModel
	}
)

// NewBotCommandModel returns a model for the database table.
func NewBotCommandModel(conn sqlx.SqlConn, c cache.CacheConf) BotCommandModel {
	return &customBotCommandModel{
		defaultBotCommandModel: newBotCommandModel(conn, c),
	}
}
