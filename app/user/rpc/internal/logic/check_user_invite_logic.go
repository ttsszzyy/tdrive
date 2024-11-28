package logic

import (
	"context"
	"github.com/Masterminds/squirrel"

	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckUserInviteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckUserInviteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckUserInviteLogic {
	return &CheckUserInviteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 校验用户是否邀请过
func (l *CheckUserInviteLogic) CheckUserInvite(in *pb.CheckUserInviteReq) (*pb.CheckUserInviteResp, error) {
	total, err := l.svcCtx.UserModel.Count(l.ctx, squirrel.Select().Where(squirrel.Eq{"pid": in.Pid, "uid": in.Uid}))
	if err != nil {
		return nil, err
	}

	return &pb.CheckUserInviteResp{Total: total}, nil
}
