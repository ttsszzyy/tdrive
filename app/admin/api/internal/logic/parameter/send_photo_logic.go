package parameter

import (
	"T-driver/app/tgbot/rpc/pb"
	"T-driver/common/errors"
	"context"
	"fmt"

	"T-driver/app/admin/api/internal/svc"
	"T-driver/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendPhotoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendPhotoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendPhotoLogic {
	return &SendPhotoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendPhotoLogic) SendPhoto(req *types.SendPhotoReq) (resp *types.Response, err error) {
	if req.FileSize > 10*1024*1024*1024 {
		return nil, errors.CustomError(fmt.Sprintf("文件大小不能超过%d MB", 10))
	}
	replyMarkup := make([]*pb.Markup, 0, len(req.ReplyMarkup))
	for _, v := range req.ReplyMarkup {
		markup := &pb.Markup{
			Markup: make([]*pb.Item, 0, len(v)),
		}
		for _, m := range v {
			markup.Markup = append(markup.Markup, &pb.Item{Button: m.Button, Url: m.Url})
		}
		replyMarkup = append(replyMarkup, markup)
	}
	_, err = l.svcCtx.BotRpc.SendPhoto(l.ctx, &pb.SendPhotoRequest{
		ChatId:           req.ChatID,
		Photo:            req.Photo,
		Caption:          req.Caption,
		ReplyMarkup:      replyMarkup,
		IsPinChatMessage: req.IsPinChatMessage,
	})
	if err != nil {
		return nil, errors.FromRpcError(err)
	}
	return &types.Response{}, nil
}
