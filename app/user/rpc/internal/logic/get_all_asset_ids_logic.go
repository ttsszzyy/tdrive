package logic

import (
	"context"

	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetAllAssetIdsLogic 获取资源目录下的所有子资源id
type GetAllAssetIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// NewGetAllAssetIdsLogic 新建 获取资源目录下的所有子资源id
func NewGetAllAssetIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllAssetIdsLogic {
	return &GetAllAssetIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetAllAssetIds 实现 获取资源目录下的所有子资源id
func (l *GetAllAssetIdsLogic) GetAllAssetIds(in *pb.GetAllAssetIDsReq) (*pb.AllAssetIDsRes, error) {
	resp := new(pb.AllAssetIDsRes)

	aids, err := l.svcCtx.AssetsModel.GetAllIDsByPID(l.ctx, in.Pid, in.Uid)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	resp.Ids = aids

	return resp, nil
}
