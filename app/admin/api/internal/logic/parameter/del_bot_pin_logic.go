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

type DelBotPinLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDelBotPinLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelBotPinLogic {
	return &DelBotPinLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DelBotPinLogic) DelBotPin(req *types.DelBotPinReq) (resp *types.Response, err error) {
	pinMsg, err := l.svcCtx.BotRpc.DelPinMsg(l.ctx, &tgbot.DelPinMsgReq{ChatID: req.ChatID, MsgID: req.MsgID})
	if err != nil {
		return nil, errors.FromRpcError(err)
	}
	if pinMsg.IsSuccess {
		_, err = l.svcCtx.Rpc.DelBotPinMsg(l.ctx, &pb.DelBotPinMsgReq{Id: req.Id})
		if err != nil {
			return nil, errors.FromRpcError(err)
		}
	}
	return &types.Response{}, nil
}
