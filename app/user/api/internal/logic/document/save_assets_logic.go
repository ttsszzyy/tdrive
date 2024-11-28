package document

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

type SaveAssetsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveAssetsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveAssetsLogic {
	return &SaveAssetsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveAssetsLogic) SaveAssets(req *types.SaveAssetsReq) (resp *types.SaveAssetsResp, err error) {
	//获取用户信息
	userData, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(fmt.Sprintf("%s", l.ctx.Value("language")))
	}
	assets, err := l.svcCtx.Rpc.FindAssets(l.ctx, &pb.FindAssetsReq{
		Ids: req.Ids,
		Uid: userData.User.ID,
	})
	if err != nil {
		return nil, err
	}
	ids := make([]int64, 0)
	for _, v := range assets.Assets {
		if v.AssetType == 1 {
			err := l.getAssetsIds(ids, v.Id, userData.User.ID)
			if err != nil {
				return nil, err
			}
		}
		ids = append(ids, v.Id)
	}

	pid, err := l.svcCtx.AssetsFolder(model.PackFromShared, userData.User.ID)
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.Rpc.UpdateAssetsCopy(l.ctx, &pb.UpdateAssetsCopyReq{
		Ids: ids,
		Pid: pid,
	})
	if err != nil {
		return nil, err
	}

	return &types.SaveAssetsResp{
		Id:        pid,
		AssetName: model.PackFromShared,
	}, nil
}

func (l *SaveAssetsLogic) getAssetsIds(ids []int64, assetId int64, uid int64) error {
	list, err := l.svcCtx.Rpc.FindAssets(l.ctx, &pb.FindAssetsReq{
		Pid: assetId,
		Uid: uid,
	})
	if err != nil {
		return err
	}
	for _, v := range list.Assets {
		if v.AssetType == 1 {
			l.getAssetsIds(ids, v.Id, uid)
		}
		ids = append(ids, v.Id)
	}
	return nil
}
