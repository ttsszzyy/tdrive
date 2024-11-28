package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/monc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var _ UserTitanTokenModel = (*customUserTitanTokenModel)(nil)

type (
	// UserTitanTokenModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserTitanTokenModel.
	UserTitanTokenModel interface {
		userTitanTokenModel
		FindOneByUid(ctx context.Context, uid int64) (*UserTitanToken, error)
		UpdateNoCache(ctx context.Context, data *UserTitanToken) (*mongo.UpdateResult, error)
	}

	customUserTitanTokenModel struct {
		*defaultUserTitanTokenModel
	}
)

// NewUserTitanTokenModel returns a model for the mongo.
func NewUserTitanTokenModel(url, db, collection string, c cache.CacheConf) UserTitanTokenModel {
	conn := monc.MustNewModel(url, db, collection, c)
	return &customUserTitanTokenModel{
		defaultUserTitanTokenModel: newDefaultUserTitanTokenModel(conn),
	}
}

func (m *defaultUserTitanTokenModel) FindOneByUid(ctx context.Context, uid int64) (*UserTitanToken, error) {

	var data UserTitanToken
	key := fmt.Sprintf("%s%d", prefixUserTitanTokenCacheKey, uid)
	err := m.conn.FindOne(ctx, key, &data, bson.M{"uid": uid})
	switch err {
	case nil:
		return &data, nil
	case monc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserTitanTokenModel) UpdateNoCache(ctx context.Context, data *UserTitanToken) (*mongo.UpdateResult, error) {
	data.UpdateAt = time.Now()
	key := prefixUserTitanTokenCacheKey + data.ID.Hex()
	ukey := fmt.Sprintf("%s%d", prefixUserTitanTokenCacheKey, data.Uid)
	res, err := m.conn.UpdateOne(ctx, key, bson.M{"_id": data.ID}, bson.M{"$set": data})
	if err != nil {
		return nil, err
	}
	//删除uid的缓存
	err = m.conn.DelCache(ctx, ukey)
	return res, err
}
