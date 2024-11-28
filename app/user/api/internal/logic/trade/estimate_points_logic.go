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

type EstimatePointsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEstimatePointsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EstimatePointsLogic {
	return &EstimatePointsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EstimatePointsLogic) EstimatePoints(req *types.EstimatePointsReq) (resp *types.EstimatePointsResp, err error) {
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
	if req.Storage < 0 {
		return nil, errors.CustomError("存储空间不能小于0")
	}
	value := dict.DictList[0].Value
	storage, _ := strconv.ParseInt(value, 10, 64)
	points := dict.DictList[0].BackupValue

	i := req.Storage / storage * points
	return &types.EstimatePointsResp{
		Points: i,
	}, nil
}
