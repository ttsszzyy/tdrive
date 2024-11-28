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

type DelUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDelUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelUserLogic {
	return &DelUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DelUserLogic) DelUser(req *types.Request) (resp *types.Response, err error) {
	//获取用户信息
	userData, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(fmt.Sprintf("%s", l.ctx.Value("language")))
	}
	storageCli, err := l.svcCtx.GetStorage(userData.User.ID)
	if err != nil {
		l.Error(err)
		return nil, err
	}
	err = mr.Finish(func() error {
		//清理资源
		assets, err := l.svcCtx.Rpc.FindAssets(l.ctx, &pb.FindAssetsReq{Uid: userData.User.ID})
		if err != nil {
			return errors.FromRpcError(err)
		}
		assetList, err := l.svcCtx.Rpc.FindAssets(l.ctx, &pb.FindAssetsReq{Uid: userData.User.ID, IsDel: true})
		if err != nil {
			return errors.FromRpcError(err)
		}
		assets.Assets = append(assets.Assets, assetList.Assets...)
		ids := make([]int64, 0, len(assets.Assets))
		for _, asset := range assets.Assets {
			if asset.Cid != "" {
				list, err := l.svcCtx.Rpc.CountAssets(l.ctx, &pb.CountAssetsReq{Cid: asset.Cid})
				if err != nil {
					return errors.FromRpcError(err)
				}
				if list.Total <= 1 {
					storageCli.Delete(l.ctx, asset.Cid)
				}
			}
			ids = append(ids, asset.Id)
		}
		_, err = l.svcCtx.Rpc.ClearAssets(l.ctx, &pb.DelAssetsReq{Ids: ids})
		if err != nil {
			return errors.FromRpcError(err)
		}
		return nil
	}, func() error {
		//清理用户
		_, err = l.svcCtx.Rpc.DelUser(l.ctx, &pb.DelUserReq{Uid: userData.User.ID})
		if err != nil {
			return errors.FromRpcError(err)
		}
		return nil
	})
	if err != nil {
		logx.Error(err)
		return nil, err
	}

	return
}
