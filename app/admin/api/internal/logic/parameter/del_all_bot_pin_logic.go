package parameter

import (
	"T-driver/app/tgbot/rpc/tgbot"
	"T-driver/app/user/rpc/pb"
	"T-driver/common/errors"
	"context"

	"T-driver/app/admin/api/internal/svc"
	"T-driver/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelAllBotPinLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDelAllBotPinLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelAllBotPinLogic {
	return &DelAllBotPinLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DelAllBotPinLogic) DelAllBotPin(req *types.DelAllBotPinReq) (resp *types.Response, err error) {
	pinMsg, err := l.svcCtx.BotRpc.DelAllPinMsg(l.ctx, &tgbot.DelAllPinMsgReq{ChatID: req.ChatID})
	if err != nil {
		return nil, errors.FromRpcError(err)
	}
	if pinMsg.IsSuccess {
		_, err = l.svcCtx.Rpc.DelBotPinMsg(l.ctx, &pb.DelBotPinMsgReq{ChatId: req.ChatID})
		if err != nil {
			return nil, errors.FromRpcError(err)
		}
	}
	return &types.Response{}, nil
}
