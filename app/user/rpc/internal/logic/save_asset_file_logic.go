package logic

import (
	"T-driver/app/user/model"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"

	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveAssetFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveAssetFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveAssetFileLogic {
	return &SaveAssetFileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 保存资源
func (l *SaveAssetFileLogic) SaveAssetFile(in *pb.SaveAssetFileReq) (*pb.SaveAssetFileResp, error) {
	ID := primitive.NewObjectID()
	one := &model.AssetFile{
		ID:        ID,
		Uid:       in.Uid,
		AssetName: in.AssetName,
		AssetSize: in.AssetSize,
		AssetType: in.AssetType,
		Pid:       in.Pid,
		Source:    in.Source,
		Status:    model.AssetStatusAfoot,
		Path:      in.Path,
		IsTag:     model.AssetStatusAfoot,
		UpdateAt:  time.Now(),
		CreateAt:  time.Now(),
	}
	err := l.svcCtx.AssetFileModel.Insert(l.ctx, one)
	if err != nil {
		return nil, err
	}

	return &pb.SaveAssetFileResp{Id: ID.Hex()}, nil
}
