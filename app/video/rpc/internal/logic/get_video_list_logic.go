package logic

import (
	"T-driver/app/video/model"
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"

	"T-driver/app/video/rpc/internal/svc"
	"T-driver/app/video/rpc/videoPb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetVideoListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetVideoListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVideoListLogic {
	return &GetVideoListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetVideoListLogic) GetVideoList(in *videoPb.GetVideoListReq) (*videoPb.GetVideoListResp, error) {
	sb := squirrel.Select().Where("deleted_time = ?", 0)
	if in.Title != "" {
		sb = sb.Where("title like ?", fmt.Sprintf("%s%%", in.Title))
	}
	if in.Desc != "" {
		sb = sb.Where("title like ?", fmt.Sprintf("%s%%", in.Desc))
	}
	if in.Status > 0 {
		sb = sb.Where("status = ?", in.Status)
	}
	if len(in.Ids) > 0 {
		sb = sb.Where(squirrel.Eq{"id": in.Ids})
	}
	var (
		listPage []*model.Videos
		total    int64
		err      error
	)
	if in.Size == 0 && in.Page == 0 {
		listPage, err = l.svcCtx.VideosModel.List(l.ctx, sb)
		if err != nil {
			return nil, err
		}
	} else {
		listPage, total, err = l.svcCtx.VideosModel.ListPage(l.ctx, in.Page, in.Size, sb)
	}

	if err != nil {
		return nil, err
	}
	VideoList := make([]*videoPb.Video, 0)
	for _, v := range listPage {
		VideoList = append(VideoList, &videoPb.Video{
			Id:          v.Id,
			Title:       v.Title,
			Desc:        v.Desc,
			Uid:         v.Uid,
			FilePath:    v.FilePath,
			CreatedTime: v.CreatedTime,
			UpdatedTime: v.UpdatedTime,
			Url:         v.Link,
			Cid:         v.Cid,
			Status:      v.Status,
		})
	}

	return &videoPb.GetVideoListResp{
		Total:     total,
		VideoList: VideoList,
	}, err
}
