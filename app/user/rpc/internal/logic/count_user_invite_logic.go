package logic

import (
	"context"
	"github.com/Masterminds/squirrel"

	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CountUserInviteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCountUserInviteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CountUserInviteLogic {
	return &CountUserInviteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询用戶邀請數量
func (l *CountUserInviteLogic) CountUserInvite(in *pb.CountUserInviteReq) (*pb.CountUserInviteResp, error) {
	sb := squirrel.Select().Where("pid = ?", in.Pid)
	count, err := l.svcCtx.UserModel.CountInvite(l.ctx, sb)
	if err != nil {
		return nil, err
	}
	return &pb.CountUserInviteResp{Total: count}, nil
}
