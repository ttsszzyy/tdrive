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
	botCommandFieldNames          = builder.RawFieldNames(&BotCommand{})
	botCommandRows                = strings.Join(botCommandFieldNames, ",")
	botCommandRowsExpectAutoSet   = strings.Join(stringx.Remove(botCommandFieldNames, "`id`", "`create_time`", "`update_time`", "`create_at`", "`update_at`"), ",")
	botCommandRowsWithPlaceHolder = strings.Join(stringx.Remove(botCommandFieldNames, "`id`", "`create_time`", "`update_time`", "`create_at`", "`update_at`"), "=?,") + "=?"

	cacheBotCommandIdPrefix                                = "cache:botCommand:id:"
	cacheBotCommandBotCommandLanguageCodeDeletedTimePrefix = "cache:botCommand:botCommand:languageCode:deletedTime:"
)

type (
	botCommandModel interface {
		Insert(ctx context.Context, data *BotCommand, s ...sqlx.Session) (sql.Result, error)
		FindOne(ctx context.Context, id int64, s ...sqlx.Session) (*BotCommand, error)
		FindOneByBotCommandLanguageCodeDeletedTime(ctx context.Context, botCommand string, languageCode string, deletedTime int64) (*BotCommand, error)
		FindOneByBuilder(ctx context.Context, sbs ...squirrel.SelectBuilder) (*BotCommand, error)
		Update(ctx context.Context, data *BotCommand, s ...sqlx.Session) error
		ListPage(ctx context.Context, page, size int64, sbs ...squirrel.SelectBuilder) ([]*BotCommand, int64, error)
		List(ctx context.Context, sbs ...squirrel.SelectBuilder) ([]*BotCommand, error)
		Count(ctx context.Context, cbs ...squirrel.SelectBuilder) (int64, error)
		Delete(ctx context.Context, id int64, s ...sqlx.Session) error
	}

	defaultBotCommandModel struct {
		sqlc.CachedConn
		table string
	}

	BotCommand struct {
		Id           int64  `db:"id"`
		LanguageCode string `db:"language_code"` // 语言代码
		BotCommand   string `db:"bot_command"`   // 指令
		Description  string `db:"description"`   // 命令描述
		Text         string `db:"text"`          // 消息
		Photo        string `db:"photo"`         // 图片字节
		SendType     int64  `db:"send_type"`     // 类型 1文本 2图片
		ButtonArray  string `db:"button_array"`  // 按钮二维数组
		Status       int64  `db:"status"`        // 状态 0关闭 1开启
		CreatedTime  int64  `db:"created_time"`
		UpdatedTime  int64  `db:"updated_time"`
		DeletedTime  int64  `db:"deleted_time"`
	}
)

func newBotCommandModel(conn sqlx.SqlConn, c cache.CacheConf) *defaultBotCommandModel {
	return &defaultBotCommandModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`bot_command`",
	}
}

func (m *defaultBotCommandModel) Delete(ctx context.Context, id int64, s ...sqlx.Session) error {
	data, err := m.FindOne(ctx, id, s...)
	if err != nil {
		return err
	}

	botCommandBotCommandLanguageCodeDeletedTimeKey := fmt.Sprintf("%s%v:%v:%v", cacheBotCommandBotCommandLanguageCodeDeletedTimePrefix, data.BotCommand, data.LanguageCode, data.DeletedTime)
	botCommandIdKey := fmt.Sprintf("%s%v", cacheBotCommandIdPrefix, id)

	var query string
	if BotCommandSoftDelete {
		query = fmt.Sprintf("update %s set `deleted_time` = %v where `id` = ?", m.table, time.Now().Unix())
	} else {
		query = fmt.Sprintf("delete from %s where `id` = ?", m.table)
	}

	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		if s != nil {
			return s[0].ExecCtx(ctx, query, id)
		}
		return conn.ExecCtx(ctx, query, id)
	}, botCommandBotCommandLanguageCodeDeletedTimeKey, botCommandIdKey)
	return err
}

func (m *defaultBotCommandModel) FindOne(ctx context.Context, id int64, s ...sqlx.Session) (*BotCommand, error) {
	botCommandIdKey := fmt.Sprintf("%s%v", cacheBotCommandIdPrefix, id)
	var resp BotCommand
	err := m.QueryRowCtx(ctx, &resp, botCommandIdKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", botCommandRows, m.table)
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

func (m *defaultBotCommandModel) FindOneByBotCommandLanguageCodeDeletedTime(ctx context.Context, botCommand string, languageCode string, deletedTime int64) (*BotCommand, error) {
	botCommandBotCommandLanguageCodeDeletedTimeKey := fmt.Sprintf("%s%v:%v:%v", cacheBotCommandBotCommandLanguageCodeDeletedTimePrefix, botCommand, languageCode, deletedTime)
	var resp BotCommand
	err := m.QueryRowIndexCtx(ctx, &resp, botCommandBotCommandLanguageCodeDeletedTimeKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `bot_command` = ? and `language_code` = ? and `deleted_time` = ? limit 1", botCommandRows, m.table)
		if err := conn.QueryRowCtx(ctx, &resp, query, botCommand, languageCode, deletedTime); err != nil {
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

func (m *defaultBotCommandModel) Insert(ctx context.Context, data *BotCommand, s ...sqlx.Session) (sql.Result, error) {
	botCommandBotCommandLanguageCodeDeletedTimeKey := fmt.Sprintf("%s%v:%v:%v", cacheBotCommandBotCommandLanguageCodeDeletedTimePrefix, data.BotCommand, data.LanguageCode, data.DeletedTime)
	botCommandIdKey := fmt.Sprintf("%s%v", cacheBotCommandIdPrefix, data.Id)
	return m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, botCommandRowsExpectAutoSet)
		if s != nil {
			return s[0].ExecCtx(ctx, query, data.LanguageCode, data.BotCommand, data.Description, data.Text, data.Photo, data.SendType, data.ButtonArray, data.Status, data.CreatedTime, data.UpdatedTime, data.DeletedTime)
		}
		return conn.ExecCtx(ctx, query, data.LanguageCode, data.BotCommand, data.Description, data.Text, data.Photo, data.SendType, data.ButtonArray, data.Status, data.CreatedTime, data.UpdatedTime, data.DeletedTime)
	}, botCommandBotCommandLanguageCodeDeletedTimeKey, botCommandIdKey)
}

func (m *defaultBotCommandModel) Update(ctx context.Context, newData *BotCommand, s ...sqlx.Session) error {
	data, err := m.FindOne(ctx, newData.Id, s...)
	if err != nil {
		return err
	}

	botCommandBotCommandLanguageCodeDeletedTimeKey := fmt.Sprintf("%s%v:%v:%v", cacheBotCommandBotCommandLanguageCodeDeletedTimePrefix, data.BotCommand, data.LanguageCode, data.DeletedTime)
	botCommandIdKey := fmt.Sprintf("%s%v", cacheBotCommandIdPrefix, data.Id)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, botCommandRowsWithPlaceHolder)
		if s != nil {
			return s[0].ExecCtx(ctx, query, newData.LanguageCode, newData.BotCommand, newData.Description, newData.Text, newData.Photo, newData.SendType, newData.ButtonArray, newData.Status, newData.CreatedTime, newData.UpdatedTime, newData.DeletedTime, newData.Id)
		}
		return conn.ExecCtx(ctx, query, newData.LanguageCode, newData.BotCommand, newData.Description, newData.Text, newData.Photo, newData.SendType, newData.ButtonArray, newData.Status, newData.CreatedTime, newData.UpdatedTime, newData.DeletedTime, newData.Id)
	}, botCommandBotCommandLanguageCodeDeletedTimeKey, botCommandIdKey)
	return err
}

func (m *defaultBotCommandModel) FindOneByBuilder(ctx context.Context, sbs ...squirrel.SelectBuilder) (*BotCommand, error) {
	var sb squirrel.SelectBuilder
	if sbs == nil {
		sb = squirrel.Select()
	} else {
		sb = sbs[0]
	}

	// query rows
	query, args, err := sb.From(m.table).Columns(botCommandRows).Limit(1).ToSql()
	if err != nil {
		return nil, err
	}
	var ret BotCommand
	if err := m.QueryRowNoCacheCtx(ctx, &ret, query, args...); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (m *defaultBotCommandModel) ListPage(ctx context.Context, page, size int64, sbs ...squirrel.SelectBuilder) ([]*BotCommand, int64, error) {
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
	query, args, err = sb.From(m.table).Columns(botCommandRows).Offset(uint64((page - 1) * size)).Limit(uint64(size)).ToSql()
	if err != nil {
		return nil, 0, err
	}
	var list []*BotCommand
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query, args...); err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (m *defaultBotCommandModel) List(ctx context.Context, sbs ...squirrel.SelectBuilder) ([]*BotCommand, error) {
	var sb squirrel.SelectBuilder
	if sbs == nil {
		sb = squirrel.Select()
	} else {
		sb = sbs[0]
	}

	// query rows
	query, args, err := sb.From(m.table).Columns(botCommandRows).ToSql()
	if err != nil {
		return nil, err
	}
	var list []*BotCommand
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query, args...); err != nil {
		return nil, err
	}

	return list, nil
}

func (m *defaultBotCommandModel) Count(ctx context.Context, cbs ...squirrel.SelectBuilder) (int64, error) {
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

func (m *defaultBotCommandModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheBotCommandIdPrefix, primary)
}

func (m *defaultBotCommandModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", botCommandRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultBotCommandModel) tableName() string {
	return m.table
}

var BotCommandSoftDelete bool

func init() {
	tp := reflect.TypeOf(BotCommand{})
	for i := 0; i < tp.NumField(); i++ {
		if tp.Field(i).Tag.Get("db") == "deleted_time" {
			BotCommandSoftDelete = true
			return
		}
	}
}
