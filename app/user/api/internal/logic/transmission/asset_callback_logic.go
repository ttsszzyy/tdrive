package transmission

import (
	"T-driver/app/user/model"
	"T-driver/app/user/rpc/pb"
	"context"

	"T-driver/app/user/api/internal/svc"
	"T-driver/app/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AssetCallbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAssetCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssetCallbackLogic {
	return &AssetCallbackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AssetCallbackLogic) AssetCallback(req *types.AssetCallbackReq) (resp *types.Response, err error) {
	assetFile, err := l.svcCtx.Rpc.FindOneAssetFile(l.ctx, &pb.FindOneAssetFileReq{Id: req.Id})
	if err != nil {
		logx.Error("AssetCallbackLogic.FindOneAssetFile:", err)
		return nil, err
	}
	//todo 文件夹特殊处理

	//文件处理
	_, err = l.svcCtx.Rpc.UpdateAssetFile(l.ctx, &pb.UpdateAssetFileReq{
		Id:        assetFile.Id,
		Cid:       req.Cid,
		Status:    model.AssetStatusEnable,
		AssetSize: assetFile.AssetSize,
		Link:      req.Link,
	})
	if err != nil {
		logx.Error("UpdateAssetFile:", err)
		return nil, err
	}
	return resp, nil
}
