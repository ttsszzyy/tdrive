package parameter

import (
	"T-driver/app/user/rpc/pb"
	"context"

	"T-driver/app/admin/api/internal/svc"
	"T-driver/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BotCommandListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBotCommandListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BotCommandListLogic {
	return &BotCommandListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BotCommandListLogic) BotCommandList(req *types.BotCommandListReq) (resp *types.BotCommandListResp, err error) {
	command, err := l.svcCtx.Rpc.FindBotCommand(l.ctx, &pb.FindBotCommandReq{
		Command:      req.BotCommand,
		Text:         req.Text,
		SendType:     req.SendType,
		Description:  req.Description,
		LanguageCode: req.LanguageCode,
		Page:         req.Page,
		Size:         req.Size,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.BotCommandListResp{
		List:  make([]*types.BotCommand, 0, len(command.List)),
		Total: command.Total,
	}
	for _, v := range command.List {
		resp.List = append(resp.List, &types.BotCommand{
			BotCommand:   v.Command,
			CreateTime:   v.CreatedTime,
			Description:  v.Description,
			Id:           v.Id,
			SendType:     v.SendType,
			Stauts:       v.Status,
			Text:         v.Text,
			UpdateTime:   v.UpdatedTime,
			LanguageCode: v.LanguageCode,
		})
	}

	return resp, nil
}
