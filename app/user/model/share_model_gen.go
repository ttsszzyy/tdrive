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
	shareFieldNames          = builder.RawFieldNames(&Share{})
	shareRows                = strings.Join(shareFieldNames, ",")
	shareRowsExpectAutoSet   = strings.Join(stringx.Remove(shareFieldNames, "`id`", "`create_time`", "`update_time`", "`create_at`", "`update_at`"), ",")
	shareRowsWithPlaceHolder = strings.Join(stringx.Remove(shareFieldNames, "`id`", "`create_time`", "`update_time`", "`create_at`", "`update_at`"), "=?,") + "=?"

	cacheShareIdPrefix              = "cache:share:id:"
	cacheShareUuidDeletedTimePrefix = "cache:share:uuid:deletedTime:"
)

type (
	shareModel interface {
		Insert(ctx context.Context, data *Share, s ...sqlx.Session) (sql.Result, error)
		FindOne(ctx context.Context, id int64, s ...sqlx.Session) (*Share, error)
		FindOneByUuidDeletedTime(ctx context.Context, uuid string, deletedTime int64) (*Share, error)
		FindOneByBuilder(ctx context.Context, sbs ...squirrel.SelectBuilder) (*Share, error)
		Update(ctx context.Context, data *Share, s ...sqlx.Session) error
		ListPage(ctx context.Context, page, size int64, sbs ...squirrel.SelectBuilder) ([]*Share, int64, error)
		List(ctx context.Context, sbs ...squirrel.SelectBuilder) ([]*Share, error)
		Count(ctx context.Context, cbs ...squirrel.SelectBuilder) (int64, error)
		Delete(ctx context.Context, id int64, s ...sqlx.Session) error
	}

	defaultShareModel struct {
		sqlc.CachedConn
		table string
	}

	Share struct {
		Id            int64  `db:"id"`
		Uid           int64  `db:"uid"`
		AssetIds      string `db:"asset_ids"`      // 资源id的集合
		Uuid          string `db:"uuid"`           // 分享的唯一标识
		AssetName     string `db:"asset_name"`     // 文件名
		AssetSize     int64  `db:"asset_size"`     // 文件大小
		AssetType     int64  `db:"asset_type"`     // 2文件3视频4图片
		Password      string `db:"password"`       // 密码
		Link          string `db:"link"`           // 链接
		EffectiveTime int64  `db:"effective_time"` // 有效时间
		ReadNum       int64  `db:"read_num"`       // 阅览数量
		SaveNum       int64  `db:"save_num"`       // 保存数量
		CreatedTime   int64  `db:"created_time"`
		UpdatedTime   int64  `db:"updated_time"`
		DeletedTime   int64  `db:"deleted_time"`
	}
)

func newShareModel(conn sqlx.SqlConn, c cache.CacheConf) *defaultShareModel {
	return &defaultShareModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`share`",
	}
}

func (m *defaultShareModel) Delete(ctx context.Context, id int64, s ...sqlx.Session) error {
	data, err := m.FindOne(ctx, id, s...)
	if err != nil {
		return err
	}

	shareIdKey := fmt.Sprintf("%s%v", cacheShareIdPrefix, id)
	shareUuidDeletedTimeKey := fmt.Sprintf("%s%v:%v", cacheShareUuidDeletedTimePrefix, data.Uuid, data.DeletedTime)

	var query string
	if ShareSoftDelete {
		query = fmt.Sprintf("update %s set `deleted_time` = %v where `id` = ?", m.table, time.Now().Unix())
	} else {
		query = fmt.Sprintf("delete from %s where `id` = ?", m.table)
	}

	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		if s != nil {
			return s[0].ExecCtx(ctx, query, id)
		}
		return conn.ExecCtx(ctx, query, id)
	}, shareIdKey, shareUuidDeletedTimeKey)
	return err
}

func (m *defaultShareModel) FindOne(ctx context.Context, id int64, s ...sqlx.Session) (*Share, error) {
	shareIdKey := fmt.Sprintf("%s%v", cacheShareIdPrefix, id)
	var resp Share
	err := m.QueryRowCtx(ctx, &resp, shareIdKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", shareRows, m.table)
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

func (m *defaultShareModel) FindOneByUuidDeletedTime(ctx context.Context, uuid string, deletedTime int64) (*Share, error) {
	shareUuidDeletedTimeKey := fmt.Sprintf("%s%v:%v", cacheShareUuidDeletedTimePrefix, uuid, deletedTime)
	var resp Share
	err := m.QueryRowIndexCtx(ctx, &resp, shareUuidDeletedTimeKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `uuid` = ? and `deleted_time` = ? limit 1", shareRows, m.table)
		if err := conn.QueryRowCtx(ctx, &resp, query, uuid, deletedTime); err != nil {
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

func (m *defaultShareModel) Insert(ctx context.Context, data *Share, s ...sqlx.Session) (sql.Result, error) {
	shareIdKey := fmt.Sprintf("%s%v", cacheShareIdPrefix, data.Id)
	shareUuidDeletedTimeKey := fmt.Sprintf("%s%v:%v", cacheShareUuidDeletedTimePrefix, data.Uuid, data.DeletedTime)
	return m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, shareRowsExpectAutoSet)
		if s != nil {
			return s[0].ExecCtx(ctx, query, data.Uid, data.AssetIds, data.Uuid, data.AssetName, data.AssetSize, data.AssetType, data.Password, data.Link, data.EffectiveTime, data.ReadNum, data.SaveNum, data.CreatedTime, data.UpdatedTime, data.DeletedTime)
		}
		return conn.ExecCtx(ctx, query, data.Uid, data.AssetIds, data.Uuid, data.AssetName, data.AssetSize, data.AssetType, data.Password, data.Link, data.EffectiveTime, data.ReadNum, data.SaveNum, data.CreatedTime, data.UpdatedTime, data.DeletedTime)
	}, shareIdKey, shareUuidDeletedTimeKey)
}

func (m *defaultShareModel) Update(ctx context.Context, newData *Share, s ...sqlx.Session) error {
	data, err := m.FindOne(ctx, newData.Id, s...)
	if err != nil {
		return err
	}

	shareIdKey := fmt.Sprintf("%s%v", cacheShareIdPrefix, data.Id)
	shareUuidDeletedTimeKey := fmt.Sprintf("%s%v:%v", cacheShareUuidDeletedTimePrefix, data.Uuid, data.DeletedTime)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, shareRowsWithPlaceHolder)
		if s != nil {
			return s[0].ExecCtx(ctx, query, newData.Uid, newData.AssetIds, newData.Uuid, newData.AssetName, newData.AssetSize, newData.AssetType, newData.Password, newData.Link, newData.EffectiveTime, newData.ReadNum, newData.SaveNum, newData.CreatedTime, newData.UpdatedTime, newData.DeletedTime, newData.Id)
		}
		return conn.ExecCtx(ctx, query, newData.Uid, newData.AssetIds, newData.Uuid, newData.AssetName, newData.AssetSize, newData.AssetType, newData.Password, newData.Link, newData.EffectiveTime, newData.ReadNum, newData.SaveNum, newData.CreatedTime, newData.UpdatedTime, newData.DeletedTime, newData.Id)
	}, shareIdKey, shareUuidDeletedTimeKey)
	return err
}

func (m *defaultShareModel) FindOneByBuilder(ctx context.Context, sbs ...squirrel.SelectBuilder) (*Share, error) {
	var sb squirrel.SelectBuilder
	if sbs == nil {
		sb = squirrel.Select()
	} else {
		sb = sbs[0]
	}

	// query rows
	query, args, err := sb.From(m.table).Columns(shareRows).Limit(1).ToSql()
	if err != nil {
		return nil, err
	}
	var ret Share
	if err := m.QueryRowNoCacheCtx(ctx, &ret, query, args...); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (m *defaultShareModel) ListPage(ctx context.Context, page, size int64, sbs ...squirrel.SelectBuilder) ([]*Share, int64, error) {
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
	query, args, err = sb.From(m.table).Columns(shareRows).Offset(uint64((page - 1) * size)).Limit(uint64(size)).ToSql()
	if err != nil {
		return nil, 0, err
	}
	var list []*Share
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query, args...); err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (m *defaultShareModel) List(ctx context.Context, sbs ...squirrel.SelectBuilder) ([]*Share, error) {
	var sb squirrel.SelectBuilder
	if sbs == nil {
		sb = squirrel.Select()
	} else {
		sb = sbs[0]
	}

	// query rows
	query, args, err := sb.From(m.table).Columns(shareRows).ToSql()
	if err != nil {
		return nil, err
	}
	var list []*Share
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query, args...); err != nil {
		return nil, err
	}

	return list, nil
}

func (m *defaultShareModel) Count(ctx context.Context, cbs ...squirrel.SelectBuilder) (int64, error) {
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

func (m *defaultShareModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheShareIdPrefix, primary)
}

func (m *defaultShareModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", shareRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultShareModel) tableName() string {
	return m.table
}

var ShareSoftDelete bool

func init() {
	tp := reflect.TypeOf(Share{})
	for i := 0; i < tp.NumField(); i++ {
		if tp.Field(i).Tag.Get("db") == "deleted_time" {
			ShareSoftDelete = true
			return
		}
	}
}
