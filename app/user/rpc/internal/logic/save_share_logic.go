package logic

import (
	"T-driver/app/user/model"
	"context"
	"time"

	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveShareLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveShareLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveShareLogic {
	return &SaveShareLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 保存分享任务
func (l *SaveShareLogic) SaveShare(in *pb.SaveShareReq) (*pb.Response, error) {
	_, err := l.svcCtx.ShareModel.Insert(l.ctx, &model.Share{
		Uid:         in.Uid,
		AssetIds:    in.AssetIds,
		Uuid:        in.Uuid,
		AssetName:   in.AssetName,
		AssetSize:   in.AssetSize,
		AssetType:   in.AssetType,
		Link:        in.Link,
		CreatedTime: time.Now().Unix(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.Response{}, nil
}
