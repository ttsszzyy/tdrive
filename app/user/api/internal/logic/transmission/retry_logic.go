package transmission

import (
	"T-driver/app/user/api/internal/svc"
	"T-driver/app/user/api/internal/types"
	"T-driver/app/user/model"
	"T-driver/app/user/rpc/pb"
	"T-driver/common/db"
	"T-driver/common/errors"
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
)

type RetryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRetryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RetryLogic {
	return &RetryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RetryLogic) Retry(req *types.RetryReq) (resp *types.Response, err error) {
	//获取用户信息
	_, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(fmt.Sprintf("%s", l.ctx.Value("language")))
	}
	one, err := l.svcCtx.Rpc.FindOneAssetFile(l.ctx, &pb.FindOneAssetFileReq{Id: req.Id})
	if err != nil {
		return nil, errors.FromRpcError(err)
	}

	assets, err := l.svcCtx.Rpc.FindOneAssets(l.ctx, &pb.FindOneAssetsReq{Id: one.AssetId})
	if err != nil {
		return nil, errors.FromRpcError(err)
	}
	if assets.Id == 0 {
		_, err = l.svcCtx.Rpc.UpdateAssetFile(l.ctx, &pb.UpdateAssetFileReq{Id: req.Id, Status: model.AssetStatusError})
		if err != nil {
			return nil, errors.FromRpcError(err)
		}
		return nil, errors.CustomError("asset not found")
	}
	_, err = l.svcCtx.Rpc.UpdateAssetFile(l.ctx, &pb.UpdateAssetFileReq{Id: req.Id, Status: model.AssetStatusEnable})
	if err != nil {
		return nil, errors.FromRpcError(err)
	}

	return &types.Response{}, nil
}
