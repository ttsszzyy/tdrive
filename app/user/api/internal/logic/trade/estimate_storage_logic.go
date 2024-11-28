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

type EstimateStorageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEstimateStorageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EstimateStorageLogic {
	return &EstimateStorageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EstimateStorageLogic) EstimateStorage(req *types.Request) (resp *types.EstimateStorageResp, err error) {
	//获取用户信息
	_, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(fmt.Sprintf("%s", l.ctx.Value("language")))
	}
	dict, err := l.svcCtx.Rpc.FindDictByName(l.ctx, &pb.FindDictByNameReq{Code: []string{model.SpaceMerchantExchangeDictCode}})
	if err != nil {
		return nil, err
	}
	if len(dict.DictList) == 0 {
		return nil, errors.CustomError("未配置空间商人兑换参数")
	}
	value := dict.DictList[0].Value
	storage, _ := strconv.ParseInt(value, 10, 64)
	points := dict.DictList[0].BackupValue

	return &types.EstimateStorageResp{
		Storage: storage,
		Points:  points,
	}, nil
}
