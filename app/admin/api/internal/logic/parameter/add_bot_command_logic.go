package parameter

import (
	"T-driver/app/tgbot/rpc/pb"
	"T-driver/app/user/rpc/user"
	"T-driver/common/utils/base64"
	"context"
	"encoding/json"

	"T-driver/app/admin/api/internal/svc"
	"T-driver/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddBotCommandLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddBotCommandLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddBotCommandLogic {
	return &AddBotCommandLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddBotCommandLogic) AddBotCommand(req *types.AddBotCommandReq) (resp *types.Response, err error) {
	if len(req.ReplyMarkup) > 0 {
		replyMarkup := make([]*pb.Markup, 0, len(req.ReplyMarkup))
		for _, v := range req.ReplyMarkup {
			markup := &pb.Markup{
				Markup: make([]*pb.Item, 0, len(v)),
			}
			for _, m := range v {
				markup.Markup = append(markup.Markup, &pb.Item{Button: m.Button, Url: m.Url, CallbackData: m.CallbackData})
			}
			replyMarkup = append(replyMarkup, markup)
		}
	}
	if req.SendType != 1 {
		command, err := l.svcCtx.Rpc.FindBotCommand(l.ctx, &user.FindBotCommandReq{})
		if err != nil {
			return nil, err
		}
		botCommand := make([]*pb.BotCommand, 0, len(command.List)+1)
		botCommandMap := make(map[string]string)
		for _, v := range command.List {
			if v.SendType != 1 {
				botCommandMap[v.Command] = v.Description
			}
		}
		botCommandMap[req.BotCommand] = req.Description
		for k, v := range botCommandMap {
			botCommand = append(botCommand, &pb.BotCommand{
				Command:     k,
				Description: v,
			})
		}
		_, err = l.svcCtx.BotRpc.SendBotCommand(l.ctx, &pb.SendBotCommandReq{BotCommand: botCommand})
		if err != nil {
			return nil, err
		}
	}
	marshal, _ := json.Marshal(req.ReplyMarkup)
	encode := ""
	if len(req.Photo) > 0 {
		encode = string(base64.Encode(req.Photo))
	}
	_, err = l.svcCtx.Rpc.SaveBotCommand(l.ctx, &user.SaveBotCommandReq{
		BotCommand:   req.BotCommand,
		LanguageCode: req.LanguageCode,
		Description:  req.Description,
		Photo:        encode,
		ButtonArray:  string(marshal),
		Text:         req.Text,
		SendType:     req.SendType,
	})
	if err != nil {
		return nil, err
	}

	return &types.Response{}, nil
}
