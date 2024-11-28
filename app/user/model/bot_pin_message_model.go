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

var _ BotPinMessageModel = (*customBotPinMessageModel)(nil)

type (
	// BotPinMessageModel is an interface to be customized, add more methods here,
	// and implement the added methods in customBotPinMessageModel.
	BotPinMessageModel interface {
		botPinMessageModel
		Deletes(ctx context.Context, ids []int64, s ...sqlx.Session) error
	}

	customBotPinMessageModel struct {
		*defaultBotPinMessageModel
	}
)

// NewBotPinMessageModel returns a model for the database table.
func NewBotPinMessageModel(conn sqlx.SqlConn, c cache.CacheConf) BotPinMessageModel {
	return &customBotPinMessageModel{
		defaultBotPinMessageModel: newBotPinMessageModel(conn, c),
	}
}

func (m *defaultBotPinMessageModel) Deletes(ctx context.Context, ids []int64, s ...sqlx.Session) error {
	shareIdKeys := make([]string, 0)
	for _, id := range ids {
		shareIdKey := fmt.Sprintf("%s%v", cacheBotPinMessageIdPrefix, id)
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
