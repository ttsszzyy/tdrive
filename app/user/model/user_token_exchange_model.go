package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserTokenExchangeModel = (*customUserTokenExchangeModel)(nil)

type (
	// UserTokenExchangeModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserTokenExchangeModel.
	UserTokenExchangeModel interface {
		userTokenExchangeModel
	}

	customUserTokenExchangeModel struct {
		*defaultUserTokenExchangeModel
	}
)

// NewUserTokenExchangeModel returns a model for the database table.
func NewUserTokenExchangeModel(conn sqlx.SqlConn, c cache.CacheConf) UserTokenExchangeModel {
	return &customUserTokenExchangeModel{
		defaultUserTokenExchangeModel: newUserTokenExchangeModel(conn, c),
	}
}
