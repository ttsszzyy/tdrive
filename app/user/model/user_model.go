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

var _ UserModel = (*customUserModel)(nil)

const oldct int64 = 1728898656

type (
	// UserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserModel.
	UserModel interface {
		userModel
		CountExchangeStorage(ctx context.Context, cbs ...squirrel.SelectBuilder) (int64, error)
		CountInvite(ctx context.Context, cbs ...squirrel.SelectBuilder) (int64, error)
		CheckIsOldUer(ctx context.Context, uid int64) (bool, error)
	}

	customUserModel struct {
		*defaultUserModel
	}
)

// NewUserModel returns a model for the database table.
func NewUserModel(conn sqlx.SqlConn, c cache.CacheConf) UserModel {
	return &customUserModel{
		defaultUserModel: newUserModel(conn, c),
	}
}

func (m *defaultUserModel) CountExchangeStorage(ctx context.Context, cbs ...squirrel.SelectBuilder) (int64, error) {
	var cb squirrel.SelectBuilder
	if cbs == nil {
		cb = squirrel.Select()
	} else {
		cb = cbs[0]
	}

	// count builder
	query, args, err := cb.Column("COALESCE(sum(exchange_storage),0) as count").From(m.table).ToSql()
	if err != nil {
		return 0, err
	}

	var total int64
	if err := m.QueryRowNoCacheCtx(ctx, &total, query, args...); err != nil {
		return 0, err
	}

	return total, err
}

func (m *defaultUserModel) CountInvite(ctx context.Context, cbs ...squirrel.SelectBuilder) (int64, error) {
	var cb squirrel.SelectBuilder
	if cbs == nil {
		cb = squirrel.Select()
	} else {
		cb = cbs[0]
	}

	// count builder
	query, args, err := cb.Column("Count(DISTINCT uid) as count").From(m.table).ToSql()
	if err != nil {
		return 0, err
	}

	var total int64
	if err := m.QueryRowNoCacheCtx(ctx, &total, query, args...); err != nil {
		return 0, err
	}

	return total, err
}

func (m *defaultUserModel) Deletes(ctx context.Context, ids []int64, s ...sqlx.Session) error {
	dataList, err := m.List(ctx, squirrel.Select().Where(squirrel.Eq{"id": ids}))
	if err != nil {
		return err
	}
	userIdKeys := make([]string, 0, len(dataList)*3)
	for _, v := range dataList {
		userIdKeys = append(userIdKeys, fmt.Sprintf("%s%v", cacheUserIdPrefix, v.Id), fmt.Sprintf("%s%v:%v", cacheUserRecommendCodeDeletedTimePrefix, v.RecommendCode, v.DeletedTime), fmt.Sprintf("%s%v:%v", cacheUserUidDeletedTimePrefix, v.Uid, v.DeletedTime))
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
	}, userIdKeys...)
	return err
}

func (m *defaultUserModel) CheckIsOldUer(ctx context.Context, uid int64) (bool, error) {
	var ct int64

	query, args, err := squirrel.Select("created_time").From(m.table).Where("uid = ?", uid).OrderBy("created_time ASC").Limit(1).ToSql()
	if err != nil {
		return false, fmt.Errorf("generate sql of check is old user error:%w", err)
	}

	err = m.QueryRowNoCacheCtx(ctx, &ct, query, args...)
	switch err {
	case nil:
		return ct <= oldct, nil
	case sqlc.ErrNotFound:
		return false, nil
	default:
		return false, fmt.Errorf("")
	}
}
