package document

import (
	"T-driver/app/user/model"
	"T-driver/app/user/rpc/pb"
	"T-driver/common/db"
	"T-driver/common/errors"
	"context"
	"fmt"
	"strconv"

	"github.com/zeromicro/go-zero/core/stores/redis"

	"T-driver/app/user/api/internal/svc"
	"T-driver/app/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type InTransitLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInTransitLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InTransitLogic {
	return &InTransitLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InTransitLogic) InTransit(req *types.InTransitReq) (resp *types.InTransitRes, err error) {
	lan := fmt.Sprintf("%s", l.ctx.Value("language"))
	//获取用户信息
	userData, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(lan)
	}
	assets, err := l.svcCtx.Rpc.FindAssets(l.ctx, &pb.FindAssetsReq{
		Uid:        userData.User.ID,
		Status:     req.Status,
		AssetTypes: []int64{2, 3, 4}, //不查询文件夹
		IsAdd:      req.IsAdd,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.InTransitRes{
		List: make([]*types.InTransit, 0, len(assets.Assets)),
	}
	for _, v := range assets.Assets {
		var completion int64
		var uploadErr string
		//查询进行中的进度
		if v.Status == model.AssetStatusAfoot {
			//计算完成度
			res, err := l.svcCtx.Redis.Get(fmt.Sprintf(model.UploadId+"%v", v.Id))
			if err != nil && err != redis.Nil {
				return nil, err
			}
			completion, _ = strconv.ParseInt(res, 10, 64)
		} else if v.Status == model.AssetStatusError {
			//查询失败原因
			resErr, err := l.svcCtx.Redis.Get(fmt.Sprintf(model.UploadErrId+"%v", v.Id))
			if err != nil && err != redis.Nil {
				return nil, err
			}
			uploadErr = resErr
		}
		resp.List = append(resp.List, &types.InTransit{
			Uid:         v.Uid,
			Cid:         v.Cid,
			TransitType: v.TransitType,
			AssetName:   v.AssetName,
			AssetSize:   v.AssetSize,
			AssetType:   v.AssetType,
			Source:      v.Source,
			Completion:  completion,
			Id:          v.Id,
			Status:      v.Status,
			CreatedTime: v.CreatedTime,
			IsTag:       v.IsTag,
			Url:         v.Link,
			UploadErr:   uploadErr,
		})
	}

	return
}
