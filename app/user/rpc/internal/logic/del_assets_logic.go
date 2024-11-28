package logic

import (
	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type DelAssetsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelAssetsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelAssetsLogic {
	return &DelAssetsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 资源删除
func (l *DelAssetsLogic) DelAssets(in *pb.DelAssetsReq) (*pb.Response, error) {
	// sq := squirrel.Select().Where(squirrel.Eq{"id": in.Ids}).Where("deleted_time = ?", 0)
	// if in.Uid > 0 {
	// 	sq = sq.Where("uid = ?", in.Uid)
	// }
	// list, err := l.svcCtx.AssetsModel.List(l.ctx, sq)
	// if err != nil {
	// 	return nil, err
	// }
	// ids := make([]int64, 0, len(list))
	// for _, v := range list {
	// 	ids = append(ids, v.Id)
	// }
	ids, _, err := l.svcCtx.AssetsModel.GetUserAssetIDs(l.ctx, in.Uid, in.Ids)
	err = l.svcCtx.MysqlConn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		err = l.svcCtx.AssetsModel.Deletes(l.ctx, ids, session)
		if err != nil {
			return err
		}
		//todo 修改mongo文件状态为可重试状态
		/*_, err = l.svcCtx.AssetFileModel.UpdateByAssetIds(l.ctx, in.Ids,5)
		if err != nil {
			logx.Errorf("删除文件失败: %v", err)
			return err
		}*/
		return nil
	})
	if err != nil {
		logx.Error("删除资源失败:", err.Error())
		return nil, err
	}
	return &pb.Response{}, nil
}
