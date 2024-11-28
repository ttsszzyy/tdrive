package model

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/monc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ Videos_replyModel = (*customVideos_replyModel)(nil)

type (
	// Videos_replyModel is an interface to be customized, add more methods here,
	// and implement the added methods in customVideos_replyModel.
	Videos_replyModel interface {
		videos_replyModel
		ListPage(ctx context.Context, page, size int64, filter any) ([]*Videos_reply, int64, error)
		Count(ctx context.Context, filter any) (int64, error)
	}

	customVideos_replyModel struct {
		*defaultVideos_replyModel
	}
)

// NewVideos_replyModel returns a model for the mongo.
func NewVideos_replyModel(url, db, collection string, c cache.CacheConf) Videos_replyModel {
	conn := monc.MustNewModel(url, db, collection, c)
	return &customVideos_replyModel{
		defaultVideos_replyModel: newDefaultVideos_replyModel(conn),
	}
}

func (m *defaultVideos_replyModel) ListPage(ctx context.Context, page, size int64, filter any) ([]*Videos_reply, int64, error) {
	total, err := m.conn.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	var data []*Videos_reply
	if err := m.conn.Find(ctx, &data, filter, options.Find().SetSort(bson.D{{"createAt", -1}}).SetSkip((page-1)*size).SetLimit(size)); err != nil {
		return nil, 0, err
	}
	return data, total, nil
}

func (m *defaultVideos_replyModel) Count(ctx context.Context, filter any) (int64, error) {
	total, err := m.conn.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return total, nil
}
