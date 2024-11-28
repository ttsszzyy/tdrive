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

type MoveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MoveLogic {
	return &MoveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MoveLogic) Move(req *types.MoveReq) (resp *types.Response, err error) {
	lan := fmt.Sprintf("%s", l.ctx.Value("language"))
	//获取用户信息
	_, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(lan)
	}
	_, err = l.svcCtx.Rpc.UpdateAssetsMove(l.ctx, &pb.UpdateAssetsMoveReq{
		Ids: req.Ids,
		Pid: req.Pid,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}
