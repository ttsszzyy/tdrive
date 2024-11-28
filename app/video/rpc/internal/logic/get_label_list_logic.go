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

type GetLabelListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetLabelListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLabelListLogic {
	return &GetLabelListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetLabelListLogic) GetLabelList(in *videoPb.GetLabelListReq) (*videoPb.GetLabelListResp, error) {
	filer := bson.M{}
	if in.Vid > 0 {
		filer["vid"] = in.Vid
	}
	if in.Id != "" {
		oid, err := primitive.ObjectIDFromHex(in.Id)
		if err != nil {
			return nil, model.ErrInvalidObjectId
		}
		filer["_id"] = oid
	}

	var (
		labels []*model.Videos_label
		total  int64
		err    error
	)
	if in.Page == 0 && in.Size == 0 {
		labels, err = l.svcCtx.VideosLabelModel.List(l.ctx, filer)
		total = int64(len(labels))
	} else {
		labels, total, err = l.svcCtx.VideosLabelModel.ListPage(l.ctx, in.Page, in.Size, filer)
	}
	if err != nil {
		return nil, err
	}

	list := make([]*videoPb.VideosLabel, 0)
	for _, label := range labels {
		list = append(list, &videoPb.VideosLabel{
			Id:          label.ID.Hex(),
			Title:       label.Title,
			Uid:         label.Uid,
			Vid:         label.Vid,
			CreatedTime: label.CreateAt.Unix(),
		})
	}
	return &videoPb.GetLabelListResp{
		List:  list,
		Total: total,
	}, nil
}
