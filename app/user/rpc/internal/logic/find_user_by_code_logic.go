package logic

import (
	"T-driver/app/user/model"
	"context"

	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindUserByCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindUserByCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindUserByCodeLogic {
	return &FindUserByCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 根据用户邀请码查询用户
func (l *FindUserByCodeLogic) FindUserByCode(in *pb.UserCodeReq) (*pb.User, error) {
	user, err := l.svcCtx.UserModel.FindOneByRecommendCodeDeletedTime(l.ctx, in.Code, 0)
	if err != nil {
		if err == model.ErrNotFound {
			return &pb.User{}, nil
		}
		return nil, err
	}

	return &pb.User{
		Id:            user.Id,
		Uid:           user.Uid,
		Name:          user.Name,
		Avatar:        user.Avatar,
		Mail:          user.Mail,
		WalletAddress: user.WalletAddress,
		Source:        user.Source,
		RecommendCode: user.RecommendCode,
		Distribution:  user.Distribution,
		Pid:           user.Pid,
		IsDisable:     user.IsDisable,
		CreatedTime:   user.CreatedTime,
		UpdatedTime:   user.UpdatedTime,
		Puid:          user.Puid,
		IsReceive:     user.IsReceive,
		LanguageCode:  user.LanguageCode,
	}, err
}
