package logic

import (
	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"
	"context"
	"github.com/Masterminds/squirrel"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindUserInviteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindUserInviteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindUserInviteLogic {
	return &FindUserInviteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 朋友列表
func (l *FindUserInviteLogic) FindUserInvite(in *pb.FindUserInviteReq) (*pb.FindUserInviteResp, error) {
	sb := squirrel.Select().Where("deleted_time = ?", 0)
	if in.Pid > 0 {
		sb = sb.Where(squirrel.Eq{"pid": in.Pid})
	}
	if len(in.Uid) > 0 {
		sb = sb.Where(squirrel.Eq{"uid": in.Uid})
	}

	list, err := l.svcCtx.UserInviteModel.List(l.ctx, sb)
	if err != nil {
		return nil, err
	}
	p := &pb.FindUserInviteResp{List: make([]*pb.UserInvite, 0, len(list))}
	for _, v := range list {
		p.List = append(p.List, &pb.UserInvite{
			Uid:          v.Uid,
			Pid:          v.Pid,
			CreatedTime:  v.CreatedTime,
			Id:           v.Id,
			InvitePoints: v.InvitePoints,
		})
	}

	return p, nil
}
