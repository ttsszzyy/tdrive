package trade

import (
	"T-driver/app/user/model"
	"T-driver/app/user/rpc/pb"
	"T-driver/common/db"
	"T-driver/common/errors"
	"context"
	"fmt"
	"strconv"

	"T-driver/app/user/api/internal/svc"
	"T-driver/app/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExchangeStorageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExchangeStorageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExchangeStorageLogic {
	return &ExchangeStorageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExchangeStorageLogic) ExchangeStorage(req *types.ExchangeStorageReq) (resp *types.Response, err error) {
	//获取用户信息
	userData, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(fmt.Sprintf("%s", l.ctx.Value("language")))
	}
	//添加分布式锁，防止重复
	err = l.svcCtx.Redis.Setex("user_exchange_storage_"+strconv.FormatInt(userData.User.ID, 10), "1", 2)
	if err != nil {
		return nil, errors.CustomError("Please try again later")
	}
	defer l.svcCtx.Redis.Del("user_exchange_storage_" + strconv.FormatInt(userData.User.ID, 10))

	u, err := l.svcCtx.Rpc.FindOneUserPoints(l.ctx, &pb.FindOneUserPointsReq{Uid: userData.User.ID})
	if err != nil {
		logx.Errorf("find user error:%v", err)
		return nil, err
	}
	dict, err := l.svcCtx.Rpc.FindDictByName(l.ctx, &pb.FindDictByNameReq{Code: []string{model.SpaceMerchantExchangeDictCode}})
	if err != nil {
		logx.Errorf("find dict error:%v", err)
		return nil, err
	}
	if len(dict.DictList) == 0 {
		return nil, errors.CustomError("未配置空间商人兑换参数")
	}
	value := dict.DictList[0].Value
	storage, _ := strconv.ParseInt(value, 10, 64)
	points := dict.DictList[0].BackupValue
	rema := req.Storage / storage * points
	if req.Storage/storage*points > u.Points {
		return nil, errors.CustomError("积分不足")
	}
	//添加兑换空间记录
	_, err = l.svcCtx.Rpc.SaveUserPointsAndStorage(l.ctx, &pb.SaveUserPointsAndStorageReq{
		Uid:             userData.User.ID,
		Points:          -rema,
		Storage:         req.Storage * storage * 1024 * 1024,
		StorageExchange: req.Storage,
	})
	if err != nil {
		logx.Errorf("save user error:%v", err)
		return nil, err
	}

	return
}
