package logic

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"

	"T-driver/app/video/rpc/internal/svc"
	"T-driver/app/video/rpc/videoPb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CountLabelUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCountLabelUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CountLabelUserLogic {
	return &CountLabelUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CountLabelUserLogic) CountLabelUser(in *videoPb.CountLabelUserReq) (*videoPb.CountLabelUserResp, error) {
	list := make([]*videoPb.CountLabelUserItem, 0, len(in.Lid))
	for _, v := range in.Lid {
		likesNum, err := l.svcCtx.VideosLabelUserModel.Count(l.ctx, bson.M{"lid": v, "likes": true})
		if err != nil {
			return nil, err
		}
		noLikesNum, err := l.svcCtx.VideosLabelUserModel.Count(l.ctx, bson.M{"lid": v, "no_likes": true})
		if err != nil {
			return nil, err
		}
		list = append(list, &videoPb.CountLabelUserItem{
			Lid:        v,
			LikesNum:   likesNum,
			NoLikesNum: noLikesNum,
		})
	}

	return &videoPb.CountLabelUserResp{
		List: list,
	}, nil
}
