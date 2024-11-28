package video

import (
	"T-driver/app/video/rpc/video"
	"context"

	"T-driver/app/admin/api/internal/svc"
	"T-driver/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type VideoDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVideoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VideoDeleteLogic {
	return &VideoDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VideoDeleteLogic) VideoDelete(req *types.VideoDeleteReq) (resp *types.Response, err error) {
	_, err = l.svcCtx.VideoRpc.DeleteVideo(l.ctx, &video.DeleteVideoReq{Id: req.Id})
	if err != nil {
		return nil, err
	}
	return &types.Response{}, nil
}
