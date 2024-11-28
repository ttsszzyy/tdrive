package model

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/monc"
)

var _ Videos_label_userModel = (*customVideos_label_userModel)(nil)

type (
	// Videos_label_userModel is an interface to be customized, add more methods here,
	// and implement the added methods in customVideos_label_userModel.
	Videos_label_userModel interface {
		videos_label_userModel
		List(ctx context.Context, filter any) ([]*Videos_label_user, error)
		Count(ctx context.Context, filter any) (int64, error)
	}

	customVideos_label_userModel struct {
		*defaultVideos_label_userModel
	}
)

// NewVideos_label_userModel returns a model for the mongo.
func NewVideos_label_userModel(url, db, collection string, c cache.CacheConf) Videos_label_userModel {
	conn := monc.MustNewModel(url, db, collection, c)
	return &customVideos_label_userModel{
		defaultVideos_label_userModel: newDefaultVideos_label_userModel(conn),
	}
}

func (m *defaultVideos_label_userModel) List(ctx context.Context, filter any) ([]*Videos_label_user, error) {
	var data []*Videos_label_user
	if err := m.conn.Find(ctx, &data, filter); err != nil {
		return nil, err
	}
	return data, nil
}

func (m *defaultVideos_label_userModel) Count(ctx context.Context, filter any) (int64, error) {
	total, err := m.conn.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return total, nil
}
