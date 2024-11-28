package trade

import (
	"context"
	"fmt"
	"strconv"

	"T-driver/app/user/api/internal/svc"
	"T-driver/app/user/api/internal/types"
	"T-driver/app/user/rpc/pb"
	"T-driver/common/db"
	"T-driver/common/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExchangeTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExchangeTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExchangeTokenLogic {
	return &ExchangeTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExchangeTokenLogic) ExchangeToken(req *types.ExChangeTokenReq) error {
	//获取用户信息
	userData, ok := db.CtxInitData(l.ctx)
	if !ok {
		return errors.UnauthError(fmt.Sprintf("%s", l.ctx.Value("language")))
	}

	//添加分布式锁，防止重复
	err := l.svcCtx.Redis.Setex("user_exchange_token_"+strconv.FormatInt(userData.User.ID, 10), "1", 2)
	if err != nil {
		return errors.CustomError("Please try again later")
	}
	defer l.svcCtx.Redis.Del("user_exchange_token_" + strconv.FormatInt(userData.User.ID, 10))

	_, err = l.svcCtx.Rpc.SaveUserToken(l.ctx, &pb.SaveUserTokenReq{Uid: userData.User.ID, Token: req.Storage})
	if err != nil {
		return errors.ErrorNotFound(fmt.Sprintf("%s", l.ctx.Value("language")))
	}

	return nil
}
