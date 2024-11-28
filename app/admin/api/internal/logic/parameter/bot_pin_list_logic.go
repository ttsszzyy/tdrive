package parameter

import (
	"T-driver/app/user/rpc/user"
	"T-driver/common/errors"
	"context"

	"T-driver/app/admin/api/internal/svc"
	"T-driver/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BotPinListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBotPinListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BotPinListLogic {
	return &BotPinListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BotPinListLogic) BotPinList(req *types.BotPinListReq) (resp *types.BotPinListResp, err error) {
	msg, err := l.svcCtx.Rpc.FindBotPinMsg(l.ctx, &user.FindBotPinMsgReq{ChatId: req.ChatID, Page: req.Page, Size: req.Size})
	if err != nil {
		return nil, errors.FromRpcError(err)
	}
	resp = &types.BotPinListResp{
		List:  make([]*types.BotPin, 0, len(msg.List)),
		Total: msg.Total,
	}
	for _, v := range msg.List {
		resp.List = append(resp.List, &types.BotPin{
			Id:     v.Id,
			ChatID: v.ChatId,
			MsgID:  v.Msg,
			Text:   v.Text,
		})
	}
	return resp, nil
}
