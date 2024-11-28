package transmission

import (
	"T-driver/app/user/api/internal/svc"
	"T-driver/app/user/api/internal/types"
	"T-driver/app/user/model"
	"T-driver/app/user/rpc/pb"
	"T-driver/app/user/rpc/user"
	"T-driver/common/db"
	"T-driver/common/errors"
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/zeromicro/go-zero/core/stores/redis"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListLogic) List(req *types.ListReq) (resp *types.ListRes, err error) {
	var (
		wg = new(sync.WaitGroup)
	)

	//获取用户信息
	userData, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(fmt.Sprintf("%s", l.ctx.Value("language")))
	}
	assets, err := l.svcCtx.Rpc.FindAssetFile(l.ctx, &pb.FindAssetFileReq{
		Uid:    userData.User.ID,
		Status: req.Status,
		Page:   req.Page,
		Size:   req.Size,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.ListRes{
		List:  make([]*types.AssetFile, len(assets.List)),
		Total: assets.Total,
	}
	tcli, err := l.svcCtx.GetStorage(userData.User.ID)
	if err != nil {
		logx.Errorf("get titan storage client error:%v", err)
	}
	for i, v := range assets.List {
		wg.Add(1)
		go func(i int, v *pb.AssetFile) {
			defer wg.Done()

			var (
				completion, isDelete int64
				uploadErr, imgBase64 string
				links                []string
			)
			//查询进行中的进度
			if v.Status == model.AssetStatusAfoot {
				//计算完成度
				res, err := l.svcCtx.Redis.Get(fmt.Sprintf(model.UploadId+"%v", v.Id))
				if err != nil && err != redis.Nil {
					logx.Errorf("get upload info error:%v", err)
					return
				}
				completion, _ = strconv.ParseInt(res, 10, 64)
			} else if v.Status == model.AssetStatusError { //查询失败原因
				resErr, err := l.svcCtx.Redis.Get(fmt.Sprintf(model.UploadErrId+"%v", v.Id))
				if err != nil && err != redis.Nil {
					logx.Errorf("get upload's error error:%v", err)
					return
				}
				uploadErr = resErr
			}
			info, err := l.svcCtx.Rpc.FindOneAssets(l.ctx, &user.FindOneAssetsReq{Id: v.AssetId})
			if (err == nil && info.Id == 0) || (err == nil && info.DeletedTime > 0) {
				isDelete = 1
			} else {
				isDelete = 2
			}
			if tcli != nil {
				imgBase64, links = l.svcCtx.GetImgBase64(l.ctx, tcli, v.AssetType, v.Cid)
			}
			resp.List[i] = &types.AssetFile{
				Uid:         v.Uid,
				Cid:         v.Cid,
				AssetName:   v.AssetName,
				AssetSize:   v.AssetSize,
				AssetType:   v.AssetType,
				Source:      v.Source,
				Completion:  completion,
				Id:          v.Id,
				Status:      v.Status,
				CreatedTime: v.CreatedTime,
				IsTag:       v.Tag,
				Url:         v.Link,
				UploadErr:   uploadErr,
				AssetId:     v.AssetId,
				IsDelete:    isDelete,
				ImgBase64:   imgBase64,
				UrlList:     links,
			}
		}(i, v)
	}
	wg.Wait()

	return resp, nil
}
