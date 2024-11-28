package logic

import (
	"T-driver/app/video/model"
	"T-driver/app/video/rpc/internal/svc"
	"T-driver/app/video/rpc/videoPb"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCommentLogic {
	return &CreateCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 评论
func (l *CreateCommentLogic) CreateComment(in *videoPb.CreateCommentReq) (*videoPb.VideoResponse, error) {
	err := l.svcCtx.VideosCommentModel.Insert(l.ctx, &model.Videos_comment{
		Uid:         in.Uid,
		Vid:         in.Vid,
		ComposeType: "",
		Content:     in.Content,
	})
	if err != nil {
		return nil, err
	}

	return &videoPb.VideoResponse{}, nil
}
