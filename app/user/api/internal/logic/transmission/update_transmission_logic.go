package transmission

import (
	"T-driver/app/user/model"
	"T-driver/app/user/rpc/pb"
	"T-driver/common/db"
	"T-driver/common/errors"
	"context"
	"fmt"

	"T-driver/app/user/api/internal/svc"
	"T-driver/app/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTransmissionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateTransmissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTransmissionLogic {
	return &UpdateTransmissionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateTransmissionLogic) UpdateTransmission(req *types.UpdateTransmissionReq) (resp *types.Response, err error) {
	resp = &types.Response{}
	//获取用户信息
	_, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(fmt.Sprintf("%s", l.ctx.Value("language")))
	}
	//保存错误到redis
	l.svcCtx.Redis.Set(fmt.Sprintf(model.UploadErrId+"%s", req.Id), req.Err)
	_, err = l.svcCtx.Rpc.UpdateAssetFile(l.ctx, &pb.UpdateAssetFileReq{Id: req.Id, Status: model.AssetStatusError})
	if err != nil {
		return nil, err
	}

	return
}
