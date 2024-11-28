package logic

import (
	"T-driver/app/user/model"
	"T-driver/common/errors"
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"T-driver/app/user/rpc/internal/svc"
	"T-driver/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveUserLogic {
	return &SaveUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SaveUserLogic) SaveUser(in *pb.User) (*pb.Response, error) {
	if in.Id > 0 {
		one, err := l.svcCtx.UserModel.FindOne(l.ctx, in.Id)
		if err != nil {
			if err == model.ErrNotFound {
				return nil, errors.SystemError()
			}
			return nil, errors.DbError()
		}
		if in.Name != "" {
			one.Name = in.Name
		}
		if in.Avatar != "" {
			one.Avatar = in.Avatar
		}
		if in.WalletAddress != "" {
			one.WalletAddress = in.WalletAddress
		}
		if in.Distribution > 0 {
			one.Distribution = in.Distribution
		}
		if in.Pid > 0 {
			one.Pid = in.Pid
		}
		if in.IsDisable > 0 {
			one.IsDisable = in.IsDisable
		}
		if in.IsReceive > 0 {
			one.IsReceive = in.IsReceive
		}
		if in.LanguageCode != "" {
			one.LanguageCode = in.LanguageCode
		}
		if in.IsRead > 0 {
			one.IsRead = in.IsRead
		}
		one.UpdatedTime = time.Now().Unix()
		err = l.svcCtx.UserModel.Update(l.ctx, one)
		if err != nil {
			return nil, errors.DbError()
		}
	} else {
		err := l.svcCtx.MysqlConn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
			_, err := l.svcCtx.UserModel.Insert(ctx, &model.User{
				Uid:           in.Uid,
				Name:          in.Name,
				Avatar:        in.Avatar,
				Mail:          in.Mail,
				WalletAddress: in.WalletAddress,
				Distribution:  in.Distribution,
				Pid:           in.Pid,
				IsDisable:     in.IsDisable,
				Source:        in.Source,
				RecommendCode: in.RecommendCode,
				LanguageCode:  "en",
				DeletedTime:   0,
				CreatedTime:   time.Now().Unix(),
				UpdatedTime:   time.Now().Unix(),
			}, session)
			if err != nil {
				return err
			}
			//创建积分和空间
			_, err = l.svcCtx.UserPointsModel.Insert(ctx, &model.UserPoints{
				Uid:         in.Uid,
				Points:      0,
				ReqPoints:   0,
				CreatedTime: time.Now().Unix(),
			}, session)
			if err != nil {
				return err
			}
			_, err = l.svcCtx.UserStorageModel.Insert(ctx, &model.UserStorage{
				Uid:         in.Uid,
				Storage:     0,
				StorageUse:  0,
				CreatedTime: time.Now().Unix(),
			}, session)
			if err != nil {
				return err
			}
			//创建默认文件夹
			_, err = l.svcCtx.AssetsModel.Insert(ctx, &model.Assets{
				Uid:         in.Uid,
				AssetName:   model.MyTondriver,
				AssetType:   1,
				Pid:         1,
				Source:      1,
				IsReport:    2,
				Status:      3,
				IsDefault:   1,
				IsTag:       2,
				UpdatedTime: time.Now().Unix(),
				CreatedTime: time.Now().Unix(),
			}, session)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return nil, errors.DbError()
		}
	}

	return &pb.Response{}, nil
}
