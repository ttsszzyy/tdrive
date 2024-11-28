package document

import (
	"T-driver/app/user/rpc/pb"
	"T-driver/common/db"
	"T-driver/common/errors"
	"context"
	"fmt"
	"sync"
	"time"

	"T-driver/app/user/api/internal/svc"
	"T-driver/app/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DocumentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDocumentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DocumentLogic {
	return &DocumentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DocumentLogic) Document(req *types.DocumentReq) (resp *types.DocumentRes, err error) {
	var (
		wg = new(sync.WaitGroup)
	)

	lan := fmt.Sprintf("%s", l.ctx.Value("language"))
	//获取用户信息
	userData, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(lan)
	}
	if req.Size > 100 {
		req.Size = 100
	}
	assets, err := l.svcCtx.Rpc.FindAssets(l.ctx, &pb.FindAssetsReq{
		Pid:        req.Pid,
		AssetName:  req.AssetName,
		Page:       req.Page,
		Size:       req.Size,
		Status:     []int64{3},
		Uid:        userData.User.ID,
		IsTag:      req.IsTag,
		IsDel:      req.IsDel,
		AssetTypes: req.AssetTypes,
		Order:      req.Order,
		Sort:       req.Sort,
		IsAdd:      req.IsAdd,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.DocumentRes{
		Total:      assets.Total,
		AssetItems: make([]*types.AssetItem, len(assets.Assets)),
	}
	tcli, err := l.svcCtx.GetStorage(userData.User.ID)
	if err != nil {
		logx.Errorf("get titan storage client error:%v", err)
	}
	for i, asset := range assets.Assets {
		wg.Add(1)
		go func(i int, asset *pb.Assets) {
			defer wg.Done()

			var (
				idDel             int64
				days              int
				errTip, imgBase64 string
				links             []string
			)

			if req.IsDel {
				unix := time.Unix(asset.DeletedTime, 0)
				duration := time.Now().Sub(unix)
				days = l.svcCtx.Config.FastReward.RecoveryDate - int(duration.Hours()/24)
			}
			/*if asset.AssetType != 1 && asset.Link == "" {
				shareAssetResult, err := l.svcCtx.Storage.GetURL(l.ctx, asset.Cid)
				if err != nil {
					logx.Error("GetURL:", err)
					errTip = err.Error()
				}
				if shareAssetResult != nil && len(shareAssetResult.URLs) > 0 {
					asset.Link = shareAssetResult.URLs[0]
					go func(svcCtx *svc.ServiceContext, id int64, link string) {
						svcCtx.Rpc.UpdateAssetsName(context.TODO(), &pb.UpdateAssetsNameReq{Id: id, Link: link})
					}(l.svcCtx, asset.Id, asset.Link)
				}
			}*/

			if asset.DeletedTime > 0 {
				idDel = 1
			}

			if tcli != nil {
				imgBase64, links = l.svcCtx.GetImgBase64(l.ctx, tcli, asset.AssetType, asset.Cid)
			}

			resp.AssetItems[i] = &types.AssetItem{
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
				AssetID:     asset.Id,
				Status:      asset.Status,
				IsReport:    asset.IsReport,
				DeletedDay:  days,
				ErrTip:      errTip,
				IsDefault:   asset.IsDefault,
				IsDelete:    idDel,
				ImgBase64:   imgBase64,
				UrlList:     links,
			}
		}(i, asset)
	}
	wg.Wait()

	/*generate := func(asset chan<- *pb.Assets) {
		for _, v := range assets.Assets {
			asset <- v
		}
	}
	mapper := func(asset *pb.Assets, writer mr.Writer[*types.AssetItem], cancel func(error)) {
		days := 0
		errTip := ""
		if req.IsDel {
			unix := time.Unix(asset.DeletedTime, 0)
			duration := time.Now().Sub(unix)
			days = l.svcCtx.Config.FastReward.RecoveryDate - int(duration.Hours()/24)
		}
		if asset.AssetType != 1 && asset.Link == "" {
			shareAssetResult, err := l.svcCtx.Storage.GetURL(l.ctx, asset.Cid)
			if err != nil {
				logx.Error("GetURL:", err)
				errTip = err.Error()
			}
			if shareAssetResult != nil && len(shareAssetResult.URLs) > 0 {
				asset.Link = shareAssetResult.URLs[0]
				go func(svcCtx *svc.ServiceContext, id int64, link string) {
					svcCtx.Rpc.UpdateAssetsName(context.TODO(), &pb.UpdateAssetsNameReq{Id: id, Link: link})
				}(l.svcCtx, asset.Id, asset.Link)
			}
		}

		writer.Write(&types.AssetItem{
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
			DeletedDay:  days,
			ErrTip:      errTip,
			IsDefault:   asset.IsDefault,
		})
	}
	reducer := func(list <-chan *types.AssetItem, writer mr.Writer[*types.DocumentRes], cancel func(error)) {
		resp = &types.DocumentRes{
			Total:      assets.Total,
			AssetItems: make([]*types.AssetItem, 0, len(assets.Assets)),
		}
		for assetItem := range list {
			resp.AssetItems = append(resp.AssetItems, assetItem)
		}
		writer.Write(resp)
	}
	reduce, err := mr.MapReduce(generate, mapper, reducer)
	if err != nil {
		return nil, err
	}*/

	return resp, nil
}
