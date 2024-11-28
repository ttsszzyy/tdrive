package logic

import (
	"T-driver/app/video/model"
	"T-driver/app/video/rpc/internal/svc"
	"T-driver/app/video/rpc/videoPb"
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateVideoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateVideoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateVideoLogic {
	return &CreateVideoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateVideoLogic) CreateVideo(in *videoPb.CreateVideoReq) (*videoPb.CreateVideoResp, error) {
	res, err := l.svcCtx.VideosModel.Insert(l.ctx, &model.Videos{
		Uid:         in.Uid,
		Title:       in.Title,
		Desc:        in.Desc,
		Filename:    in.Filename,
		Status:      in.Status,
		CreatedTime: time.Now().Unix(),
		UpdatedTime: time.Now().Unix(),
	})
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return &videoPb.CreateVideoResp{Id: id}, nil
}
