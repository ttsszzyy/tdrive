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

func (l *IsTagLogic) IsTag(req *types.IsTagReq) (resp *types.Response, err error) {
	lan := fmt.Sprintf("%s", l.ctx.Value("language"))
	//获取用户信息
	_, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(lan)
	}
	_, err = l.svcCtx.Rpc.UpdateAssetsName(l.ctx, &pb.UpdateAssetsNameReq{
		Id:    req.Id,
		IsTag: req.IsTag,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}
