package document

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

type CheckRecycleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckRecycleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckRecycleLogic {
	return &CheckRecycleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckRecycleLogic) CheckRecycle() (resp *types.CheckRecycleRes, err error) {
	lan := fmt.Sprintf("%s", l.ctx.Value("language"))
	//获取用户信息
	userData, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(lan)
	}
	assets, err := l.svcCtx.Rpc.CountAssets(l.ctx, &pb.CountAssetsReq{
		Uid:    userData.User.ID,
		IsDel:  true,
		Status: 3,
	})
	if err != nil {
		logx.Error(err)
		return nil, err
	}
	resp = &types.CheckRecycleRes{}
	if assets.Total > 0 {
		resp.IsRecycle = true
	}

	return resp, nil
}
