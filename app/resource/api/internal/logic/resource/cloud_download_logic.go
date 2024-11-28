package resource

import (
	"T-driver/app/user/rpc/pb"
	"T-driver/common/db"
	"T-driver/common/errors"
	"context"
	"fmt"
	"io"

	"T-driver/app/resource/api/internal/svc"
	"T-driver/app/resource/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CloudDownloadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCloudDownloadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CloudDownloadLogic {
	return &CloudDownloadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CloudDownloadLogic) CloudDownload(req *types.CloudDownloadReq) (resp *types.CloudDownloadResp, err error) {
	lan := fmt.Sprintf("%s", l.ctx.Value("language"))
	//获取用户信息
	_, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(lan)
	}
	assets, err := l.svcCtx.Rpc.FindOneAssets(l.ctx, &pb.FindOneAssetsReq{Id: req.Id})
	if err != nil {
		return nil, err
	}
	file, _, err := l.svcCtx.Storage.GetFileWithCid(l.ctx, assets.Cid)
	if err != nil {
		logx.Error("Titan下载文件失败:", err)
		return nil, err
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return &types.CloudDownloadResp{
		AssetName: assets.AssetName,
		Flie:      bytes,
		AssetSize: assets.AssetSize,
	}, nil
}
