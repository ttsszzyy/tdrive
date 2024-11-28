package logic

import (
	"T-driver/app/user/model"
	"context"
	"encoding/json"
	"errors"

	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindOneShareLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindOneShareLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindOneShareLogic {
	return &FindOneShareLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查看分享资源
func (l *FindOneShareLogic) FindOneShare(in *pb.FindOneShareReq) (*pb.Share, error) {
	var (
		one      *model.Share
		assetIds = make([]int64, 0)

		err error
	)

	switch {
	case in.Id > 0:
		one, err = l.svcCtx.ShareModel.FindOne(l.ctx, in.Id)
		if err != nil {
			if err == model.ErrNotFound {
				return &pb.Share{}, nil
			}
			return nil, err
		}
	case in.Uuid != "":
		one, err = l.svcCtx.ShareModel.FindOneByUUID(l.ctx, in.Uuid)
		if err != nil {
			if err == model.ErrNotFound {
				return &pb.Share{}, nil
			}
			return nil, err
		}
	default:
		return nil, errors.New("param error")
	}

	json.Unmarshal([]byte(one.AssetIds), &assetIds)
	return &pb.Share{
		Id:            one.Id,
		Uid:           one.Uid,
		AssetIds:      assetIds,
		AssetName:     one.AssetName,
		AssetSize:     one.AssetSize,
		AssetType:     one.AssetType,
		Link:          one.Link,
		EffectiveTime: one.EffectiveTime,
		ReadNum:       one.ReadNum,
		SaveNum:       one.SaveNum,
		CreatedTime:   one.CreatedTime,
		DeletedTime:   one.DeletedTime,
		Password:      one.Password,
	}, nil
}
