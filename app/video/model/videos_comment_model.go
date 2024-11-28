package model

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/monc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ Videos_commentModel = (*customVideos_commentModel)(nil)

type (
	// Videos_commentModel is an interface to be customized, add more methods here,
	// and implement the added methods in customVideos_commentModel.
	Videos_commentModel interface {
		videos_commentModel
		ListPage(ctx context.Context, page, size int64, filter any) ([]*Videos_comment, int64, error)
	}

	customVideos_commentModel struct {
		*defaultVideos_commentModel
	}
)

// NewVideos_commentModel returns a model for the mongo.
func NewVideos_commentModel(url, db, collection string, c cache.CacheConf) Videos_commentModel {
	conn := monc.MustNewModel(url, db, collection, c)
	return &customVideos_commentModel{
		defaultVideos_commentModel: newDefaultVideos_commentModel(conn),
	}
}

func (m *defaultVideos_commentModel) ListPage(ctx context.Context, page, size int64, filter any) ([]*Videos_comment, int64, error) {
	total, err := m.conn.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	var data []*Videos_comment
	if err := m.conn.Find(ctx, &data, filter, options.Find().SetSort(bson.D{{"createAt", -1}}).SetSkip((page-1)*size).SetLimit(size)); err != nil {
		return nil, 0, err
	}
	return data, total, nil
}
