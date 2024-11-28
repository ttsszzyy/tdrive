package model

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ActionRecordModel = (*customActionRecordModel)(nil)

type (
	// ActionRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customActionRecordModel.
	ActionRecordModel interface {
		actionRecordModel
		GetIPNumByAction(ctx context.Context, action, ip string) (int64, error)
	}

	customActionRecordModel struct {
		*defaultActionRecordModel
	}
)

// NewActionRecordModel returns a model for the database table.
func NewActionRecordModel(conn sqlx.SqlConn, c cache.CacheConf) ActionRecordModel {
	return &customActionRecordModel{
		defaultActionRecordModel: newActionRecordModel(conn, c),
	}
}

func (m *customActionRecordModel) GetIPNumByAction(ctx context.Context, action, ip string) (int64, error) {
	var num int64

	query, args, err := squirrel.Select("COUNT(id)").From(m.table).Where("action = ? AND ip = ?", action, ip).ToSql()
	if err != nil {
		return 0, fmt.Errorf("generate sql of get ip's count error:%w", err)
	}

	err = m.QueryRowNoCacheCtx(ctx, &num, query, args...)
	if err != nil {
		return 0, fmt.Errorf("get ip's count error:%w", err)
	}

	return num, nil
}
