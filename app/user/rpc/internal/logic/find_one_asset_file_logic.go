package logic

import (
	"T-driver/app/user/model"
	"context"

	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindOneAssetFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindOneAssetFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindOneAssetFileLogic {
	return &FindOneAssetFileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindOneAssetFileLogic) FindOneAssetFile(in *pb.FindOneAssetFileReq) (*pb.AssetFile, error) {
	v := &model.AssetFile{}
	var err error
	if in.Id != "" {
		v, err = l.svcCtx.AssetFileModel.FindOne(l.ctx, in.Id)
	} else {
		v, err = l.svcCtx.AssetFileModel.FindOneAssetId(l.ctx, in.AssetId)
	}
	if err != nil {
		return nil, err
	}

	return &pb.AssetFile{
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
	}, nil
}
