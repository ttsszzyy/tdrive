package logic

import (
	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DelAssetFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelAssetFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelAssetFileLogic {
	return &DelAssetFileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 资源删除
func (l *DelAssetFileLogic) DelAssetFile(in *pb.DelAssetFileReq) (*pb.Response, error) {
	oids := make([]primitive.ObjectID, 0, len(in.Ids))
	for _, v := range in.Ids {
		oid, err := primitive.ObjectIDFromHex(v)
		if err != nil {
			return nil, err
		}
		oids = append(oids, oid)
	}
	list, err := l.svcCtx.AssetFileModel.List(l.ctx, bson.M{"_id": bson.M{"$in": oids}})
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(list))
	assetIds := make([]int64, 0, len(list))
	for _, v := range list {
		ids = append(ids, v.ID.Hex())
		assetIds = append(assetIds, v.AssetId)
	}
	if in.IsSource {
		err = l.svcCtx.MysqlConn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
			err = l.svcCtx.AssetsModel.Deletes(ctx, assetIds, session)
			if err != nil {
				return err
			}
			_, err = l.svcCtx.AssetFileModel.Deletes(l.ctx, ids)
			return err
		})
		if err != nil {
			logx.Error("删除资源失败:", err.Error())
			return nil, err
		}
	} else {
		_, err := l.svcCtx.AssetFileModel.Deletes(l.ctx, ids)
		if err != nil {
			return nil, err
		}
	}

	return &pb.Response{}, nil
}
