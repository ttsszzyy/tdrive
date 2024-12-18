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

var prefixVideos_label_userCacheKey = "cache:videos_label_user:"

type videos_label_userModel interface {
	Insert(ctx context.Context, data *Videos_label_user) error
	FindOne(ctx context.Context, id string) (*Videos_label_user, error)
	Update(ctx context.Context, data *Videos_label_user) (*mongo.UpdateResult, error)
	Delete(ctx context.Context, id string) (int64, error)
}

type defaultVideos_label_userModel struct {
	conn *monc.Model
}

func newDefaultVideos_label_userModel(conn *monc.Model) *defaultVideos_label_userModel {
	return &defaultVideos_label_userModel{conn: conn}
}

func (m *defaultVideos_label_userModel) Insert(ctx context.Context, data *Videos_label_user) error {
	if data.ID.IsZero() {
		data.ID = primitive.NewObjectID()
		data.CreateAt = time.Now()
		data.UpdateAt = time.Now()
	}

	key := prefixVideos_label_userCacheKey + data.ID.Hex()
	_, err := m.conn.InsertOne(ctx, key, data)
	return err
}

func (m *defaultVideos_label_userModel) FindOne(ctx context.Context, id string) (*Videos_label_user, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrInvalidObjectId
	}

	var data Videos_label_user
	key := prefixVideos_label_userCacheKey + id
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

func (m *defaultVideos_label_userModel) Update(ctx context.Context, data *Videos_label_user) (*mongo.UpdateResult, error) {
	data.UpdateAt = time.Now()
	key := prefixVideos_label_userCacheKey + data.ID.Hex()
	res, err := m.conn.UpdateOne(ctx, key, bson.M{"_id": data.ID}, bson.M{"$set": data})
	return res, err
}

func (m *defaultVideos_label_userModel) Delete(ctx context.Context, id string) (int64, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, ErrInvalidObjectId
	}
	key := prefixVideos_label_userCacheKey + id
	res, err := m.conn.DeleteOne(ctx, key, bson.M{"_id": oid})
	return res, err
}
