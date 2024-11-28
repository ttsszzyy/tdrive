package logic

import (
	"T-driver/app/user/model"
	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
)

type SaveUserTitanTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveUserTitanTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveUserTitanTokenLogic {
	return &SaveUserTitanTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SaveUserTitanTokenLogic) SaveUserTitanToken(in *pb.SaveUserTitanTokenReq) (*pb.Response, error) {
	err := l.svcCtx.UserTitanTokenModel.Insert(l.ctx, &model.UserTitanToken{
		Uid:   in.Uid,
		Token: in.Token,
		Exp:   in.Expire,
	})
	return &pb.Response{}, err
}
