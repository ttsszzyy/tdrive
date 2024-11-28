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

var _ UserInviteRewardModel = (*customUserInviteRewardModel)(nil)

type (
	// UserInviteRewardModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserInviteRewardModel.
	UserInviteRewardModel interface {
		userInviteRewardModel
		Deletes(ctx context.Context, ids []int64, uid int64, s ...sqlx.Session) error
	}

	customUserInviteRewardModel struct {
		*defaultUserInviteRewardModel
	}
)

// NewUserInviteRewardModel returns a model for the database table.
func NewUserInviteRewardModel(conn sqlx.SqlConn, c cache.CacheConf) UserInviteRewardModel {
	return &customUserInviteRewardModel{
		defaultUserInviteRewardModel: newUserInviteRewardModel(conn, c),
	}
}

func (m *defaultUserInviteRewardModel) Deletes(ctx context.Context, ids []int64, uid int64, s ...sqlx.Session) error {
	shareIdKeys := make([]string, 0, len(ids)+1)
	shareIdKeys = append(shareIdKeys, fmt.Sprintf("%s%v:%v", cacheUserInviteRewardUidDeletedTimePrefix, uid, 0))
	for _, id := range ids {
		shareIdKey := fmt.Sprintf("%s%v", cacheUserInviteRewardIdPrefix, id)
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
