package resource

import (
	"context"

	"T-driver/app/resource/api/internal/svc"
	"T-driver/app/resource/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RelayTitanLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRelayTitanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RelayTitanLogic {
	return &RelayTitanLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RelayTitanLogic) RelayTitan(req *types.RelayTitanReq) (resp string, err error) {

	return
}
