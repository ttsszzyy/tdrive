package user

import (
	"T-driver/app/user/api/internal/svc"
	"T-driver/app/user/api/internal/types"
	"T-driver/app/user/rpc/pb"
	"T-driver/common/db"
	"T-driver/common/errors"
	"context"
	"fmt"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/sync/singleflight"
)

type ClaimRewardLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	sg     singleflight.Group
}

func NewClaimRewardLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClaimRewardLogic {
	return &ClaimRewardLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ClaimRewardLogic) ClaimReward(req *types.Request) (resp *types.ClaimRewardResp, err error) {
	resp = &types.ClaimRewardResp{}
	//获取用户信息
	userData, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(fmt.Sprintf("%s", l.ctx.Value("language")))
	}
	//添加分布式锁，防止重复领取
	err = l.svcCtx.Redis.Setex("user_claim_reward_"+strconv.FormatInt(userData.User.ID, 10), "1", 2)
	if err != nil {
		return nil, errors.CustomError("Please try again later")
	}
	defer l.svcCtx.Redis.Del("user_claim_reward_" + strconv.FormatInt(userData.User.ID, 10))
	_, err, _ = l.sg.Do(strconv.FormatInt(userData.User.ID, 10), func() (interface{}, error) {
		_, err = l.svcCtx.Rpc.ClaimInviteReward(l.ctx, &pb.ClaimInviteRewardReq{Uid: userData.User.ID})
		if err != nil {
			logx.Error(err)
			return nil, errors.FromRpcError(err)
		}
		return nil, nil
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}
