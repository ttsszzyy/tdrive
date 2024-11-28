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

var _ UserInviteModel = (*customUserInviteModel)(nil)

type (
	// UserInviteModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserInviteModel.
	UserInviteModel interface {
		userInviteModel
		Deletes(ctx context.Context, ids []int64, s ...sqlx.Session) error
	}

	customUserInviteModel struct {
		*defaultUserInviteModel
	}
)

// NewUserInviteModel returns a model for the database table.
func NewUserInviteModel(conn sqlx.SqlConn, c cache.CacheConf) UserInviteModel {
	return &customUserInviteModel{
		defaultUserInviteModel: newUserInviteModel(conn, c),
	}
}

func (m *defaultUserInviteModel) Deletes(ctx context.Context, ids []int64, s ...sqlx.Session) error {
	shareIdKeys := make([]string, 0)
	for _, id := range ids {
		shareIdKey := fmt.Sprintf("%s%v", cacheUserInviteIdPrefix, id)
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
