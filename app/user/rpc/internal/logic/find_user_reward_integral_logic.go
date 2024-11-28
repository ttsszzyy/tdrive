package logic

import (
	"T-driver/app/user/model"
	"context"

	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindUserRewardIntegralLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindUserRewardIntegralLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindUserRewardIntegralLogic {
	return &FindUserRewardIntegralLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查看邀请奖励
func (l *FindUserRewardIntegralLogic) FindUserRewardIntegral(in *pb.FindUserRewardIntegralReq) (*pb.FindUserRewardIntegralResp, error) {
	userInviteReward, err := l.svcCtx.UserInviteRewardModel.FindOneByUidDeletedTime(l.ctx, in.Uid, 0)
	if err != nil {
		if err == model.ErrNotFound {
			return &pb.FindUserRewardIntegralResp{}, nil
		}
		return nil, err
	}
	return &pb.FindUserRewardIntegralResp{
		Id:          userInviteReward.Id,
		Uid:         userInviteReward.Uid,
		InviteNum:   userInviteReward.InviteNum,
		CreatedTime: userInviteReward.CreatedTime,
	}, nil
}
