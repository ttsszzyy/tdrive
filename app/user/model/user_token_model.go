package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserTokenModel = (*customUserTokenModel)(nil)

type (
	// UserTokenModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserTokenModel.
	UserTokenModel interface {
		userTokenModel
	}

	customUserTokenModel struct {
		*defaultUserTokenModel
	}
)

// NewUserTokenModel returns a model for the database table.
func NewUserTokenModel(conn sqlx.SqlConn, c cache.CacheConf) UserTokenModel {
	return &customUserTokenModel{
		defaultUserTokenModel: newUserTokenModel(conn, c),
	}
}
