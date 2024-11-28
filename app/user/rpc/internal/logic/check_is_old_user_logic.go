package logic

import (
	"context"

	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"
	"T-driver/common/errors"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CheckIsOldUserLogic 校验当前用户是否为老用户
type CheckIsOldUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// NewCheckIsOldUserLogic 新建 校验当前用户是否为老用户
func NewCheckIsOldUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckIsOldUserLogic {
	return &CheckIsOldUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CheckIsOldUser 实现 校验当前用户是否为老用户
func (l *CheckIsOldUserLogic) CheckIsOldUser(in *pb.UidReq) (*pb.IsOldUserResp, error) {
	rsp, err := l.svcCtx.UserModel.CheckIsOldUer(l.ctx, in.Uid)
	if err != nil {
		return nil, status.Errorf(codes.Internal, errors.DbError().Msg())
	}

	return &pb.IsOldUserResp{
		IsOldUser: rsp,
	}, nil
}
