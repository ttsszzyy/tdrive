package user

import (
	"T-driver/app/user/rpc/pb"
	"T-driver/common/db"
	"T-driver/common/errors"
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/mr"

	"T-driver/app/user/api/internal/svc"
	"T-driver/app/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserRecommendPlanLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserRecommendPlanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRecommendPlanLogic {
	return &UserRecommendPlanLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserRecommendPlanLogic) UserRecommendPlan(req *types.Request) (resp *types.UserRecommendPlanRes, err error) {
	resp = &types.UserRecommendPlanRes{
		DictItems: make([]*types.DictItem, 0, 12),
		List:      make([]*types.Friend, 0, 10),
	}
	//获取用户信息
	userData, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(fmt.Sprintf("%s", l.ctx.Value("language")))
	}
	user := &pb.User{}
	userRewardIntegral := &pb.FindUserRewardIntegralResp{}
	userList := &pb.FindUserInviteResp{}
	err = mr.Finish(func() error {
		user, err = l.svcCtx.Rpc.FindOneByUid(l.ctx, &pb.UidReq{Uid: userData.User.ID})
		if err != nil {
			return errors.DbError()
		}
		return nil
	}, func() error {
		//获取用户邀请奖励数量
		userRewardIntegral, err = l.svcCtx.Rpc.FindUserRewardIntegral(l.ctx, &pb.FindUserRewardIntegralReq{Uid: userData.User.ID})
		if err != nil {
			logx.Error("获取用户邀请奖励数量失败：", err)
			return err
		}
		return nil
	}, func() error {
		//朋友列表
		userList, err = l.svcCtx.Rpc.FindUserInvite(l.ctx, &pb.FindUserInviteReq{Pid: userData.User.ID})
		if err != nil {
			return errors.DbError()
		}
		return nil
	})
	if err != nil {
		return nil, errors.FromRpcError(err)
	}

	err = mr.Finish(func() error {
		ids := make([]int64, 0, len(userList.List))
		umap := make(map[int64]*pb.User, len(userList.List))
		pmap := make(map[int64]*pb.FindOneUserPointsResp, len(userList.List))
		for _, v := range userList.List {
			ids = append(ids, v.Uid)
		}
		if len(ids) > 0 {
			findUser, err := l.svcCtx.Rpc.FindUser(l.ctx, &pb.QueryUserReq{Uid: ids})
			if err != nil {
				return errors.DbError()
			}
			for _, u := range findUser.Users {
				umap[u.Uid] = u
			}
			findPoints, err := l.svcCtx.Rpc.FindUserPoints(l.ctx, &pb.FindUserPointsReq{Uid: ids})
			if err != nil {
				return errors.DbError()
			}
			for _, p := range findPoints.List {
				pmap[p.Uid] = p
			}
		}
		p, err := l.svcCtx.Rpc.FindOneUserPoints(l.ctx, &pb.FindOneUserPointsReq{Uid: userData.User.ID})
		if err != nil {
			return errors.DbError()
		}

		//获取推荐用户数据
		resp.RewardPoints = p.ReqPoints
		resp.Total = int64(len(userList.List))
		for _, v := range userList.List {
			u := umap[v.Uid]
			var points int64
			po, ok := pmap[v.Uid]
			if ok {
				points = po.Points
			}
			//todo 用户奖励积分
			resp.List = append(resp.List, &types.Friend{
				Uid:          u.Uid,
				Name:         u.Name,
				Points:       points,
				RewardPoints: v.InvitePoints,
				CreatedTime:  u.CreatedTime,
			})
		}
		resp.Url = l.svcCtx.Config.FastReward.TgUrl + "?" + "startapp=" + user.RecommendCode
		resp.Describe = l.svcCtx.Config.FastReward.Describe
		return nil
	}, func() error {
		//获取活动规则
		//引荐奖励
		list, err := l.svcCtx.Rpc.FindDictByName(l.ctx, &pb.FindDictByNameReq{ParamType: []int64{2, 5}})
		if err != nil {
			return errors.DbError()
		}
		if len(list.DictList) == 0 {
			return errors.CustomError("未配置引荐奖励")
		}
		for _, dict := range list.DictList {
			resp.DictItems = append(resp.DictItems, &types.DictItem{
				Id:          dict.Id,
				ParamType:   dict.ParamType,
				Name:        dict.Name,
				Value:       dict.Value,
				BackupValue: dict.BackupValue,
				Code:        dict.Code,
				Status:      userRewardIntegral.InviteNum >= dict.BackupValue,
			})
		}
		return nil
	})
	if err != nil {
		return nil, errors.FromRpcError(err)
	}

	return
}
