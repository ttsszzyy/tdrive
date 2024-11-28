package logic

import (
	"T-driver/app/user/model"
	"context"
	"time"

	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveBotCommandLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveBotCommandLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveBotCommandLogic {
	return &SaveBotCommandLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 保存bot命令
func (l *SaveBotCommandLogic) SaveBotCommand(in *pb.SaveBotCommandReq) (*pb.Response, error) {
	_, err := l.svcCtx.BotCommandModel.Insert(l.ctx, &model.BotCommand{
		BotCommand:   in.BotCommand,
		Text:         in.Text,
		Photo:        in.Photo,
		ButtonArray:  in.ButtonArray,
		Status:       in.Status,
		CreatedTime:  time.Now().Unix(),
		Description:  in.Description,
		SendType:     in.SendType,
		LanguageCode: in.LanguageCode,
	})
	if err != nil {
		return nil, err
	}

	return &pb.Response{}, nil
}
