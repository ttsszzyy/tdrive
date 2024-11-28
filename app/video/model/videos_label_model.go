package model

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/monc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ Videos_labelModel = (*customVideos_labelModel)(nil)

type (
	// Videos_labelModel is an interface to be customized, add more methods here,
	// and implement the added methods in customVideos_labelModel.
	Videos_labelModel interface {
		videos_labelModel
		ListPage(ctx context.Context, page, size int64, filter any) ([]*Videos_label, int64, error)
		List(ctx context.Context, filter any) ([]*Videos_label, error)
	}

	customVideos_labelModel struct {
		*defaultVideos_labelModel
	}
)

// NewVideos_labelModel returns a model for the mongo.
func NewVideos_labelModel(url, db, collection string, c cache.CacheConf) Videos_labelModel {
	conn := monc.MustNewModel(url, db, collection, c)
	return &customVideos_labelModel{
		defaultVideos_labelModel: newDefaultVideos_labelModel(conn),
	}
}

func (m *defaultVideos_labelModel) ListPage(ctx context.Context, page, size int64, filter any) ([]*Videos_label, int64, error) {
	total, err := m.conn.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	var data []*Videos_label
	if err := m.conn.Find(ctx, &data, filter, options.Find().SetSort(bson.D{{"createAt", -1}}).SetSkip((page-1)*size).SetLimit(size)); err != nil {
		return nil, 0, err
	}
	return data, total, nil
}
func (m *defaultVideos_labelModel) List(ctx context.Context, filter any) ([]*Videos_label, error) {
	var data []*Videos_label
	if err := m.conn.Find(ctx, &data, filter); err != nil {
		return nil, err
	}
	return data, nil
}
