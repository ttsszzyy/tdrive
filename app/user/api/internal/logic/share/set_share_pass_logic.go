package share

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

type SetSharePassLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetSharePassLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetSharePassLogic {
	return &SetSharePassLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetSharePassLogic) SetSharePass(req *types.SetPassReq) (resp *types.Response, err error) {
	//获取用户信息
	_, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(fmt.Sprintf("%s", l.ctx.Value("language")))
	}
	_, err = l.svcCtx.Rpc.UpdateShare(l.ctx, &pb.UpdateShareReq{
		Id:       req.Id,
		Password: req.Pass,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}
