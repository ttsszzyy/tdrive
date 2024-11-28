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

type SetLanguageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetLanguageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetLanguageLogic {
	return &SetLanguageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetLanguageLogic) SetLanguage(req *types.SetLanguageReq) (resp *types.Response, err error) {
	//获取用户信息
	userData, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(fmt.Sprintf("%s", l.ctx.Value("language")))
	}

	user, err := l.svcCtx.Rpc.FindOneByUid(l.ctx, &pb.UidReq{Uid: userData.User.ID})
	if err != nil {
		return nil, errors.DbError()
	}
	_, err = l.svcCtx.Rpc.SaveUser(l.ctx, &pb.User{Id: user.Id, LanguageCode: req.LanguageCode})
	if err != nil {
		logx.Error(err)
		return nil, errors.DbError()
	}

	return &types.Response{}, nil
}
