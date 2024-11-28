package user

import (
	"T-driver/app/user/model"
	"T-driver/app/user/rpc/pb"
	"T-driver/common/db"
	"T-driver/common/errors"
	"T-driver/common/utils/rand"
	"context"
	"fmt"
	"strconv"
	"unicode/utf8"

	"github.com/zeromicro/go-zero/core/stores/redis"

	"T-driver/app/user/api/internal/svc"
	"T-driver/app/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	newUserFixedPoints int64 = 10_0000 // 新用户固定积分
	luckyBasePoints    int64 = 1_0000  // 随机幸运积分基础积分
)

type GuideRewardLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGuideRewardLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GuideRewardLogic {
	return &GuideRewardLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GuideRewardLogic) GuideReward(req *types.Request) (resp *types.GuideRewardRes, err error) {
	var (
		tgIntegral int64
	)

	resp = &types.GuideRewardRes{}
	//获取用户信息
	userData, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(fmt.Sprintf("%s", l.ctx.Value("language")))
	}
	//获取注册年限 获取TG会员
	dictList, err := l.svcCtx.Rpc.FindDictByName(l.ctx, &pb.FindDictByNameReq{Code: []string{model.RegisterYearDictCode, model.RegisterTGDictCode, model.RegisterRewardDictCode}})
	if err != nil {
		return nil, err
	}

	if dictList != nil && len(dictList.DictList) != 3 {
		return nil, errors.CustomError("未配置注册奖励")
	}
	// 获取年限的积分配置
	yearIntegral, _ := strconv.ParseInt(dictList.DictList[0].Value, 10, 64)
	// 获取tg会员
	if userData.User.IsPremium {
		tgIntegral, _ = strconv.ParseInt(dictList.DictList[1].Value, 10, 64)
	}

	//随机积分
	random, err := l.svcCtx.Redis.GetCtx(l.ctx, model.UserReward+strconv.FormatInt(userData.User.ID, 10))
	if err != nil && err != redis.Nil {
		return nil, err
	}
	randomIntegral, _ := strconv.ParseInt(random, 10, 64)
	if random == "" {
		value, _ := strconv.ParseInt(dictList.DictList[2].Value, 10, 64)
		randomIntegral = luckyBasePoints + rand.RandomInt64(0, value)
		//缓存积分
		l.svcCtx.Redis.SetCtx(l.ctx, model.UserReward+strconv.FormatInt(userData.User.ID, 10), strconv.FormatInt(randomIntegral, 10))
	}

	resp.Uid = userData.User.ID
	resp.Year = getTgYearGear(userData.User.ID)
	resp.Top = 50
	resp.YearPoints = int64(float64(yearIntegral) * resp.Year)
	resp.LuckyPoints = randomIntegral //幸运空间
	resp.TgPoints = tgIntegral
	resp.NewUser = newUserFixedPoints
	return resp, nil
}

func getTgYearGear(tgID int64) float64 {
	count := utf8.RuneCountInString(strconv.FormatInt(tgID, 10))

	switch {
	case count <= 8:
		return 2
	case count == 9:
		return 1.5
	default:
		return 1
	}
}
