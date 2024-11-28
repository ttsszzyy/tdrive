package logic

import (
	"T-driver/app/user/model"
	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindOneByUidLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindOneByUidLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindOneByUidLogic {
	return &FindOneByUidLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户端
func (l *FindOneByUidLogic) FindOneByUid(in *pb.UidReq) (*pb.User, error) {
	user, err := l.svcCtx.UserModel.FindOneByUidDeletedTime(l.ctx, in.Uid, 0)
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
		IsRead:        user.IsRead,
	}, err
}
