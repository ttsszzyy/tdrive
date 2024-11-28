package logic

import (
	"context"

	"T-driver/app/video/rpc/internal/svc"
	"T-driver/app/video/rpc/videoPb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteVideoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteVideoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteVideoLogic {
	return &DeleteVideoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteVideoLogic) DeleteVideo(in *videoPb.DeleteVideoReq) (*videoPb.VideoResponse, error) {
	err := l.svcCtx.VideosModel.Delete(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return &videoPb.VideoResponse{}, nil
}
