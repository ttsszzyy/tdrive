package logic

import (
	"context"

	"T-driver/app/video/rpc/internal/svc"
	"T-driver/app/video/rpc/videoPb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateVideoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateVideoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateVideoLogic {
	return &UpdateVideoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateVideoLogic) UpdateVideo(in *videoPb.UpdateVideoReq) (*videoPb.VideoResponse, error) {
	one, err := l.svcCtx.VideosModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	if in.Title != "" {
		one.Title = in.Title
	}
	if in.Desc != "" {
		one.Desc = in.Desc
	}
	if in.FilePath != "" {
		one.FilePath = in.FilePath
	}
	if in.Cid != "" {
		one.Cid = in.Cid
	}
	if in.Url != "" {
		one.Link = in.Url
	}
	if in.Status > 0 {
		one.Status = in.Status
	}
	err = l.svcCtx.VideosModel.Update(l.ctx, one)
	if err != nil {
		return nil, err
	}
	return &videoPb.VideoResponse{}, nil
}
