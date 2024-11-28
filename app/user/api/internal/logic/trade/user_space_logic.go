package trade

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

type UserSpaceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserSpaceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserSpaceLogic {
	return &UserSpaceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserSpaceLogic) UserSpace(req *types.UserSpaceReq) (resp *types.UserSpaceResp, err error) {
	//获取用户信息
	userData, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(fmt.Sprintf("%s", l.ctx.Value("language")))
	}
	if req.Size > 100 {
		req.Size = 100
	}
	list, err := l.svcCtx.Rpc.FindUserStorageExchange(l.ctx, &pb.FindUserStorageExchangeReq{Uid: userData.User.ID, Page: req.Page, Size: req.Size})
	if err != nil {
		return nil, err
	}
	spaceItems := make([]*types.SpaceItem, 0)
	for _, v := range list.List {
		spaceItems = append(spaceItems, &types.SpaceItem{
			CreatedTime:     v.CreatedTime,
			Id:              v.Id,
			StorageExchange: v.StorageExchange,
			Uid:             v.Uid,
		})
	}
	return &types.UserSpaceResp{
		Total: list.Total,
		List:  spaceItems,
	}, nil
}
