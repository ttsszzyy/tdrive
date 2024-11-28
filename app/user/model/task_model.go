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

var _ TaskModel = (*customTaskModel)(nil)

type (
	// TaskModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTaskModel.
	TaskModel interface {
		taskModel
		UpdateTask(ctx context.Context, id int64, isComplete int64, s ...sqlx.Session) error
		Deletes(ctx context.Context, ids []int64, s ...sqlx.Session) error
	}

	customTaskModel struct {
		*defaultTaskModel
	}
)

// NewTaskModel returns a model for the database table.
func NewTaskModel(conn sqlx.SqlConn, c cache.CacheConf) TaskModel {
	return &customTaskModel{
		defaultTaskModel: newTaskModel(conn, c),
	}
}

func (m *defaultTaskModel) UpdateTask(ctx context.Context, id int64, isComplete int64, s ...sqlx.Session) error {
	data, err := m.FindOne(ctx, id)
	if err != nil {
		return err
	}

	var query string
	var args []interface{}
	sb := squirrel.Update(m.table).Set("updated_time", time.Now().Unix()).Set("is_complete", isComplete).Where(squirrel.Eq{"id": data.Id})
	query, args, err = sb.ToSql()
	if err != nil {
		return err
	}

	taskIdKey := fmt.Sprintf("%s%v", cacheTaskIdPrefix, data.Id)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		if s != nil {
			return s[0].ExecCtx(ctx, query, args...)
		}
		return conn.ExecCtx(ctx, query, args...)
	}, taskIdKey)
	return err
}

func (m *defaultTaskModel) Deletes(ctx context.Context, ids []int64, s ...sqlx.Session) error {
	cacheIdKeys := make([]string, 0, len(ids))
	for _, id := range ids {
		taskIdKey := fmt.Sprintf("%s%v", cacheTaskIdPrefix, id)
		cacheIdKeys = append(cacheIdKeys, taskIdKey)
	}
	query, args, err := squirrel.Delete(m.table).Where(squirrel.Eq{"id": ids}).ToSql()
	if err != nil {
		return err
	}

	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		if s != nil {
			return s[0].ExecCtx(ctx, query, args...)
		}
		return conn.ExecCtx(ctx, query, args...)
	}, cacheIdKeys...)
	return err
}
