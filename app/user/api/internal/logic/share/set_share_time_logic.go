package share

import (
	"T-driver/app/user/rpc/pb"
	"T-driver/common/db"
	"T-driver/common/errors"
	"context"
	"fmt"
	"time"

	"T-driver/app/user/api/internal/svc"
	"T-driver/app/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetShareTimeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetShareTimeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetShareTimeLogic {
	return &SetShareTimeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetShareTimeLogic) SetShareTime(req *types.SetTimeReq) (resp *types.Response, err error) {
	//获取用户信息
	_, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(fmt.Sprintf("%s", l.ctx.Value("language")))
	}
	unix := time.Now().AddDate(0, 0, req.Day).Unix()
	_, err = l.svcCtx.Rpc.UpdateShare(l.ctx, &pb.UpdateShareReq{
		Id:            req.Id,
		EffectiveTime: unix,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}
