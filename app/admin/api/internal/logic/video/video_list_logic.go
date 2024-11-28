package video

import (
	"T-driver/app/video/rpc/videoPb"
	"context"

	"T-driver/app/admin/api/internal/svc"
	"T-driver/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type VideoListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVideoListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VideoListLogic {
	return &VideoListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VideoListLogic) VideoList(req *types.VideoListReq) (resp *types.VideoListRes, err error) {
	list, err := l.svcCtx.VideoRpc.GetVideoList(l.ctx, &videoPb.GetVideoListReq{
		Title: req.Titel,
		Page:  req.Page,
		Size:  req.Size,
		Desc:  req.Desc,
	})
	if err != nil {
		return nil, err
	}
	List := make([]*types.Video, 0)
	for _, v := range list.VideoList {
		List = append(List, &types.Video{
			Id:          v.Id,
			Title:       v.Title,
			Desc:        v.Desc,
			FilePath:    v.FilePath,
			CreatedTime: v.CreatedTime,
			Status:      v.Status,
		})
	}
	return &types.VideoListRes{
		List:  List,
		Total: list.Total,
	}, nil
}
