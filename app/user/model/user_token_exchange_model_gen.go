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
	userTokenExchangeFieldNames          = builder.RawFieldNames(&UserTokenExchange{})
	userTokenExchangeRows                = strings.Join(userTokenExchangeFieldNames, ",")
	userTokenExchangeRowsExpectAutoSet   = strings.Join(stringx.Remove(userTokenExchangeFieldNames, "`id`", "`create_time`", "`update_time`", "`create_at`", "`update_at`"), ",")
	userTokenExchangeRowsWithPlaceHolder = strings.Join(stringx.Remove(userTokenExchangeFieldNames, "`id`", "`create_time`", "`update_time`", "`create_at`", "`update_at`"), "=?,") + "=?"

	cacheUserTokenExchangeIdPrefix = "cache:userTokenExchange:id:"
)

type (
	userTokenExchangeModel interface {
		Insert(ctx context.Context, data *UserTokenExchange, s ...sqlx.Session) (sql.Result, error)
		FindOne(ctx context.Context, id int64, s ...sqlx.Session) (*UserTokenExchange, error)
		FindOneByBuilder(ctx context.Context, sbs ...squirrel.SelectBuilder) (*UserTokenExchange, error)
		Update(ctx context.Context, data *UserTokenExchange, s ...sqlx.Session) error
		ListPage(ctx context.Context, page, size int64, sbs ...squirrel.SelectBuilder) ([]*UserTokenExchange, int64, error)
		List(ctx context.Context, sbs ...squirrel.SelectBuilder) ([]*UserTokenExchange, error)
		Count(ctx context.Context, cbs ...squirrel.SelectBuilder) (int64, error)
		Delete(ctx context.Context, id int64, s ...sqlx.Session) error
	}

	defaultUserTokenExchangeModel struct {
		sqlc.CachedConn
		table string
	}

	UserTokenExchange struct {
		Id            int64 `db:"id"`
		Uid           int64 `db:"uid"`            // 用户uid
		ExchangeToken int64 `db:"exchange_token"` // 预定的空投
		CreatedTime   int64 `db:"created_time"`
		DeletedTime   int64 `db:"deleted_time"`
	}
)

func newUserTokenExchangeModel(conn sqlx.SqlConn, c cache.CacheConf) *defaultUserTokenExchangeModel {
	return &defaultUserTokenExchangeModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`user_token_exchange`",
	}
}

func (m *defaultUserTokenExchangeModel) Delete(ctx context.Context, id int64, s ...sqlx.Session) error {
	userTokenExchangeIdKey := fmt.Sprintf("%s%v", cacheUserTokenExchangeIdPrefix, id)

	var query string
	if UserTokenExchangeSoftDelete {
		query = fmt.Sprintf("update %s set `deleted_time` = %v where `id` = ?", m.table, time.Now().Unix())
	} else {
		query = fmt.Sprintf("delete from %s where `id` = ?", m.table)
	}

	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		if s != nil {
			return s[0].ExecCtx(ctx, query, id)
		}
		return conn.ExecCtx(ctx, query, id)
	}, userTokenExchangeIdKey)
	return err
}

func (m *defaultUserTokenExchangeModel) FindOne(ctx context.Context, id int64, s ...sqlx.Session) (*UserTokenExchange, error) {
	userTokenExchangeIdKey := fmt.Sprintf("%s%v", cacheUserTokenExchangeIdPrefix, id)
	var resp UserTokenExchange
	err := m.QueryRowCtx(ctx, &resp, userTokenExchangeIdKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", userTokenExchangeRows, m.table)
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

func (m *defaultUserTokenExchangeModel) Insert(ctx context.Context, data *UserTokenExchange, s ...sqlx.Session) (sql.Result, error) {
	userTokenExchangeIdKey := fmt.Sprintf("%s%v", cacheUserTokenExchangeIdPrefix, data.Id)
	return m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?)", m.table, userTokenExchangeRowsExpectAutoSet)
		if s != nil {
			return s[0].ExecCtx(ctx, query, data.Uid, data.ExchangeToken, data.CreatedTime, data.DeletedTime)
		}
		return conn.ExecCtx(ctx, query, data.Uid, data.ExchangeToken, data.CreatedTime, data.DeletedTime)
	}, userTokenExchangeIdKey)
}

func (m *defaultUserTokenExchangeModel) Update(ctx context.Context, data *UserTokenExchange, s ...sqlx.Session) error {
	userTokenExchangeIdKey := fmt.Sprintf("%s%v", cacheUserTokenExchangeIdPrefix, data.Id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, userTokenExchangeRowsWithPlaceHolder)
		if s != nil {
			return s[0].ExecCtx(ctx, query, data.Uid, data.ExchangeToken, data.CreatedTime, data.DeletedTime, data.Id)
		}
		return conn.ExecCtx(ctx, query, data.Uid, data.ExchangeToken, data.CreatedTime, data.DeletedTime, data.Id)
	}, userTokenExchangeIdKey)
	return err
}

func (m *defaultUserTokenExchangeModel) FindOneByBuilder(ctx context.Context, sbs ...squirrel.SelectBuilder) (*UserTokenExchange, error) {
	var sb squirrel.SelectBuilder
	if sbs == nil {
		sb = squirrel.Select()
	} else {
		sb = sbs[0]
	}

	// query rows
	query, args, err := sb.From(m.table).Columns(userTokenExchangeRows).Limit(1).ToSql()
	if err != nil {
		return nil, err
	}
	var ret UserTokenExchange
	if err := m.QueryRowNoCacheCtx(ctx, &ret, query, args...); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (m *defaultUserTokenExchangeModel) ListPage(ctx context.Context, page, size int64, sbs ...squirrel.SelectBuilder) ([]*UserTokenExchange, int64, error) {
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
	query, args, err = sb.From(m.table).Columns(userTokenExchangeRows).Offset(uint64((page - 1) * size)).Limit(uint64(size)).ToSql()
	if err != nil {
		return nil, 0, err
	}
	var list []*UserTokenExchange
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query, args...); err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (m *defaultUserTokenExchangeModel) List(ctx context.Context, sbs ...squirrel.SelectBuilder) ([]*UserTokenExchange, error) {
	var sb squirrel.SelectBuilder
	if sbs == nil {
		sb = squirrel.Select()
	} else {
		sb = sbs[0]
	}

	// query rows
	query, args, err := sb.From(m.table).Columns(userTokenExchangeRows).ToSql()
	if err != nil {
		return nil, err
	}
	var list []*UserTokenExchange
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query, args...); err != nil {
		return nil, err
	}

	return list, nil
}

func (m *defaultUserTokenExchangeModel) Count(ctx context.Context, cbs ...squirrel.SelectBuilder) (int64, error) {
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

func (m *defaultUserTokenExchangeModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheUserTokenExchangeIdPrefix, primary)
}

func (m *defaultUserTokenExchangeModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", userTokenExchangeRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultUserTokenExchangeModel) tableName() string {
	return m.table
}

var UserTokenExchangeSoftDelete bool

func init() {
	tp := reflect.TypeOf(UserTokenExchange{})
	for i := 0; i < tp.NumField(); i++ {
		if tp.Field(i).Tag.Get("db") == "deleted_time" {
			UserTokenExchangeSoftDelete = true
			return
		}
	}
}
