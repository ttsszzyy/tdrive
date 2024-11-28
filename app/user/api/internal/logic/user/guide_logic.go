package user

import (
	"T-driver/app/user/api/internal/svc"
	"T-driver/app/user/api/internal/types"
	"T-driver/app/user/model"
	"T-driver/app/user/rpc/pb"
	"T-driver/common/db"
	"T-driver/common/errors"
	"context"
	"fmt"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type GuideLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGuideLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GuideLogic {
	return &GuideLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GuideLogic) Guide(req *types.GuideReq) (resp *types.GuideResponse, err error) {
	var tgIntegral int64

	resp = &types.GuideResponse{}
	lan := fmt.Sprintf("%s", l.ctx.Value("language"))
	//获取用户信息
	userData, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(lan)
	}

	user, err := l.svcCtx.Rpc.FindOneByUid(l.ctx, &pb.UidReq{Uid: userData.User.ID})
	if err != nil {
		return nil, errors.DbError(lan)
	}
	//校验奖励是否领取
	if user.IsReceive != 1 {
		//添加分布式锁，防止重复领取
		err = l.svcCtx.Redis.Setex("user_guide_"+strconv.FormatInt(userData.User.ID, 10), "1", 2)
		if err != nil {
			return nil, errors.CustomError("Please try again later")
		}
		defer l.svcCtx.Redis.Del("user_guide_" + strconv.FormatInt(userData.User.ID, 10))

		//获取注册年限 获取TG会员
		dictList, err := l.svcCtx.Rpc.FindDictByName(l.ctx, &pb.FindDictByNameReq{Code: []string{model.RegisterYearDictCode, model.RegisterTGDictCode, model.RegisterRewardDictCode}})
		if err != nil {
			return nil, err
		}

		if dictList != nil && len(dictList.DictList) != 3 {
			return nil, errors.CustomError("未配置注册奖励")
		}
		//获取年限
		yearIntegral, _ := strconv.ParseInt(dictList.DictList[0].Value, 10, 64)
		//获取tg会员
		if userData.User.IsPremium {
			tgIntegral, _ = strconv.ParseInt(dictList.DictList[1].Value, 10, 64)
		}
		//随机积分
		randomIntegral, _ := strconv.ParseInt(dictList.DictList[2].Value, 10, 64)
		if req.LuckyPoints > randomIntegral {
			return nil, errors.CustomError("随机积分有误")
		}
		getCtx, err := l.svcCtx.Redis.GetCtx(l.ctx, model.UserReward+strconv.FormatInt(userData.User.ID, 10))
		if err != nil || err == redis.Nil {
			return nil, err
		}
		luckyPoints, _ := strconv.ParseInt(getCtx, 10, 64)
		if luckyPoints == 0 {
			luckyPoints = req.LuckyPoints
		}

		year := getTgYearGear(userData.User.ID)

		//用户奖励
		_, err = l.svcCtx.Rpc.ReceiveUser(l.ctx, &pb.ReceiveUserReq{
			Uid:       user.Uid,
			Points:    newUserFixedPoints + luckyPoints + tgIntegral + int64(float64(yearIntegral)*year),
			IsReceive: 1,
		})
		if err != nil {
			return nil, errors.DbError(lan)
		}
		l.svcCtx.Redis.DelCtx(l.ctx, model.UserReward+strconv.FormatInt(userData.User.ID, 10))
	}
	resp.IsReceive = true

	return resp, nil
}
