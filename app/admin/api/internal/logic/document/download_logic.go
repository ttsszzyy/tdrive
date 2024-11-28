package document

import (
	"T-driver/app/user/rpc/pb"
	"context"
	"io"

	"T-driver/app/admin/api/internal/svc"
	"T-driver/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DownloadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDownloadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DownloadLogic {
	return &DownloadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DownloadLogic) Download(req *types.DownloadReq) (resp *types.DownloadResp, err error) {
	assets, err := l.svcCtx.Rpc.FindOneAssets(l.ctx, &pb.FindOneAssetsReq{Id: req.Id})
	if err != nil {
		return nil, err
	}
	file, _, err := l.svcCtx.Storage.GetFileWithCid(l.ctx, assets.Cid)

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return &types.DownloadResp{
		Flie: bytes,
	}, nil
}
