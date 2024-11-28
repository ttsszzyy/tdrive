package share

import (
	"T-driver/app/user/rpc/pb"
	"T-driver/common/db"
	"T-driver/common/errors"
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"T-driver/app/user/api/internal/svc"
	"T-driver/app/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddShareLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddShareLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddShareLogic {
	return &AddShareLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddShareLogic) AddShare(req *types.AddShareReq) (resp *types.AddShareResp, err error) {
	resp = &types.AddShareResp{}
	//获取用户信息
	userData, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(fmt.Sprintf("%s", l.ctx.Value("language")))
	}
	list, err := l.svcCtx.Rpc.FindAssets(l.ctx, &pb.FindAssetsReq{Ids: req.Ids})
	if err != nil {
		return nil, err
	}
	assetIds, _ := json.Marshal(req.Ids)
	assetType := list.Assets[0].AssetType
	assetName := list.Assets[0].AssetName
	if len(list.Assets) > 1 {
		assetType = 2
		assetName = fmt.Sprintf("%s等%d项", assetName, len(list.Assets))
	}
	//生成一个uuid
	u := uuid.New().String()
	// link := fmt.Sprintf("%s/s/%s", l.svcCtx.Config.FastReward.TgUrl, u)

	//获取资源链接
	link := list.Assets[0].Link
	_, err = l.svcCtx.Rpc.SaveShare(l.ctx, &pb.SaveShareReq{
		Uid:       userData.User.ID,
		AssetIds:  string(assetIds),
		AssetName: assetName,
		AssetSize: list.Assets[0].AssetSize,
		Link:      link,
		AssetType: assetType,
		Uuid:      u,
	})
	if err != nil {
		return nil, err
	}
	resp.Link = fmt.Sprintf("%s/api/v1/user/share/resource/%s", l.svcCtx.Config.BaseUrl, u)
	return resp, nil
}
