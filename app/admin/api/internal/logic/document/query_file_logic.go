package document

import (
	"T-driver/app/user/rpc/pb"
	"T-driver/app/user/rpc/user"
	"context"

	"T-driver/app/admin/api/internal/svc"
	"T-driver/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueryFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryFileLogic {
	return &QueryFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryFileLogic) QueryFile(req *types.QueryFileReq) (resp *types.QueryFileResp, err error) {
	assets, err := l.svcCtx.Rpc.FindAssets(l.ctx, &user.FindAssetsReq{
		AssetName:  req.AssetName,
		Page:       req.Page,
		Size:       req.Size,
		Ids:        []int64{req.Id},
		Status:     []int64{req.Status},
		Order:      1,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
		IsReport:   req.IsReport,
		AssetTypes: req.AssetTypes,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.QueryFileResp{
		Total:      assets.Total,
		AssetItems: make([]*types.AssetItem, 0, len(assets.Assets)),
	}
	for _, asset := range assets.Assets {
		user, err := l.svcCtx.Rpc.FindOneByUid(l.ctx, &pb.UidReq{Uid: asset.Uid})
		if err != nil {
			return nil, err
		}
		//获取资源链接
		url := ""
		/*if asset.Cid != "" {
			shareAssetResult, err := l.svcCtx.Storage.GetURL(l.ctx, asset.Cid)
			if err != nil {
				return nil, err
			}
			if shareAssetResult != nil && len(shareAssetResult.URLs) > 0 {
				url = shareAssetResult.URLs[0]
			}
		}*/

		resp.AssetItems = append(resp.AssetItems, &types.AssetItem{
			Uid:         asset.Uid,
			Cid:         asset.Cid,
			TransitType: asset.TransitType,
			AssetName:   asset.AssetName,
			AssetSize:   asset.AssetSize,
			AssetType:   asset.AssetType,
			CreatedTime: asset.CreatedTime,
			Name:        user.Name,
			UpdatedTime: asset.UpdatedTime,
			IsTag:       asset.IsTag,
			Source:      asset.Source,
			Url:         url,
			Id:          asset.Id,
			Status:      asset.Status,
			IsReport:    asset.IsReport,
		})
	}

	return resp, nil
}
