package message

import (
	"T-driver/app/user/rpc/pb"
	"T-driver/common/db"
	"T-driver/common/errors"
	"context"
	"fmt"

	"T-driver/app/user/api/internal/svc"
	"T-driver/app/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MessageListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMessageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MessageListLogic {
	return &MessageListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MessageListLogic) MessageList(req *types.RewardListReq) (resp *types.MessageListResp, err error) {
	//获取用户信息
	userData, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(fmt.Sprintf("%s", l.ctx.Value("language")))
	}
	if req.Size > 100 {
		req.Size = 100
	}
	message, err := l.svcCtx.Rpc.FindMessage(l.ctx, &pb.FindMessageReq{
		Puid: userData.User.ID,
		Page: req.Page,
		Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.MessageListResp{
		Total: message.Total,
		Items: make([]*types.MessageItem, 0),
	}
	for _, item := range message.Messages {
		resp.Items = append(resp.Items, &types.MessageItem{
			Id:          item.Id,
			Uid:         item.Uid,
			Name:        item.Name,
			Puid:        item.Puid,
			Remark:      item.Remark,
			Status:      item.Status,
			CreatedTime: item.CreatedTime,
			UpdatedTime: item.UpdatedTime,
		})
	}
	return
}
