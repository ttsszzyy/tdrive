package logic

import (
	"T-driver/app/video/model"
	"T-driver/app/video/rpc/internal/svc"
	"T-driver/app/video/rpc/videoPb"
	"context"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLabelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateLabelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLabelLogic {
	return &CreateLabelLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 标签
func (l *CreateLabelLogic) CreateLabel(in *videoPb.CreateLabelReq) (*videoPb.CreateLabelResp, error) {
	labels, err := l.svcCtx.VideosLabelModel.List(l.ctx, bson.M{"vid": in.Vid, "title": in.Title})
	if err != nil {
		return nil, err
	}
	if len(labels) > 0 {
		_, err = NewSaveLabelUserLogic(l.ctx, l.svcCtx).SaveLabelUser(&videoPb.SaveLabelUserReq{
			Lid:     labels[0].ID.Hex(),
			Uid:     in.Uid,
			Likes:   true,
			NoLikes: false,
		})
		if err != nil {
			return nil, err
		}
		return &videoPb.CreateLabelResp{Repeat: true}, nil
	}
	err = l.svcCtx.VideosLabelModel.Insert(l.ctx, &model.Videos_label{
		Uid:   in.Uid,
		Vid:   in.Vid,
		Title: in.Title,
	})
	if err != nil {
		return nil, err
	}

	return &videoPb.CreateLabelResp{}, nil
}
