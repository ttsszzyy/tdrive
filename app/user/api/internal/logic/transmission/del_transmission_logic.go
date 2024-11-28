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

type DelTransmissionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDelTransmissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelTransmissionLogic {
	return &DelTransmissionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DelTransmissionLogic) DelTransmission(req *types.DelTransmissionReq) (resp *types.Response, err error) {
	resp = &types.Response{}
	//获取用户信息
	_, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(fmt.Sprintf("%s", l.ctx.Value("language")))
	}
	_, err = l.svcCtx.Rpc.DelAssetFile(l.ctx, &pb.DelAssetFileReq{Ids: req.Id, IsSource: req.IsSource})
	if err != nil {
		return nil, err
	}
	return
}
