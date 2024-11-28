package logic

import (
	"context"
	"github.com/Masterminds/squirrel"

	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CountUserByExchangeStorageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCountUserByExchangeStorageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CountUserByExchangeStorageLogic {
	return &CountUserByExchangeStorageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询奖励空间总数
func (l *CountUserByExchangeStorageLogic) CountUserByExchangeStorage(in *pb.CountUserByExchangeStorageReq) (*pb.CountUserByExchangeStorageResp, error) {
	storage, err := l.svcCtx.UserModel.CountExchangeStorage(l.ctx, squirrel.Select().Where("deleted_time = ?", 0))
	if err != nil {
		return nil, err
	}

	return &pb.CountUserByExchangeStorageResp{Total: storage}, nil
}
