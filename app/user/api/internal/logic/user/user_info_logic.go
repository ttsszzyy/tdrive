package user

import (
	"T-driver/app/user/api/internal/svc"
	"T-driver/app/user/api/internal/types"
	"T-driver/app/user/rpc/pb"
	"T-driver/common/db"
	"T-driver/common/errors"
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo(req *types.Request) (resp *types.UserRes, err error) {
	resp = &types.UserRes{}
	//获取用户信息
	userData, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(fmt.Sprintf("%s", l.ctx.Value("language")))
	}
	user, err := l.svcCtx.Rpc.FindOneByUid(l.ctx, &pb.UidReq{Uid: userData.User.ID})
	if err != nil {
		return nil, errors.DbError()
	}
	p, err := l.svcCtx.Rpc.FindOneUserPoints(l.ctx, &pb.FindOneUserPointsReq{Uid: userData.User.ID})
	if err != nil {
		return nil, errors.DbError()
	}
	s, err := l.svcCtx.Rpc.FindOneUserStorage(l.ctx, &pb.FindOneUserStorageReq{Uid: userData.User.ID})
	if err != nil {
		return nil, errors.DbError()
	}
	u, err := l.svcCtx.Rpc.CheckIsOldUser(l.ctx, &pb.UidReq{Uid: userData.User.ID})
	if err != nil {
		return nil, errors.DbError()
	}
	resp.IsTips = u.IsOldUser
	resp.Uid = user.Uid
	resp.Name = user.Name
	resp.Avatar = user.Avatar
	resp.Integral = p.Points
	resp.Storage = s.Storage
	resp.StorageUse = s.StorageUse
	resp.LanguageCode = user.LanguageCode
	resp.IsRead = user.IsRead > 0
	return resp, nil
}
