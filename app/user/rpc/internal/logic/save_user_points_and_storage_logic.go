package logic

import (
	"T-driver/app/user/model"
	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"
	"T-driver/common/errors"
	"context"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"google.golang.org/grpc/status"
	"strconv"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveUserPointsAndStorageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveUserPointsAndStorageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveUserPointsAndStorageLogic {
	return &SaveUserPointsAndStorageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 保存用户积分和空间
func (l *SaveUserPointsAndStorageLogic) SaveUserPointsAndStorage(in *pb.SaveUserPointsAndStorageReq) (*pb.Response, error) {
	err := l.svcCtx.MysqlConn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		if in.Points != 0 || in.ReqPoints != 0 {
			one, err := l.svcCtx.UserPointsModel.FindOneByUidDeletedTime(l.ctx, in.Uid, 0)
			if err != nil {
				if err == model.ErrNotFound {
					one = &model.UserPoints{
						Uid:         in.Uid,
						CreatedTime: time.Now().Unix(),
					}
				} else {
					return err
				}
			}
			if in.Points < 0 && one.Points+in.Points <= 0 {
				return status.Error(errors.ErrUserPointsNotEnough, "积分不足")
			}
			if in.Points != 0 {
				one.Points += in.Points
			}
			if in.ReqPoints != 0 {
				one.ReqPoints += in.ReqPoints
			}
			if one.Id == 0 {
				_, err = l.svcCtx.UserPointsModel.Insert(l.ctx, one, session)
				if err != nil {
					return err
				}
			} else {
				err = l.svcCtx.UserPointsModel.Update(l.ctx, one, session)
				if err != nil {
					return err
				}
			}

			//更新用户积分
			l.svcCtx.Redis.ZaddCtx(ctx, model.UserIntegral, one.Points, strconv.FormatInt(in.Uid, 10))
		}
		if in.Storage != 0 || in.StorageUse != 0 {
			one, err := l.svcCtx.UserStorageModel.FindOneByUidDeletedTime(l.ctx, in.Uid, 0)
			if err != nil {
				if err == model.ErrNotFound {
					one = &model.UserStorage{
						Uid:         in.Uid,
						CreatedTime: time.Now().Unix(),
					}
				} else {
					return err
				}
			}
			if in.Storage != 0 {
				one.Storage += in.Storage
			}
			if one.Storage < 0 {
				one.Storage = 0
			}
			if in.ReqPoints != 0 {
				one.StorageUse += in.StorageUse
			}
			if one.StorageUse < 0 {
				one.StorageUse = 0
			}
			if one.Id == 0 {
				_, err = l.svcCtx.UserStorageModel.Insert(l.ctx, one, session)
				if err != nil {
					return err
				}
			} else {
				err = l.svcCtx.UserStorageModel.Update(l.ctx, one, session)
				if err != nil {
					return err
				}
			}
		}
		if in.StorageExchange > 0 {
			_, err := l.svcCtx.UserStorageExchangeModel.Insert(l.ctx, &model.UserStorageExchange{
				Uid:             in.Uid,
				ExchangeStorage: in.StorageExchange,
				CreatedTime:     time.Now().Unix(),
			}, session)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &pb.Response{}, nil
}
