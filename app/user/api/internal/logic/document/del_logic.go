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

type DelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelLogic {
	return &DelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DelLogic) Del(req *types.DelReq) (resp *types.Response, err error) {
	lan := fmt.Sprintf("%s", l.ctx.Value("language"))
	resp = &types.Response{}
	//获取用户信息
	userData, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(lan)
	}
	// 获取所有的包含父级文件id的文件id
	rsp, err := l.svcCtx.Rpc.GetAllAssetIds(l.ctx, &pb.GetAllAssetIDsReq{Pid: req.Ids, Uid: userData.User.ID})
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.Rpc.DelAssets(l.ctx, &pb.DelAssetsReq{Ids: rsp.Ids, Uid: userData.User.ID})
	if err != nil {
		return nil, err
	}
	return resp, nil
}
