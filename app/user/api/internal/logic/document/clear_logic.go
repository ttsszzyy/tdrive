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
	"sync"

	"github.com/zeromicro/go-zero/core/logx"
)

type ClearLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewClearLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClearLogic {
	return &ClearLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ClearLogic) Clear(req *types.DelReq) (resp *types.Response, err error) {
	var (
		wg = new(sync.WaitGroup)
		mu = new(sync.Mutex)
	)

	lan := fmt.Sprintf("%s", l.ctx.Value("language"))
	//获取用户信息
	userData, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(lan)
	}
	assets := &pb.FindAssetsResp{
		Assets: make([]*pb.Assets, 0),
	}
	//清空titan资源
	if len(req.Ids) == 0 {
		assets, err = l.svcCtx.Rpc.FindAssets(l.ctx, &pb.FindAssetsReq{Uid: userData.User.ID, IsDel: true})
	} else {
		assets, err = l.svcCtx.Rpc.FindAssetsNoDel(l.ctx, &pb.FindAssetsReq{Ids: req.Ids})
	}
	if err != nil {
		logx.Error(err)
		return nil, err
	}
	storageCli, err := l.svcCtx.GetStorage(userData.User.ID)
	if err != nil {
		l.Error(err)
		return nil, err
	}
	ids := make([]int64, 0, len(assets.Assets))
	for _, asset := range assets.Assets {
		wg.Add(1)
		go func(asset *pb.Assets) {
			defer wg.Done()
			//资源删除
			if asset.AssetType != 1 {
				//清理Titan资源
				if asset.Cid != "" {
					count, err := l.svcCtx.Rpc.CountAssets(l.ctx, &pb.CountAssetsReq{Cid: asset.Cid, Status: 3})
					if err != nil {
						logx.Error(err)
						return
					}
					if count.Total == 1 {
						storageCli.Delete(l.ctx, asset.Cid)
					}
				} else {
					//删除上传失败的redis
					l.svcCtx.Redis.Del(fmt.Sprintf(model.UploadErrId+"%v", asset.Id))
				}
			}
			mu.Lock()
			defer mu.Unlock()
			ids = append(ids, asset.Id)
		}(asset)
	}
	wg.Wait()
	_, err = l.svcCtx.Rpc.ClearAssets(l.ctx, &pb.DelAssetsReq{Ids: ids})
	if err != nil {
		return nil, err
	}
	return
}
