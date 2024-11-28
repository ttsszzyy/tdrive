// Code generated by goctl. DO NOT EDIT!

package model

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	taskFieldNames          = builder.RawFieldNames(&Task{})
	taskRows                = strings.Join(taskFieldNames, ",")
	taskRowsExpectAutoSet   = strings.Join(stringx.Remove(taskFieldNames, "`id`", "`create_time`", "`update_time`", "`create_at`", "`update_at`"), ",")
	taskRowsWithPlaceHolder = strings.Join(stringx.Remove(taskFieldNames, "`id`", "`create_time`", "`update_time`", "`create_at`", "`update_at`"), "=?,") + "=?"

	cacheTaskIdPrefix = "cache:task:id:"
)

type (
	taskModel interface {
		Insert(ctx context.Context, data *Task, s ...sqlx.Session) (sql.Result, error)
		FindOne(ctx context.Context, id int64, s ...sqlx.Session) (*Task, error)
		FindOneByBuilder(ctx context.Context, sbs ...squirrel.SelectBuilder) (*Task, error)
		Update(ctx context.Context, data *Task, s ...sqlx.Session) error
		ListPage(ctx context.Context, page, size int64, sbs ...squirrel.SelectBuilder) ([]*Task, int64, error)
		List(ctx context.Context, sbs ...squirrel.SelectBuilder) ([]*Task, error)
		Count(ctx context.Context, cbs ...squirrel.SelectBuilder) (int64, error)
		Delete(ctx context.Context, id int64, s ...sqlx.Session) error
	}

	defaultTaskModel struct {
		sqlc.CachedConn
		table string
	}

	Task struct {
		Id          int64 `db:"id"`
		TaskPoolId  int64 `db:"task_pool_id"` // 任务池id
		Uid         int64 `db:"uid"`
		FinishTime  int64 `db:"finish_time"`
		CreatedTime int64 `db:"created_time"`
		UpdatedTime int64 `db:"updated_time"`
		DeletedTime int64 `db:"deleted_time"`
	}
)

func newTaskModel(conn sqlx.SqlConn, c cache.CacheConf) *defaultTaskModel {
	return &defaultTaskModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`task`",
	}
}

func (m *defaultTaskModel) Delete(ctx context.Context, id int64, s ...sqlx.Session) error {
	taskIdKey := fmt.Sprintf("%s%v", cacheTaskIdPrefix, id)

	var query string
	if TaskSoftDelete {
		query = fmt.Sprintf("update %s set `deleted_time` = %v where `id` = ?", m.table, time.Now().Unix())
	} else {
		query = fmt.Sprintf("delete from %s where `id` = ?", m.table)
	}

	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		if s != nil {
			return s[0].ExecCtx(ctx, query, id)
		}
		return conn.ExecCtx(ctx, query, id)
	}, taskIdKey)
	return err
}

func (m *defaultTaskModel) FindOne(ctx context.Context, id int64, s ...sqlx.Session) (*Task, error) {
	taskIdKey := fmt.Sprintf("%s%v", cacheTaskIdPrefix, id)
	var resp Task
	err := m.QueryRowCtx(ctx, &resp, taskIdKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", taskRows, m.table)
		if s != nil {
			return s[0].QueryRowCtx(ctx, v, query, id)
		}
		return conn.QueryRowCtx(ctx, v, query, id)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultTaskModel) Insert(ctx context.Context, data *Task, s ...sqlx.Session) (sql.Result, error) {
	taskIdKey := fmt.Sprintf("%s%v", cacheTaskIdPrefix, data.Id)
	return m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?)", m.table, taskRowsExpectAutoSet)
		if s != nil {
			return s[0].ExecCtx(ctx, query, data.TaskPoolId, data.Uid, data.FinishTime, data.CreatedTime, data.UpdatedTime, data.DeletedTime)
		}
		return conn.ExecCtx(ctx, query, data.TaskPoolId, data.Uid, data.FinishTime, data.CreatedTime, data.UpdatedTime, data.DeletedTime)
	}, taskIdKey)
}

func (m *defaultTaskModel) Update(ctx context.Context, data *Task, s ...sqlx.Session) error {
	taskIdKey := fmt.Sprintf("%s%v", cacheTaskIdPrefix, data.Id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, taskRowsWithPlaceHolder)
		if s != nil {
			return s[0].ExecCtx(ctx, query, data.TaskPoolId, data.Uid, data.FinishTime, data.CreatedTime, data.UpdatedTime, data.DeletedTime, data.Id)
		}
		return conn.ExecCtx(ctx, query, data.TaskPoolId, data.Uid, data.FinishTime, data.CreatedTime, data.UpdatedTime, data.DeletedTime, data.Id)
	}, taskIdKey)
	return err
}

func (m *defaultTaskModel) FindOneByBuilder(ctx context.Context, sbs ...squirrel.SelectBuilder) (*Task, error) {
	var sb squirrel.SelectBuilder
	if sbs == nil {
		sb = squirrel.Select()
	} else {
		sb = sbs[0]
	}

	// query rows
	query, args, err := sb.From(m.table).Columns(taskRows).Limit(1).ToSql()
	if err != nil {
		return nil, err
	}
	var ret Task
	if err := m.QueryRowNoCacheCtx(ctx, &ret, query, args...); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (m *defaultTaskModel) ListPage(ctx context.Context, page, size int64, sbs ...squirrel.SelectBuilder) ([]*Task, int64, error) {
	var sb squirrel.SelectBuilder
	if sbs == nil {
		sb = squirrel.Select()
	} else {
		sb = sbs[0]
	}

	// count builder
	cb := sb.Column("Count(*) as count").From(m.table)
	query, args, err := cb.ToSql()
	if err != nil {
		return nil, 0, err
	}

	var total int64
	if err := m.QueryRowNoCacheCtx(ctx, &total, query, args...); err != nil {
		return nil, 0, err
	}

	// query rows
	if page <= 0 {
		page = 1
	}
	query, args, err = sb.From(m.table).Columns(taskRows).Offset(uint64((page - 1) * size)).Limit(uint64(size)).ToSql()
	if err != nil {
		return nil, 0, err
	}
	var list []*Task
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query, args...); err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (m *defaultTaskModel) List(ctx context.Context, sbs ...squirrel.SelectBuilder) ([]*Task, error) {
	var sb squirrel.SelectBuilder
	if sbs == nil {
		sb = squirrel.Select()
	} else {
		sb = sbs[0]
	}

	// query rows
	query, args, err := sb.From(m.table).Columns(taskRows).ToSql()
	if err != nil {
		return nil, err
	}
	var list []*Task
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query, args...); err != nil {
		return nil, err
	}

	return list, nil
}

func (m *defaultTaskModel) Count(ctx context.Context, cbs ...squirrel.SelectBuilder) (int64, error) {
	var cb squirrel.SelectBuilder
	if cbs == nil {
		cb = squirrel.Select()
	} else {
		cb = cbs[0]
	}

	// count builder
	query, args, err := cb.Column("Count(*) as count").From(m.table).ToSql()
	if err != nil {
		return 0, err
	}

	var total int64
	if err := m.QueryRowNoCacheCtx(ctx, &total, query, args...); err != nil {
		return 0, err
	}

	return total, err
}

func (m *defaultTaskModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheTaskIdPrefix, primary)
}

func (m *defaultTaskModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", taskRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultTaskModel) tableName() string {
	return m.table
}

var TaskSoftDelete bool

func init() {
	tp := reflect.TypeOf(Task{})
	for i := 0; i < tp.NumField(); i++ {
		if tp.Field(i).Tag.Get("db") == "deleted_time" {
			TaskSoftDelete = true
			return
		}
	}
}
