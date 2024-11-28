package logic

import (
	"T-driver/app/video/model"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"T-driver/app/video/rpc/internal/svc"
	"T-driver/app/video/rpc/videoPb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLabelUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetLabelUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLabelUserLogic {
	return &GetLabelUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetLabelUserLogic) GetLabelUser(in *videoPb.GetLabelUserReq) (*videoPb.GetLabelUserResp, error) {

	filer := bson.M{}
	if len(in.Lid) > 0 {
		filer["lid"] = bson.M{"$in": in.Lid}
	}
	if in.Uid > 0 {
		filer["uid"] = in.Uid
	}
	if in.Id != "" {
		oid, err := primitive.ObjectIDFromHex(in.Id)
		if err != nil {
			return nil, model.ErrInvalidObjectId
		}
		filer["_id"] = oid
	}
	list, err := l.svcCtx.VideosLabelUserModel.List(l.ctx, filer)
	if err != nil {
		return nil, err
	}
	objs := make([]*videoPb.VideosLabelUser, 0)
	for _, v := range list {
		objs = append(objs, &videoPb.VideosLabelUser{
			CreatedTime: v.CreateAt.Unix(),
			Id:          v.ID.Hex(),
			Lid:         v.Lid,
			Likes:       v.Likes,
			NoLikes:     v.NoLikes,
			Uid:         v.Uid,
		})
	}
	return &videoPb.GetLabelUserResp{
		List: objs,
	}, nil
}
