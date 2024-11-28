package logic

import (
	"context"

	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DisableUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDisableUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DisableUserLogic {
	return &DisableUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 禁用用户
func (l *DisableUserLogic) DisableUser(in *pb.DisableUserReq) (*pb.Response, error) {
	one, err := l.svcCtx.UserModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	one.IsDisable = in.IsDisable
	err = l.svcCtx.UserModel.Update(l.ctx, one)
	if err != nil {
		return nil, err
	}

	return &pb.Response{}, nil
}
