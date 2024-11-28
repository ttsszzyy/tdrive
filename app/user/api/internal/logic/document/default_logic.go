package document

import (
	"T-driver/app/user/api/internal/svc"
	"T-driver/app/user/api/internal/types"
	"T-driver/app/user/rpc/pb"
	"T-driver/common/db"
	"T-driver/common/errors"
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
)

type DefaultLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDefaultLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DefaultLogic {
	return &DefaultLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DefaultLogic) Default(req *types.Request) (resp *types.DefaultRes, err error) {
	lan := fmt.Sprintf("%s", l.ctx.Value("language"))
	//获取用户信息
	userData, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(lan)
	}
	//查询默认文件夹
	oneAssets, err := l.svcCtx.Rpc.FindOneAssets(l.ctx, &pb.FindOneAssetsReq{Uid: userData.User.ID, IsDefault: 1})
	if err != nil || oneAssets.Id == 0 {
		return nil, errors.ErrorNotFound()
	}
	return &types.DefaultRes{
		AssetName: oneAssets.AssetName,
		Id:        oneAssets.Id,
	}, nil
}
