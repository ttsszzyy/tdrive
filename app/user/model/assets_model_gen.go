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
	assetsFieldNames          = builder.RawFieldNames(&Assets{})
	assetsRows                = strings.Join(assetsFieldNames, ",")
	assetsRowsExpectAutoSet   = strings.Join(stringx.Remove(assetsFieldNames, "`id`", "`create_time`", "`update_time`", "`create_at`", "`update_at`"), ",")
	assetsRowsWithPlaceHolder = strings.Join(stringx.Remove(assetsFieldNames, "`id`", "`create_time`", "`update_time`", "`create_at`", "`update_at`"), "=?,") + "=?"

	cacheAssetsIdPrefix = "cache:assets:id:"
)

type (
	assetsModel interface {
		Insert(ctx context.Context, data *Assets, s ...sqlx.Session) (sql.Result, error)
		FindOne(ctx context.Context, id int64, s ...sqlx.Session) (*Assets, error)
		FindOneByBuilder(ctx context.Context, sbs ...squirrel.SelectBuilder) (*Assets, error)
		Update(ctx context.Context, data *Assets, s ...sqlx.Session) error
		ListPage(ctx context.Context, page, size int64, sbs ...squirrel.SelectBuilder) ([]*Assets, int64, error)
		List(ctx context.Context, sbs ...squirrel.SelectBuilder) ([]*Assets, error)
		Count(ctx context.Context, cbs ...squirrel.SelectBuilder) (int64, error)
		Delete(ctx context.Context, id int64, s ...sqlx.Session) error
	}

	defaultAssetsModel struct {
		sqlc.CachedConn
		table string
	}

	Assets struct {
		Id          int64  `db:"id"`
		Uid         int64  `db:"uid"`
		Cid         string `db:"cid"`          // 资源id
		AssetName   string `db:"asset_name"`   // 文件名
		AssetSize   int64  `db:"asset_size"`   // 文件大小
		AssetType   int64  `db:"asset_type"`   // 1文件夹2文件3视频4图片
		TransitType int64  `db:"transit_type"` // 上传类型 1云链接 2TG 3X 4TK 5种子
		IsTag       int64  `db:"is_tag"`       // 是否标记 1是2否
		Pid         int64  `db:"pid"`          // 资源所属父资源ID
		Source      int64  `db:"source"`       // 来源 1本地上次2云上传3TDriver
		IsReport    int64  `db:"is_report"`    // 是否举报 1是2否
		ReportType  int64  `db:"report_type"`  // 举报类型 1色情2恐怖3暴力4虐待
		Status      int64  `db:"status"`       // 状态 1禁用2进行中3完成4失败
		Link        string `db:"link"`         // 链接
		IsDefault   int64  `db:"is_default"`   // 是否默认 1是0否
		CreatedTime int64  `db:"created_time"`
		UpdatedTime int64  `db:"updated_time"`
		DeletedTime int64  `db:"deleted_time"`
	}
)

func newAssetsModel(conn sqlx.SqlConn, c cache.CacheConf) *defaultAssetsModel {
	return &defaultAssetsModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`assets`",
	}
}

func (m *defaultAssetsModel) Delete(ctx context.Context, id int64, s ...sqlx.Session) error {
	assetsIdKey := fmt.Sprintf("%s%v", cacheAssetsIdPrefix, id)

	var query string
	if AssetsSoftDelete {
		query = fmt.Sprintf("update %s set `deleted_time` = %v where `id` = ?", m.table, time.Now().Unix())
	} else {
		query = fmt.Sprintf("delete from %s where `id` = ?", m.table)
	}

	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		if s != nil {
			return s[0].ExecCtx(ctx, query, id)
		}
		return conn.ExecCtx(ctx, query, id)
	}, assetsIdKey)
	return err
}

func (m *defaultAssetsModel) FindOne(ctx context.Context, id int64, s ...sqlx.Session) (*Assets, error) {
	assetsIdKey := fmt.Sprintf("%s%v", cacheAssetsIdPrefix, id)
	var resp Assets
	err := m.QueryRowCtx(ctx, &resp, assetsIdKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", assetsRows, m.table)
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

func (m *defaultAssetsModel) Insert(ctx context.Context, data *Assets, s ...sqlx.Session) (sql.Result, error) {
	assetsIdKey := fmt.Sprintf("%s%v", cacheAssetsIdPrefix, data.Id)
	return m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, assetsRowsExpectAutoSet)
		if s != nil {
			return s[0].ExecCtx(ctx, query, data.Uid, data.Cid, data.AssetName, data.AssetSize, data.AssetType, data.TransitType, data.IsTag, data.Pid, data.Source, data.IsReport, data.ReportType, data.Status, data.Link, data.IsDefault, data.CreatedTime, data.UpdatedTime, data.DeletedTime)
		}
		return conn.ExecCtx(ctx, query, data.Uid, data.Cid, data.AssetName, data.AssetSize, data.AssetType, data.TransitType, data.IsTag, data.Pid, data.Source, data.IsReport, data.ReportType, data.Status, data.Link, data.IsDefault, data.CreatedTime, data.UpdatedTime, data.DeletedTime)
	}, assetsIdKey)
}

func (m *defaultAssetsModel) Update(ctx context.Context, data *Assets, s ...sqlx.Session) error {
	assetsIdKey := fmt.Sprintf("%s%v", cacheAssetsIdPrefix, data.Id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, assetsRowsWithPlaceHolder)
		if s != nil {
			return s[0].ExecCtx(ctx, query, data.Uid, data.Cid, data.AssetName, data.AssetSize, data.AssetType, data.TransitType, data.IsTag, data.Pid, data.Source, data.IsReport, data.ReportType, data.Status, data.Link, data.IsDefault, data.CreatedTime, data.UpdatedTime, data.DeletedTime, data.Id)
		}
		return conn.ExecCtx(ctx, query, data.Uid, data.Cid, data.AssetName, data.AssetSize, data.AssetType, data.TransitType, data.IsTag, data.Pid, data.Source, data.IsReport, data.ReportType, data.Status, data.Link, data.IsDefault, data.CreatedTime, data.UpdatedTime, data.DeletedTime, data.Id)
	}, assetsIdKey)
	return err
}

func (m *defaultAssetsModel) FindOneByBuilder(ctx context.Context, sbs ...squirrel.SelectBuilder) (*Assets, error) {
	var sb squirrel.SelectBuilder
	if sbs == nil {
		sb = squirrel.Select()
	} else {
		sb = sbs[0]
	}

	// query rows
	query, args, err := sb.From(m.table).Columns(assetsRows).Limit(1).ToSql()
	if err != nil {
		return nil, err
	}
	var ret Assets
	if err := m.QueryRowNoCacheCtx(ctx, &ret, query, args...); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (m *defaultAssetsModel) ListPage(ctx context.Context, page, size int64, sbs ...squirrel.SelectBuilder) ([]*Assets, int64, error) {
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
	query, args, err = sb.From(m.table).Columns(assetsRows).Offset(uint64((page - 1) * size)).Limit(uint64(size)).ToSql()
	if err != nil {
		return nil, 0, err
	}
	var list []*Assets
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query, args...); err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (m *defaultAssetsModel) List(ctx context.Context, sbs ...squirrel.SelectBuilder) ([]*Assets, error) {
	var sb squirrel.SelectBuilder
	if sbs == nil {
		sb = squirrel.Select()
	} else {
		sb = sbs[0]
	}

	// query rows
	query, args, err := sb.From(m.table).Columns(assetsRows).ToSql()
	if err != nil {
		return nil, err
	}
	var list []*Assets
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query, args...); err != nil {
		return nil, err
	}

	return list, nil
}

func (m *defaultAssetsModel) Count(ctx context.Context, cbs ...squirrel.SelectBuilder) (int64, error) {
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

func (m *defaultAssetsModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheAssetsIdPrefix, primary)
}

func (m *defaultAssetsModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", assetsRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultAssetsModel) tableName() string {
	return m.table
}

var AssetsSoftDelete bool

func init() {
	tp := reflect.TypeOf(Assets{})
	for i := 0; i < tp.NumField(); i++ {
		if tp.Field(i).Tag.Get("db") == "deleted_time" {
			AssetsSoftDelete = true
			return
		}
	}
}
