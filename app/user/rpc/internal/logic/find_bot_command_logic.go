package logic

import (
	"T-driver/app/user/model"
	"context"
	"github.com/Masterminds/squirrel"

	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindBotCommandLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindBotCommandLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindBotCommandLogic {
	return &FindBotCommandLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询bot命令
func (l *FindBotCommandLogic) FindBotCommand(in *pb.FindBotCommandReq) (*pb.FindBotCommandResp, error) {
	sb := squirrel.Select().Where("deleted_time = ?", 0)
	if in.Command != "" {
		sb = sb.Where("bot_command = ?", in.Command)
	}
	if in.Text != "" {
		sb = sb.Where("text = ?", in.Text)
	}
	if in.LanguageCode != "" {
		sb = sb.Where("language_code = ?", in.LanguageCode)
	}
	list := make([]*model.BotCommand, 0)
	var err error
	var total int64
	if in.Page == 0 && in.Size == 0 {
		list, err = l.svcCtx.BotCommandModel.List(l.ctx, sb)
	} else {
		list, total, err = l.svcCtx.BotCommandModel.ListPage(l.ctx, in.Page, in.Size, sb)
	}
	if err != nil {
		return nil, err
	}
	resp := &pb.FindBotCommandResp{
		Total: total,
		List:  make([]*pb.BotCommandItem, 0, len(list)),
	}
	for _, command := range list {
		resp.List = append(resp.List, &pb.BotCommandItem{
			Command:      command.BotCommand,
			ButtonArray:  command.ButtonArray,
			CreatedTime:  command.CreatedTime,
			Id:           command.Id,
			Photo:        command.Photo,
			Status:       command.Status,
			Text:         command.Text,
			UpdatedTime:  command.UpdatedTime,
			Description:  command.Description,
			SendType:     command.SendType,
			LanguageCode: command.LanguageCode,
		})
	}

	return resp, nil
}
