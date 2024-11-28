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

type AddReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddReadLogic {
	return &AddReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddReadLogic) AddRead(req *types.DelShareReq) (resp *types.Response, err error) {
	lan := fmt.Sprintf("%s", l.ctx.Value("language"))
	//获取用户信息
	_, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(lan)
	}
	_, err = l.svcCtx.Rpc.UpdateShare(l.ctx, &pb.UpdateShareReq{Id: req.Id, ReadNum: 1})
	if err != nil {
		return nil, err
	}

	return resp, nil
}
