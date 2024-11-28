package trade

import (
	"context"
	"fmt"

	"T-driver/app/user/api/internal/svc"
	"T-driver/app/user/api/internal/types"
	"T-driver/app/user/rpc/pb"
	"T-driver/common/db"
	"T-driver/common/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAirdropTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAirdropTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAirdropTokenLogic {
	return &GetAirdropTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAirdropTokenLogic) GetAirdropToken() (resp *types.AirdropTokenDetail, err error) {
	//获取用户信息
	userData, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(fmt.Sprintf("%s", l.ctx.Value("language")))
	}

	resp = new(types.AirdropTokenDetail)

	info, err := l.svcCtx.Rpc.GetUserStorageAndToken(l.ctx, &pb.UidReq{Uid: userData.User.ID})
	if err != nil {
		return nil, errors.ErrorNotFound(fmt.Sprintf("%s", l.ctx.Value("language")))
	}
	resp.BookedToken = info.Token
	resp.BookableToken = info.Storage

	return resp, nil
}
