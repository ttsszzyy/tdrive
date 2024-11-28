package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ DictModel = (*customDictModel)(nil)

type (
	// DictModel is an interface to be customized, add more methods here,
	// and implement the added methods in customDictModel.
	DictModel interface {
		dictModel
	}

	customDictModel struct {
		*defaultDictModel
	}
)

// NewDictModel returns a model for the database table.
func NewDictModel(conn sqlx.SqlConn, c cache.CacheConf) DictModel {
	return &customDictModel{
		defaultDictModel: newDictModel(conn, c),
	}
}
