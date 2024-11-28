// Code generated by goctl. DO NOT EDIT.
package model

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/stores/monc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var prefixAssetFileCacheKey = "cache:assetFile:"

type assetFileModel interface {
	Insert(ctx context.Context, data *AssetFile) error
	FindOne(ctx context.Context, id string) (*AssetFile, error)
	Update(ctx context.Context, data *AssetFile) (*mongo.UpdateResult, error)
	Delete(ctx context.Context, id string) (int64, error)
}

type defaultAssetFileModel struct {
	conn *monc.Model
}

func newDefaultAssetFileModel(conn *monc.Model) *defaultAssetFileModel {
	return &defaultAssetFileModel{conn: conn}
}

func (m *defaultAssetFileModel) Insert(ctx context.Context, data *AssetFile) error {
	if data.ID.IsZero() {
		data.ID = primitive.NewObjectID()
		data.CreateAt = time.Now()
		data.UpdateAt = time.Now()
	}

	key := prefixAssetFileCacheKey + data.ID.Hex()
	_, err := m.conn.InsertOne(ctx, key, data)
	return err
}

func (m *defaultAssetFileModel) FindOne(ctx context.Context, id string) (*AssetFile, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrInvalidObjectId
	}

	var data AssetFile
	key := prefixAssetFileCacheKey + id
	err = m.conn.FindOne(ctx, key, &data, bson.M{"_id": oid})
	switch err {
	case nil:
		return &data, nil
	case monc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultAssetFileModel) Update(ctx context.Context, data *AssetFile) (*mongo.UpdateResult, error) {
	data.UpdateAt = time.Now()
	key := prefixAssetFileCacheKey + data.ID.Hex()
	res, err := m.conn.UpdateOne(ctx, key, bson.M{"_id": data.ID}, bson.M{"$set": data})
	return res, err
}

func (m *defaultAssetFileModel) Delete(ctx context.Context, id string) (int64, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, ErrInvalidObjectId
	}
	key := prefixAssetFileCacheKey + id
	res, err := m.conn.DeleteOne(ctx, key, bson.M{"_id": oid})
	return res, err
}
