package document

import (
	"T-driver/app/user/api/internal/svc"
	"T-driver/app/user/api/internal/types"
	"T-driver/app/user/rpc/pb"
	"T-driver/common/db"
	"T-driver/common/errors"
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
)

type AssetIdsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAssetIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssetIdsLogic {
	return &AssetIdsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AssetIdsLogic) AssetIds(req *types.AssetIdsReq) (resp *types.DocumentRes, err error) {
	lan := fmt.Sprintf("%s", l.ctx.Value("language"))
	//获取用户信息
	userData, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(lan)
	}
	assets, err := l.svcCtx.Rpc.FindAssets(l.ctx, &pb.FindAssetsReq{
		Ids: req.Ids,
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
