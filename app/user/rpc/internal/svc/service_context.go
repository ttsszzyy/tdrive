package svc

import (
	"T-driver/app/user/model"
	"T-driver/app/user/rpc/internal/config"
	"T-driver/app/user/rpc/internal/svc/migrate"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ServiceContext struct {
	Config    config.Config
	MysqlConn sqlx.SqlConn
	Redis     *redis.Redis
	MongoCli  *mongo.Client
	// --------- models -------------
	UserModel                model.UserModel
	AssetsModel              model.AssetsModel
	ShareModel               model.ShareModel
	TaskPoolModel            model.TaskPoolModel
	TaskModel                model.TaskModel
	AdminModel               model.AdminModel
	DictModel                model.DictModel
	MessageModel             model.MessageModel
	UserInviteModel          model.UserInviteModel
	UserInviteRewardModel    model.UserInviteRewardModel
	UserPointsModel          model.UserPointsModel
	UserStorageModel         model.UserStorageModel
	UserStorageExchangeModel model.UserStorageExchangeModel
	BotPinMessageModel       model.BotPinMessageModel
	BotCommandModel          model.BotCommandModel
	ActionRecordModel        model.ActionRecordModel
	UserTokenModel           model.UserTokenModel
	UserTokenExchangeModel   model.UserTokenExchangeModel
	// --------- mongodb -------------
	AssetFileModel      model.AssetFileModel
	UserTitanTokenModel model.UserTitanTokenModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.Dns)
	redisCli := redis.MustNewRedis(c.RedisConf)
	connect, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(c.MonDb.Uri))
	if err != nil {
		logx.Must(err)
	}

	svcCtx := &ServiceContext{
		Config:                   c,
		MysqlConn:                conn,
		Redis:                    redisCli,
		MongoCli:                 connect,
		UserModel:                model.NewUserModel(conn, c.Mysql.DbCache),
		AssetsModel:              model.NewAssetsModel(conn, c.Mysql.DbCache),
		ShareModel:               model.NewShareModel(conn, c.Mysql.DbCache),
		TaskPoolModel:            model.NewTaskPoolModel(conn, c.Mysql.DbCache),
		TaskModel:                model.NewTaskModel(conn, c.Mysql.DbCache),
		AdminModel:               model.NewAdminModel(conn, c.Mysql.DbCache),
		DictModel:                model.NewDictModel(conn, c.Mysql.DbCache),
		MessageModel:             model.NewMessageModel(conn, c.Mysql.DbCache),
		UserInviteModel:          model.NewUserInviteModel(conn, c.Mysql.DbCache),
		UserInviteRewardModel:    model.NewUserInviteRewardModel(conn, c.Mysql.DbCache),
		UserPointsModel:          model.NewUserPointsModel(conn, c.Mysql.DbCache),
		UserStorageModel:         model.NewUserStorageModel(conn, c.Mysql.DbCache),
		UserStorageExchangeModel: model.NewUserStorageExchangeModel(conn, c.Mysql.DbCache),
		BotPinMessageModel:       model.NewBotPinMessageModel(conn, c.Mysql.DbCache),
		BotCommandModel:          model.NewBotCommandModel(conn, c.Mysql.DbCache),
		ActionRecordModel:        model.NewActionRecordModel(conn, c.Mysql.DbCache),
		UserTokenModel:           model.NewUserTokenModel(conn, c.Mysql.DbCache),
		UserTokenExchangeModel:   model.NewUserTokenExchangeModel(conn, c.Mysql.DbCache),
		AssetFileModel:           model.NewAssetFileModel(c.MonDb.Uri, c.MonDb.Db, "asset_file", c.Mysql.DbCache),
		UserTitanTokenModel:      model.NewUserTitanTokenModel(c.MonDb.Uri, c.MonDb.Db, "user_titan_token", c.Mysql.DbCache),
	}
	//初始化管理员
	migrate.SeedAdmin(svcCtx.AdminModel, svcCtx.Redis)
	// 初始化字典数据
	migrate.SeedDict(svcCtx.DictModel, svcCtx.Redis)
	// 初始化任务池
	migrate.SeedTaskPool(svcCtx.TaskPoolModel, svcCtx.Redis)
	return svcCtx
}
