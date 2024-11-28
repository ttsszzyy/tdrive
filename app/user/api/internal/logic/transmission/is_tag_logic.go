package transmission

import (
	"T-driver/app/user/rpc/pb"
	"T-driver/common/db"
	"T-driver/common/errors"
	"context"
	"fmt"

	"T-driver/app/user/api/internal/svc"
	"T-driver/app/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type IsTagLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIsTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsTagLogic {
	return &IsTagLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IsTagLogic) IsTag(req *types.IsTagAssetFileReq) (resp *types.Response, err error) {
	//获取用户信息
	_, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(fmt.Sprintf("%s", l.ctx.Value("language")))
	}
	_, err = l.svcCtx.Rpc.UpdateAssetFile(l.ctx, &pb.UpdateAssetFileReq{
		Id:      req.Id,
		AssetId: req.AssetId,
		Tag:     req.IsTag,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}
