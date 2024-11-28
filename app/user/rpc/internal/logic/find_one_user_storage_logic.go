package logic

import (
	"T-driver/app/user/model"
	"context"

	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindOneUserStorageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindOneUserStorageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindOneUserStorageLogic {
	return &FindOneUserStorageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户存储空间
func (l *FindOneUserStorageLogic) FindOneUserStorage(in *pb.FindOneUserStorageReq) (*pb.FindOneUserStorageResp, error) {
	one, err := l.svcCtx.UserStorageModel.FindOneByUidDeletedTime(l.ctx, in.Uid, 0)
	if err != nil {
		if err == model.ErrNotFound {
			return &pb.FindOneUserStorageResp{}, nil
		}
		return nil, err
	}

	return &pb.FindOneUserStorageResp{
		Id:         one.Id,
		Uid:        one.Uid,
		Storage:    one.Storage,
		StorageUse: one.StorageUse,
		SurStorage: one.Storage - one.StorageUse,
	}, nil
}
