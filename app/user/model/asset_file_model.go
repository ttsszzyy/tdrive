package model

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/monc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ AssetFileModel = (*customAssetFileModel)(nil)

type (
	// AssetFileModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAssetFileModel.
	AssetFileModel interface {
		assetFileModel
		ListPage(ctx context.Context, page, size int64, filter any) ([]*AssetFile, int64, error)
		List(ctx context.Context, filter any) ([]*AssetFile, error)
		Deletes(ctx context.Context, ids []string) (int64, error)
		FindOneAssetId(ctx context.Context, assetId int64) (*AssetFile, error)
		UpdateByAssetIds(ctx context.Context, ids []string, status int64) (int64, error)
	}

	customAssetFileModel struct {
		*defaultAssetFileModel
	}
)

// NewAssetFileModel returns a model for the mongo.
func NewAssetFileModel(url, db, collection string, c cache.CacheConf) AssetFileModel {
	conn := monc.MustNewModel(url, db, collection, c)
	return &customAssetFileModel{
		defaultAssetFileModel: newDefaultAssetFileModel(conn),
	}
}

func (m *defaultAssetFileModel) ListPage(ctx context.Context, page, size int64, filter any) ([]*AssetFile, int64, error) {
	total, err := m.conn.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	var data []*AssetFile
	if err := m.conn.Find(ctx, &data, filter, options.Find().SetSort(bson.D{{"createAt", -1}}).SetSkip((page-1)*size).SetLimit(size)); err != nil {
		return nil, 0, err
	}
	return data, total, nil
}
func (m *defaultAssetFileModel) List(ctx context.Context, filter any) ([]*AssetFile, error) {
	var data []*AssetFile
	if err := m.conn.Find(ctx, &data, filter); err != nil {
		return nil, err
	}
	return data, nil
}

func (m *defaultAssetFileModel) Deletes(ctx context.Context, ids []string) (int64, error) {
	oids := make([]primitive.ObjectID, len(ids))
	keys := make([]string, len(ids))
	for _, id := range ids {
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return 0, ErrInvalidObjectId
		}
		oids = append(oids, oid)
		keys = append(keys, prefixAssetFileCacheKey+id)
	}

	res, err := m.conn.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": oids}})
	err = m.conn.DelCache(ctx, keys...)
	return res, err
}

func (m *defaultAssetFileModel) UpdateByAssetIds(ctx context.Context, ids []string, status int64) (int64, error) {
	var data []*AssetFile
	err := m.conn.Find(ctx, data, bson.M{"assetId": bson.M{"$in": ids}})
	if err != nil {
		return 0, err
	}
	keys := make([]string, 0, len(ids))
	oids := make([]primitive.ObjectID, 0, len(ids))
	for _, v := range data {
		oids = append(oids, v.ID)
		keys = append(keys, prefixAssetFileCacheKey+v.ID.Hex())
	}

	res, err := m.conn.UpdateMany(ctx, keys, bson.M{"_id": bson.M{"$in": oids}}, bson.M{"$set": bson.M{"status": status}})
	return res.ModifiedCount, err
}

func (m *defaultAssetFileModel) FindOneAssetId(ctx context.Context, assetId int64) (*AssetFile, error) {
	var data AssetFile
	err := m.conn.FindOneNoCache(ctx, &data, bson.M{"assetId": assetId})
	switch err {
	case nil:
		return &data, nil
	case monc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
