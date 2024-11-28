package message

import (
	"T-driver/app/user/rpc/user"
	"T-driver/common/db"
	"T-driver/common/errors"
	"context"
	"fmt"

	"T-driver/app/user/api/internal/svc"
	"T-driver/app/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BingdConfirmLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBingdConfirmLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BingdConfirmLogic {
	return &BingdConfirmLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BingdConfirmLogic) BingdConfirm(req *types.BingConfirmReq) (resp *types.Response, err error) {
	//获取用户信息
	_, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(fmt.Sprintf("%s", l.ctx.Value("language")))
	}
	_, err = l.svcCtx.Rpc.SaveMessage(l.ctx, &user.SaveMessageReq{
		Id:     req.Id,
		Status: req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &types.Response{}, nil
}
