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
	userTokenFieldNames          = builder.RawFieldNames(&UserToken{})
	userTokenRows                = strings.Join(userTokenFieldNames, ",")
	userTokenRowsExpectAutoSet   = strings.Join(stringx.Remove(userTokenFieldNames, "`id`", "`create_time`", "`update_time`", "`create_at`", "`update_at`"), ",")
	userTokenRowsWithPlaceHolder = strings.Join(stringx.Remove(userTokenFieldNames, "`id`", "`create_time`", "`update_time`", "`create_at`", "`update_at`"), "=?,") + "=?"

	cacheUserTokenIdPrefix             = "cache:userToken:id:"
	cacheUserTokenUidDeletedTimePrefix = "cache:userToken:uid:deletedTime:"
)

type (
	userTokenModel interface {
		Insert(ctx context.Context, data *UserToken, s ...sqlx.Session) (sql.Result, error)
		FindOne(ctx context.Context, id int64, s ...sqlx.Session) (*UserToken, error)
		FindOneByUidDeletedTime(ctx context.Context, uid int64, deletedTime int64) (*UserToken, error)
		FindOneByBuilder(ctx context.Context, sbs ...squirrel.SelectBuilder) (*UserToken, error)
		Update(ctx context.Context, data *UserToken, s ...sqlx.Session) error
		ListPage(ctx context.Context, page, size int64, sbs ...squirrel.SelectBuilder) ([]*UserToken, int64, error)
		List(ctx context.Context, sbs ...squirrel.SelectBuilder) ([]*UserToken, error)
		Count(ctx context.Context, cbs ...squirrel.SelectBuilder) (int64, error)
		Delete(ctx context.Context, id int64, s ...sqlx.Session) error
	}

	defaultUserTokenModel struct {
		sqlc.CachedConn
		table string
	}

	UserToken struct {
		Id          int64 `db:"id"`
		Uid         int64 `db:"uid"`   // 用户uid
		Token       int64 `db:"token"` // 预定的空投
		CreatedTime int64 `db:"created_time"`
		UpdatedTime int64 `db:"updated_time"`
		DeletedTime int64 `db:"deleted_time"`
	}
)

func newUserTokenModel(conn sqlx.SqlConn, c cache.CacheConf) *defaultUserTokenModel {
	return &defaultUserTokenModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`user_token`",
	}
}

func (m *defaultUserTokenModel) Delete(ctx context.Context, id int64, s ...sqlx.Session) error {
	data, err := m.FindOne(ctx, id, s...)
	if err != nil {
		return err
	}

	userTokenIdKey := fmt.Sprintf("%s%v", cacheUserTokenIdPrefix, id)
	userTokenUidDeletedTimeKey := fmt.Sprintf("%s%v:%v", cacheUserTokenUidDeletedTimePrefix, data.Uid, data.DeletedTime)

	var query string
	if UserTokenSoftDelete {
		query = fmt.Sprintf("update %s set `deleted_time` = %v where `id` = ?", m.table, time.Now().Unix())
	} else {
		query = fmt.Sprintf("delete from %s where `id` = ?", m.table)
	}

	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		if s != nil {
			return s[0].ExecCtx(ctx, query, id)
		}
		return conn.ExecCtx(ctx, query, id)
	}, userTokenIdKey, userTokenUidDeletedTimeKey)
	return err
}

func (m *defaultUserTokenModel) FindOne(ctx context.Context, id int64, s ...sqlx.Session) (*UserToken, error) {
	userTokenIdKey := fmt.Sprintf("%s%v", cacheUserTokenIdPrefix, id)
	var resp UserToken
	err := m.QueryRowCtx(ctx, &resp, userTokenIdKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", userTokenRows, m.table)
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

func (m *defaultUserTokenModel) FindOneByUidDeletedTime(ctx context.Context, uid int64, deletedTime int64) (*UserToken, error) {
	userTokenUidDeletedTimeKey := fmt.Sprintf("%s%v:%v", cacheUserTokenUidDeletedTimePrefix, uid, deletedTime)
	var resp UserToken
	err := m.QueryRowIndexCtx(ctx, &resp, userTokenUidDeletedTimeKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `uid` = ? and `deleted_time` = ? limit 1", userTokenRows, m.table)
		if err := conn.QueryRowCtx(ctx, &resp, query, uid, deletedTime); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, m.queryPrimary)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserTokenModel) Insert(ctx context.Context, data *UserToken, s ...sqlx.Session) (sql.Result, error) {
	userTokenIdKey := fmt.Sprintf("%s%v", cacheUserTokenIdPrefix, data.Id)
	userTokenUidDeletedTimeKey := fmt.Sprintf("%s%v:%v", cacheUserTokenUidDeletedTimePrefix, data.Uid, data.DeletedTime)
	return m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?)", m.table, userTokenRowsExpectAutoSet)
		if s != nil {
			return s[0].ExecCtx(ctx, query, data.Uid, data.Token, data.CreatedTime, data.UpdatedTime, data.DeletedTime)
		}
		return conn.ExecCtx(ctx, query, data.Uid, data.Token, data.CreatedTime, data.UpdatedTime, data.DeletedTime)
	}, userTokenIdKey, userTokenUidDeletedTimeKey)
}

func (m *defaultUserTokenModel) Update(ctx context.Context, newData *UserToken, s ...sqlx.Session) error {
	data, err := m.FindOne(ctx, newData.Id, s...)
	if err != nil {
		return err
	}

	userTokenIdKey := fmt.Sprintf("%s%v", cacheUserTokenIdPrefix, data.Id)
	userTokenUidDeletedTimeKey := fmt.Sprintf("%s%v:%v", cacheUserTokenUidDeletedTimePrefix, data.Uid, data.DeletedTime)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, userTokenRowsWithPlaceHolder)
		if s != nil {
			return s[0].ExecCtx(ctx, query, newData.Uid, newData.Token, newData.CreatedTime, newData.UpdatedTime, newData.DeletedTime, newData.Id)
		}
		return conn.ExecCtx(ctx, query, newData.Uid, newData.Token, newData.CreatedTime, newData.UpdatedTime, newData.DeletedTime, newData.Id)
	}, userTokenIdKey, userTokenUidDeletedTimeKey)
	return err
}

func (m *defaultUserTokenModel) FindOneByBuilder(ctx context.Context, sbs ...squirrel.SelectBuilder) (*UserToken, error) {
	var sb squirrel.SelectBuilder
	if sbs == nil {
		sb = squirrel.Select()
	} else {
		sb = sbs[0]
	}

	// query rows
	query, args, err := sb.From(m.table).Columns(userTokenRows).Limit(1).ToSql()
	if err != nil {
		return nil, err
	}
	var ret UserToken
	if err := m.QueryRowNoCacheCtx(ctx, &ret, query, args...); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (m *defaultUserTokenModel) ListPage(ctx context.Context, page, size int64, sbs ...squirrel.SelectBuilder) ([]*UserToken, int64, error) {
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
	query, args, err = sb.From(m.table).Columns(userTokenRows).Offset(uint64((page - 1) * size)).Limit(uint64(size)).ToSql()
	if err != nil {
		return nil, 0, err
	}
	var list []*UserToken
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query, args...); err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (m *defaultUserTokenModel) List(ctx context.Context, sbs ...squirrel.SelectBuilder) ([]*UserToken, error) {
	var sb squirrel.SelectBuilder
	if sbs == nil {
		sb = squirrel.Select()
	} else {
		sb = sbs[0]
	}

	// query rows
	query, args, err := sb.From(m.table).Columns(userTokenRows).ToSql()
	if err != nil {
		return nil, err
	}
	var list []*UserToken
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query, args...); err != nil {
		return nil, err
	}

	return list, nil
}

func (m *defaultUserTokenModel) Count(ctx context.Context, cbs ...squirrel.SelectBuilder) (int64, error) {
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

func (m *defaultUserTokenModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheUserTokenIdPrefix, primary)
}

func (m *defaultUserTokenModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", userTokenRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultUserTokenModel) tableName() string {
	return m.table
}

var UserTokenSoftDelete bool

func init() {
	tp := reflect.TypeOf(UserToken{})
	for i := 0; i < tp.NumField(); i++ {
		if tp.Field(i).Tag.Get("db") == "deleted_time" {
			UserTokenSoftDelete = true
			return
		}
	}
}