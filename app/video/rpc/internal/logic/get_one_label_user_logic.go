package logic

import (
	"T-driver/app/video/model"
	"context"

	"T-driver/app/video/rpc/internal/svc"
	"T-driver/app/video/rpc/videoPb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOneLabelUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOneLabelUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOneLabelUserLogic {
	return &GetOneLabelUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetOneLabelUserLogic) GetOneLabelUser(in *videoPb.GetOneLabelUserReq) (*videoPb.VideosLabelUser, error) {
	one, err := l.svcCtx.VideosLabelUserModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if model.ErrNotFound == err {
			return &videoPb.VideosLabelUser{}, err
		}
		return nil, err
	}
	return &videoPb.VideosLabelUser{
		Id:      one.ID.Hex(),
		Lid:     one.Lid,
		Uid:     one.Uid,
		Likes:   one.Likes,
		NoLikes: one.NoLikes,
	}, nil
}
