package logic

import (
	"T-driver/app/user/model"
	"context"

	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserTitanTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserTitanTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserTitanTokenLogic {
	return &UpdateUserTitanTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserTitanTokenLogic) UpdateUserTitanToken(in *pb.UpdateUserTitanTokenReq) (*pb.Response, error) {
	one, err := l.svcCtx.UserTitanTokenModel.FindOneByUid(l.ctx, in.Uid)
	if err != nil {
		if err == model.ErrNotFound {
			err := l.svcCtx.UserTitanTokenModel.Insert(l.ctx, &model.UserTitanToken{
				Uid:   in.Uid,
				Token: in.Token,
				Exp:   in.Expire,
			})
			return &pb.Response{}, err
		}
		return nil, err
	}
	if in.Token != "" {
		one.Token = in.Token
	}
	if in.Expire > 0 {
		one.Exp = in.Expire
	}
	_, err = l.svcCtx.UserTitanTokenModel.UpdateNoCache(l.ctx, one)
	if err != nil {
		return nil, err
	}
	return &pb.Response{}, err
}
