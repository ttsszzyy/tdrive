package logic

import (
	"T-driver/app/user/rpc/user"
	"bytes"
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"T-driver/app/tgbot/rpc/internal/svc"
	"T-driver/app/tgbot/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendPhotoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendPhotoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendPhotoLogic {
	return &SendPhotoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 文本发送
func (l *SendPhotoLogic) SendPhoto(in *pb.SendPhotoRequest) (*pb.SendBotMsgResponse, error) {
	msg := &models.Message{}
	var err error
	if len(in.ReplyMarkup) > 0 {
		InlineKeyboardMarkup := make([][]models.InlineKeyboardButton, 0, len(in.ReplyMarkup))
		for _, markup := range in.ReplyMarkup {
			inlineKeyboardMarkup := make([]models.InlineKeyboardButton, 0, len(markup.Markup))
			for _, item := range markup.Markup {
				inlineKeyboardMarkup = append(inlineKeyboardMarkup, models.InlineKeyboardButton{
					Text:         item.Button,
					URL:          item.Url,
					CallbackData: item.CallbackData,
				})
			}
			InlineKeyboardMarkup = append(InlineKeyboardMarkup, inlineKeyboardMarkup)
		}
		msg, err = l.svcCtx.TgBot.SendPhoto(l.ctx, &bot.SendPhotoParams{
			ChatID:      in.ChatId,
			Photo:       &models.InputFileUpload{Data: bytes.NewReader(in.Photo)},
			Caption:     in.Caption,
			ReplyMarkup: models.InlineKeyboardMarkup{InlineKeyboard: InlineKeyboardMarkup},
		})
		if err != nil {
			logx.Error(err)
			return nil, err
		}
	} else {
		msg, err = l.svcCtx.TgBot.SendPhoto(l.ctx, &bot.SendPhotoParams{
			ChatID:  in.ChatId,
			Photo:   &models.InputFileUpload{Data: bytes.NewReader(in.Photo)},
			Caption: in.Caption,
		})
		if err != nil {
			logx.Error(err)
			return nil, err
		}
	}
	if in.IsPinChatMessage {
		_, err = l.svcCtx.TgBot.PinChatMessage(l.ctx, &bot.PinChatMessageParams{ChatID: msg.Chat.ID, MessageID: msg.ID})
		if err != nil {
			logx.Error(err)
			return nil, err
		}
		l.svcCtx.Rpc.SaveBotPinMsg(l.ctx, &user.SaveBotPinMsgReq{ChatId: msg.Chat.ID, Msg: int64(msg.ID), Text: in.Caption})
	}

	return &pb.SendBotMsgResponse{}, nil
}
