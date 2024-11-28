package logic

import (
	"T-driver/app/video/model"
	"context"

	"T-driver/app/video/rpc/internal/svc"
	"T-driver/app/video/rpc/videoPb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOneLabelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOneLabelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOneLabelLogic {
	return &GetOneLabelLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetOneLabelLogic) GetOneLabel(in *videoPb.GetOneLabelReq) (*videoPb.VideosLabel, error) {
	one, err := l.svcCtx.VideosLabelModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if model.ErrNotFound == err {
			return &videoPb.VideosLabel{}, err
		}
		return nil, err
	}
	return &videoPb.VideosLabel{
		Id:          one.ID.Hex(),
		Title:       one.Title,
		Uid:         one.Uid,
		Vid:         one.Vid,
		CreatedTime: one.CreateAt.Unix(),
	}, nil
}
