package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ShareModel = (*customShareModel)(nil)

type (
	// ShareModel is an interface to be customized, add more methods here,
	// and implement the added methods in customShareModel.
	ShareModel interface {
		shareModel
		Deletes(ctx context.Context, ids []int64, s ...sqlx.Session) error
		FindOneByUUID(ctx context.Context, uuid string) (*Share, error)
	}

	customShareModel struct {
		*defaultShareModel
	}
)

// NewShareModel returns a model for the database table.
func NewShareModel(conn sqlx.SqlConn, c cache.CacheConf) ShareModel {
	return &customShareModel{
		defaultShareModel: newShareModel(conn, c),
	}
}

func (m *defaultShareModel) Deletes(ctx context.Context, ids []int64, s ...sqlx.Session) error {
	shareIdKeys := make([]string, 0)
	for _, id := range ids {
		shareIdKey := fmt.Sprintf("%s%v", cacheShareIdPrefix, id)
		shareIdKeys = append(shareIdKeys, shareIdKey)
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
	}, shareIdKeys...)
	return err
}

func (m *defaultShareModel) FindOneByUUID(ctx context.Context, uuid string) (*Share, error) {
	var resp Share

	query := fmt.Sprintf("select %s from %s where `uuid` = ? limit 1", shareRows, m.table)
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, uuid)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
