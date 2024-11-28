package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AssetsModel = (*customAssetsModel)(nil)

type (
	// AssetsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAssetsModel.
	AssetsModel interface {
		assetsModel
		Clear(ctx context.Context, ids []int64, s ...sqlx.Session) error
		Deletes(ctx context.Context, ids []int64, s ...sqlx.Session) error
		GetAllIDsByPID(ctx context.Context, pid []int64, uid int64) ([]int64, error)
		GetUserAssetIDs(ctx context.Context, uid int64, ids []int64) ([]int64, int64, error)
		GetLastAssetInfo(ctx context.Context, uid int64) (*Assets, error)
	}

	customAssetsModel struct {
		*defaultAssetsModel
	}
)

// NewAssetsModel returns a model for the database table.
func NewAssetsModel(conn sqlx.SqlConn, c cache.CacheConf) AssetsModel {
	return &customAssetsModel{
		defaultAssetsModel: newAssetsModel(conn, c),
	}
}

func (m *defaultAssetsModel) Clear(ctx context.Context, ids []int64, s ...sqlx.Session) error {
	assetsIdKeys := make([]string, 0, len(ids))
	for _, id := range ids {
		assetsIdKey := fmt.Sprintf("%s%v", cacheAssetsIdPrefix, id)
		assetsIdKeys = append(assetsIdKeys, assetsIdKey)
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
	}, assetsIdKeys...)
	return err
}

func (m *defaultAssetsModel) Deletes(ctx context.Context, ids []int64, s ...sqlx.Session) error {
	assetsIdKeys := make([]string, 0)
	for _, id := range ids {
		assetsIdKey := fmt.Sprintf("%s%v", cacheAssetsIdPrefix, id)
		assetsIdKeys = append(assetsIdKeys, assetsIdKey)
	}
	query, args, err := squirrel.Update(m.table).Set("deleted_time", time.Now().Unix()).Where(squirrel.Eq{"id": ids}).ToSql()
	if err != nil {
		return err
	}

	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		if len(s) > 0 {
			return s[0].ExecCtx(ctx, query, args...)
		}
		return conn.ExecCtx(ctx, query, args...)
	}, assetsIdKeys...)
	return err
}

// GetAllIDsByPID 根据文件父id获取所有文件id
func (m *defaultAssetsModel) GetAllIDsByPID(ctx context.Context, pid []int64, uid int64) ([]int64, error) {
	var ids, pids []int64

	ids = append(ids, pid...)
	pids = append(pids, pid...)
	for {
		subIds := make([]int64, 0)
		query, args, err := squirrel.Select("id").From(m.table).Where(squirrel.Eq{"pid": pids, "uid": uid}).ToSql()
		if err != nil {
			return nil, fmt.Errorf("generate sql of get asset ids by pid error:%w", err)
		}
		err = m.QueryRowsNoCacheCtx(ctx, &subIds, query, args...)
		if err != nil {
			return nil, fmt.Errorf("execute sql of get asset ids by pid error:%w", err)
		}
		pids = subIds
		logx.Errorf("subIds:%v %v", subIds, pids)
		ids = append(ids, subIds...)
		if len(subIds) == 0 {
			logx.Errorf("ids:%v", ids)
			return ids, nil
		}
	}
}

// GetUserAssetIDs 获取用户自己的ids
func (m *defaultAssetsModel) GetUserAssetIDs(ctx context.Context, uid int64, ids []int64) ([]int64, int64, error) {
	var (
		nids []int64
		size int64
	)

	sq := squirrel.Select("id").From(m.table).Where(squirrel.Eq{"id": ids, "deleted_time": 0})
	if uid > 0 {
		sq = sq.Where("uid = ?", uid)
	}
	query, args, err := sq.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("generate sql of get asset ids by pid error:%w", err)
	}
	err = m.QueryRowsNoCacheCtx(ctx, &nids, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("execute sql of get asset ids by pid error:%w", err)
	}

	query, args, err = squirrel.Select("SUM(asset_size)").From(m.table).Where(squirrel.Eq{"id": nids}).ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("generate sql of get asset's size error:%w", err)
	}
	err = m.QueryRowNoCacheCtx(ctx, &size, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("execute sql of get asset's size error:%w", err)
	}

	return nids, size, nil
}

// GetLastAssetInfo 获取最后一次用户文件信息
func (m *defaultAssetsModel) GetLastAssetInfo(ctx context.Context, uid int64) (*Assets, error) {
	var info Assets

	query, args, err := squirrel.Select("*").From(m.table).Where("uid = ?", uid).OrderBy("created_time DESC").Limit(1).ToSql()
	if err != nil {
		return nil, fmt.Errorf("generate sql of get asset's info error:%w", err)
	}

	err = m.QueryRowNoCacheCtx(ctx, &info, query, args...)
	if err != nil {
		return nil, fmt.Errorf("get asset's info error:%w", err)
	}
	switch err {
	case nil:
		return &info, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
