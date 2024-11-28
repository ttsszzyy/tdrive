package logic

import (
	"T-driver/app/user/model"
	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindOneBotCommandLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindOneBotCommandLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindOneBotCommandLogic {
	return &FindOneBotCommandLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询bot命令
func (l *FindOneBotCommandLogic) FindOneBotCommand(in *pb.FindOneBotCommandReq) (*pb.BotCommandItem, error) {
	one := &model.BotCommand{}
	var err error
	if in.BotCommand != "" {
		one, err = l.svcCtx.BotCommandModel.FindOneByBotCommandLanguageCodeDeletedTime(l.ctx, in.BotCommand, in.LanguageCode, 0)
	}
	if in.Id > 0 {
		one, err = l.svcCtx.BotCommandModel.FindOne(l.ctx, in.Id)
	}
	if err != nil {
		if err == model.ErrNotFound {
			return &pb.BotCommandItem{}, nil
		}
		return nil, err
	}
	return &pb.BotCommandItem{
		Command:     one.BotCommand,
		ButtonArray: one.ButtonArray,
		CreatedTime: one.CreatedTime,
		Id:          one.Id,
		Photo:       one.Photo,
		Status:      one.Status,
		Text:        one.Text,
		UpdatedTime: one.UpdatedTime,
		Description: one.Description,
		SendType:    one.SendType,
	}, nil
}
