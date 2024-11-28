package model

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"time"
)

var _ UserStorageExchangeModel = (*customUserStorageExchangeModel)(nil)

type (
	// UserStorageExchangeModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserStorageExchangeModel.
	UserStorageExchangeModel interface {
		userStorageExchangeModel
		Deletes(ctx context.Context, ids []int64, s ...sqlx.Session) error
	}

	customUserStorageExchangeModel struct {
		*defaultUserStorageExchangeModel
	}
)

// NewUserStorageExchangeModel returns a model for the database table.
func NewUserStorageExchangeModel(conn sqlx.SqlConn, c cache.CacheConf) UserStorageExchangeModel {
	return &customUserStorageExchangeModel{
		defaultUserStorageExchangeModel: newUserStorageExchangeModel(conn, c),
	}
}

func (m *defaultUserStorageExchangeModel) Deletes(ctx context.Context, ids []int64, s ...sqlx.Session) error {
	assetsIdKeys := make([]string, 0)
	for _, id := range ids {
		assetsIdKey := fmt.Sprintf("%s%v", cacheUserStorageExchangeIdPrefix, id)
		assetsIdKeys = append(assetsIdKeys, assetsIdKey)
	}
	query, args, err := squirrel.Update(m.table).Set("deleted_time", time.Now().Unix()).Where(squirrel.Eq{"id": ids}).ToSql()
	if err != nil {
		return err
	}

	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		if s != nil {
			return s[0].ExecCtx(ctx, query, args...)
		}
		return conn.ExecCtx(ctx, query, args...)
	}, assetsIdKeys...)
	return err
}
