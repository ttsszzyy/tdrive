package document

import (
	"T-driver/app/admin/api/internal/svc"
	"T-driver/app/admin/api/internal/types"
	"T-driver/app/user/rpc/pb"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDelFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelFileLogic {
	return &DelFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DelFileLogic) DelFile(req *types.DelFileReq) (resp *types.Response, err error) {
	_, err = l.svcCtx.Rpc.DelAssets(l.ctx, &pb.DelAssetsReq{Ids: []int64{req.Id}})
	if err != nil {
		return nil, err
	}
	//清空titan资源
	assets, err := l.svcCtx.Rpc.FindAssets(l.ctx, &pb.FindAssetsReq{Ids: []int64{req.Id}})
	if err != nil {
		return nil, err
	}
	for _, asset := range assets.Assets {
		//清理Titan资源
		list, err := l.svcCtx.Rpc.FindAssets(l.ctx, &pb.FindAssetsReq{Cid: asset.Cid})
		if err != nil {
			return nil, err
		}
		if len(list.Assets) == 1 {
			l.svcCtx.Storage.Delete(l.ctx, asset.Cid)
		}
	}
	_, err = l.svcCtx.Rpc.ClearAssets(l.ctx, &pb.DelAssetsReq{Ids: []int64{req.Id}})
	if err != nil {
		return nil, err
	}

	return
}
