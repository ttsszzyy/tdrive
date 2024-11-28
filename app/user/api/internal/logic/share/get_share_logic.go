package share

import (
	"T-driver/app/user/rpc/pb"
	"T-driver/common/db"
	"T-driver/common/errors"
	"context"
	"fmt"

	"T-driver/app/user/api/internal/svc"
	"T-driver/app/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetShareLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetShareLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetShareLogic {
	return &GetShareLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetShareLogic) GetShare(req *types.GetShareReq) (resp *types.DocumentRes, err error) {
	//获取用户信息
	userData, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(fmt.Sprintf("%s", l.ctx.Value("language")))
	}
	share, err := l.svcCtx.Rpc.FindOneShare(l.ctx, &pb.FindOneShareReq{Id: req.Id})
	if err != nil {
		return
	}
	assets, err := l.svcCtx.Rpc.FindAssets(l.ctx, &pb.FindAssetsReq{
		Ids: share.AssetIds,
		Uid: userData.User.ID,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.DocumentRes{
		Total:      assets.Total,
		AssetItems: make([]*types.AssetItem, 0, len(assets.Assets)),
	}
	for _, asset := range assets.Assets {
		// fmt.Sprintf("%s/api/v1/user/share/resource/%s", l.svcCtx.Config.BaseUrl, v.Uuid)
		resp.AssetItems = append(resp.AssetItems, &types.AssetItem{
			Uid:         asset.Uid,
			Cid:         asset.Cid,
			TransitType: asset.TransitType,
			AssetName:   asset.AssetName,
			AssetSize:   asset.AssetSize,
			AssetType:   asset.AssetType,
			CreatedTime: asset.CreatedTime,
			Pid:         asset.Pid,
			UpdatedTime: asset.UpdatedTime,
			IsTag:       asset.IsTag,
			Source:      asset.Source,
			Url:         asset.Link,
			Id:          asset.Id,
			Status:      asset.Status,
			IsReport:    asset.IsReport,
		})
	}

	return
}
