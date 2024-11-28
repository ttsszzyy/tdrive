package model

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TaskPoolModel = (*customTaskPoolModel)(nil)

var cacheTaskPoolAllPrefix = "cache:taskPool:all"

type (
	// TaskPoolModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTaskPoolModel.
	TaskPoolModel interface {
		taskPoolModel
		FindAll(ctx context.Context, sbs ...squirrel.SelectBuilder) ([]*TaskPool, error)
		ClearCache() error
	}

	customTaskPoolModel struct {
		*defaultTaskPoolModel
	}
)

// NewTaskPoolModel returns a model for the database table.
func NewTaskPoolModel(conn sqlx.SqlConn, c cache.CacheConf) TaskPoolModel {
	return &customTaskPoolModel{
		defaultTaskPoolModel: newTaskPoolModel(conn, c),
	}
}

// 清理查询缓存
func (m *defaultTaskPoolModel) ClearCache() error {
	err := m.DelCache(cacheTaskPoolAllPrefix)
	if err != nil {
		return err
	}
	return nil
}

func (m *defaultTaskPoolModel) FindAll(ctx context.Context, sbs ...squirrel.SelectBuilder) ([]*TaskPool, error) {
	var sb squirrel.SelectBuilder
	if sbs == nil {
		sb = squirrel.Select()
	} else {
		sb = sbs[0]
	}

	// query rows
	query, args, err := sb.From(m.table).Columns(taskPoolRows).ToSql()
	if err != nil {
		return nil, err
	}

	var list []*TaskPool
	if err := m.QueryRowCtx(ctx, &list, cacheTaskPoolIdPrefix, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		return conn.QueryRowCtx(ctx, v, query, args)
	}); err != nil {
		return nil, err
	}

	return list, nil
}
