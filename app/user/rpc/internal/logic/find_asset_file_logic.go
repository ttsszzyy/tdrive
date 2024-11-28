package logic

import (
	"T-driver/app/user/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindAssetFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindAssetFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindAssetFileLogic {
	return &FindAssetFileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查看资源
func (l *FindAssetFileLogic) FindAssetFile(in *pb.FindAssetFileReq) (*pb.FindAssetFileResp, error) {
	filer := bson.M{}
	if len(in.Status) > 0 {
		filer["status"] = bson.M{"$in": in.Status}
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
	var (
		list  []*model.AssetFile
		total int64
		err   error
	)
	if in.Page == 0 && in.Size == 0 {
		list, err = l.svcCtx.AssetFileModel.List(l.ctx, filer)
		total = int64(len(list))
	} else {
		list, total, err = l.svcCtx.AssetFileModel.ListPage(l.ctx, in.Page, in.Size, filer)
	}
	if err != nil {
		return nil, err
	}
	resp := &pb.FindAssetFileResp{
		List:  make([]*pb.AssetFile, 0, len(list)),
		Total: total,
	}
	for _, v := range list {
		resp.List = append(resp.List, &pb.AssetFile{
			Id:          v.ID.Hex(),
			Uid:         v.Uid,
			Path:        v.Path,
			Link:        v.Link,
			Cid:         v.Cid,
			AssetId:     v.AssetId,
			AssetName:   v.AssetName,
			AssetSize:   v.AssetSize,
			AssetType:   v.AssetType,
			Source:      v.Source,
			Status:      v.Status,
			Tag:         v.IsTag,
			Pid:         v.Pid,
			CreatedTime: v.CreateAt.Unix(),
			UpdatedTime: v.UpdateAt.Unix(),
		})
	}

	return resp, nil
}
