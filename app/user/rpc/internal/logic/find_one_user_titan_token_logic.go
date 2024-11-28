package logic

import (
	"T-driver/app/user/model"
	"context"

	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindOneUserTitanTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindOneUserTitanTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindOneUserTitanTokenLogic {
	return &FindOneUserTitanTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindOneUserTitanTokenLogic) FindOneUserTitanToken(in *pb.FindOneUserTitanTokenReq) (*pb.UserTitanToken, error) {
	one, err := l.svcCtx.UserTitanTokenModel.FindOneByUid(l.ctx, in.Uid)
	if err != nil {
		if err == model.ErrNotFound {
			return &pb.UserTitanToken{}, nil
		}
		return nil, err
	}
	return &pb.UserTitanToken{
		Uid:    one.Uid,
		Token:  one.Token,
		Expire: one.Exp,
	}, nil
}
