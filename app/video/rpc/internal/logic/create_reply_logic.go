package logic

import (
	"T-driver/app/video/model"
	"T-driver/app/video/rpc/internal/svc"
	"T-driver/app/video/rpc/videoPb"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateReplyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateReplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateReplyLogic {
	return &CreateReplyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 回复评论
func (l *CreateReplyLogic) CreateReply(in *videoPb.CreateReplyReq) (*videoPb.CreateReplyResp, error) {
	ID := primitive.NewObjectID()
	createAt := time.Now()
	err := l.svcCtx.VideosReplyModel.Insert(l.ctx, &model.Videos_reply{
		FromUserid: in.FromUserId,
		ToUserid:   in.ToUserId,
		ReplyType:  "",
		Content:    in.Content,
		CommentId:  in.CommentId,
		ID:         ID,
		CreateAt:   createAt,
		UpdateAt:   createAt,
	})
	if err != nil {
		return nil, err
	}

	return &videoPb.CreateReplyResp{
		Id:          ID.Hex(),
		Content:     in.Content,
		CreatedTime: createAt.Unix(),
		FromUserId:  in.FromUserId,
		ToUserId:    in.ToUserId,
		CommentId:   in.CommentId,
	}, nil
}
