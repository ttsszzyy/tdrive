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

type BingdUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBingdUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BingdUserLogic {
	return &BingdUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BingdUserLogic) BingdUser(req *types.BingUserReq) (resp *types.Response, err error) {
	var (
		msg string
		lan = fmt.Sprintf("%s", l.ctx.Value("language"))
	)

	//获取用户信息
	userData, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(lan)
	}
	//校验是否被绑定
	one, err := l.svcCtx.Rpc.FindOneByUid(l.ctx, &pb.UidReq{Uid: userData.User.ID})
	if err != nil {
		return nil, err
	}
	if one.Puid > 0 {
		switch lan {
		case errors.LanEn:
			msg = "Do not repeatedly bind users"
		case errors.LanTw:
			msg = "請勿重複綁定用戶"
		default:
			msg = "Do not repeatedly bind users"
		}
		return nil, errors.CustomError(msg)
	}
	//校验被绑定用户是否被绑定
	pone, err := l.svcCtx.Rpc.FindOneByUid(l.ctx, &pb.UidReq{Uid: req.Puid})
	if err != nil {
		return nil, err
	}
	if pone.Puid > 0 {
		switch lan {
		case errors.LanEn:
			msg = "The bound user has already been bound"
		case errors.LanTw:
			msg = "綁定的用戶已經被綁定"
		default:
			msg = "The bound user has already been bound"
		}
		return nil, errors.CustomError(msg)
	}
	//校验自己是否被绑定
	list, err := l.svcCtx.Rpc.FindUserByPid(l.ctx, &pb.FindUserByPidReq{Puid: userData.User.ID})
	if err != nil {
		return nil, err
	}
	if list.Total > 0 {
		switch lan {
		case errors.LanEn:
			msg = "You have already been unbound before you can be bound"
		case errors.LanTw:
			msg = "您已經被綁定解除綁定才能被綁定"
		default:
			msg = "You have already been unbound before you can be bound"
		}
		return nil, errors.CustomError(msg)
	}

	_, err = l.svcCtx.Rpc.SaveMessage(l.ctx, &pb.SaveMessageReq{
		Uid:    userData.User.ID,
		Name:   userData.User.Username,
		Puid:   req.Puid,
		Status: 1,
		Remark: "绑定",
	})
	if err != nil {
		return nil, err
	}

	return
}
