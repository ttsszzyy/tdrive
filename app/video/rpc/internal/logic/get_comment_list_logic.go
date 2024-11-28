package logic

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"

	"T-driver/app/video/rpc/internal/svc"
	"T-driver/app/video/rpc/videoPb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCommentListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentListLogic {
	return &GetCommentListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCommentListLogic) GetCommentList(in *videoPb.GetCommentListReq) (*videoPb.GetCommentListResp, error) {
	filer := bson.M{}
	if in.Vid > 0 {
		filer["vid"] = in.Vid
	}
	if in.Uid > 0 {
		filer["uid"] = in.Uid
	}
	page, total, err := l.svcCtx.VideosCommentModel.ListPage(l.ctx, in.Page, in.Size, filer)
	if err != nil {
		return nil, err
	}
	resp := &videoPb.GetCommentListResp{
		List:  make([]*videoPb.VideosComment, 0, len(page)),
		Total: total,
	}
	for _, v := range page {
		count, err := l.svcCtx.VideosReplyModel.Count(l.ctx, bson.M{"comment_id": v.ID.Hex()})
		if err != nil {
			return nil, err
		}
		resp.List = append(resp.List, &videoPb.VideosComment{
			Id:          v.ID.Hex(),
			Uid:         v.Uid,
			Vid:         v.Vid,
			Content:     v.Content,
			CreatedTime: v.CreateAt.Unix(),
			ReplyCount:  count,
		})
	}

	return &videoPb.GetCommentListResp{
		List:  resp.List,
		Total: resp.Total,
	}, nil
}
