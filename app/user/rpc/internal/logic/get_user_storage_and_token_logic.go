package logic

import (
	"context"
	"fmt"

	"T-driver/app/user/model"
	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetUserStorageAndTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserStorageAndTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserStorageAndTokenLogic {
	return &GetUserStorageAndTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户预订空投
func (l *GetUserStorageAndTokenLogic) GetUserStorageAndToken(in *pb.UidReq) (*pb.UserStorageAndToken, error) {
	var resp = new(pb.UserStorageAndToken)

	// 获取用户的存储空间
	storage, err := l.svcCtx.UserStorageModel.FindOneByUidDeletedTime(l.ctx, in.Uid, 0)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Errorf("get user storage's info error:%w", err).Error())
	}

	// 获取用户的预订积分
	token, err := l.svcCtx.UserTokenModel.FindOneByUidDeletedTime(l.ctx, in.Uid, 0)
	switch err {
	case model.ErrNotFound:
		resp.Storage = storage.Storage
	case nil:
		resp.Token = token.Token
		resp.Storage = storage.Storage - token.Token
	default:
		return nil, status.Errorf(codes.Internal, fmt.Errorf("get user token's info error:%w", err).Error())
	}

	return resp, nil
}
