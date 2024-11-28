package logic

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"

	"T-driver/app/video/rpc/internal/svc"
	"T-driver/app/video/rpc/videoPb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetReplyListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetReplyListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetReplyListLogic {
	return &GetReplyListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetReplyListLogic) GetReplyList(in *videoPb.GetReplyListReq) (*videoPb.GetReplyListResp, error) {
	filer := bson.M{}
	if in.CommentId != "" {
		filer["comment_id"] = in.CommentId
	}
	if in.ToUserid > 0 {
		filer["to_userid"] = in.ToUserid
	}
	page, total, err := l.svcCtx.VideosReplyModel.ListPage(l.ctx, in.Page, in.Size, filer)
	if err != nil {
		return nil, err
	}
	resp := &videoPb.GetReplyListResp{
		List:  make([]*videoPb.VideosReply, 0, len(page)),
		Total: total,
	}
	for _, v := range page {
		resp.List = append(resp.List, &videoPb.VideosReply{
			Id:          v.ID.Hex(),
			CommentId:   v.CommentId,
			ToUserId:    v.ToUserid,
			FromUserId:  v.FromUserid,
			Content:     v.Content,
			CreatedTime: v.CreateAt.Unix(),
		})
	}
	return &videoPb.GetReplyListResp{
		List:  resp.List,
		Total: resp.Total,
	}, nil
}
