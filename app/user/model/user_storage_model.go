package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserStorageModel = (*customUserStorageModel)(nil)

type (
	// UserStorageModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserStorageModel.
	UserStorageModel interface {
		userStorageModel
		UpdateUsedSize(ctx context.Context, uid, size int64, s ...sqlx.Session) error
	}

	customUserStorageModel struct {
		*defaultUserStorageModel
	}
)

// NewUserStorageModel returns a model for the database table.
func NewUserStorageModel(conn sqlx.SqlConn, c cache.CacheConf) UserStorageModel {
	return &customUserStorageModel{
		defaultUserStorageModel: newUserStorageModel(conn, c),
	}
}

// UpdateUsedSize 更新用户使用的存储大小
func (m *customUserStorageModel) UpdateUsedSize(ctx context.Context, uid, size int64, s ...sqlx.Session) error {
	data, err := m.FindOneByUidDeletedTime(ctx, uid, 0)
	if err != nil {
		return err
	}

	userStorageIdKey := fmt.Sprintf("%s%v", cacheUserStorageIdPrefix, data.Id)
	userStorageUidDeletedTimeKey := fmt.Sprintf("%s%v:%v", cacheUserStorageUidDeletedTimePrefix, data.Uid, data.DeletedTime)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set `storage_use` = ? where `id` = ?", m.table)
		if s != nil {
			return s[0].ExecCtx(ctx, query, size, data.Id)
		}
		return conn.ExecCtx(ctx, query, size, data.Id)
	}, userStorageIdKey, userStorageUidDeletedTimeKey)
	return err
}
