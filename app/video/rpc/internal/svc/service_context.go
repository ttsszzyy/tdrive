package svc

import (
	"T-driver/app/video/model"
	"T-driver/app/video/rpc/internal/config"
	"context"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config    config.Config
	MysqlConn sqlx.SqlConn
	//MonCli    *mon.Model
	Redis *redis.Redis
	// --------- models -------------
	VideosModel model.VideosModel
	// --------- mongodb -------------
	VideosLabelModel     model.Videos_labelModel
	VideosLabelUserModel model.Videos_label_userModel
	VideosCommentModel   model.Videos_commentModel
	VideosReplyModel     model.Videos_replyModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.Dns)
	redisCli := redis.MustNewRedis(c.RedisConf)
	/*monCli, err := mon.NewModel(c.MonDb.Uri, c.MonDb.Db, c.MonDb.Collection)
	if err != nil {
		logx.Must(err)
	}*/
	return &ServiceContext{
		Config:    c,
		MysqlConn: conn,
		//MonCli:               monCli,
		Redis:                redisCli,
		VideosModel:          model.NewVideosModel(conn, c.Mysql.DbCache),
		VideosLabelModel:     model.NewVideos_labelModel(c.MonDb.Uri, c.MonDb.Db, "videos_label", c.Mysql.DbCache),
		VideosLabelUserModel: model.NewVideos_label_userModel(c.MonDb.Uri, c.MonDb.Db, "videos_label_user", c.Mysql.DbCache),
		VideosCommentModel:   model.NewVideos_commentModel(c.MonDb.Uri, c.MonDb.Db, "videos_comment", c.Mysql.DbCache),
		VideosReplyModel:     model.NewVideos_replyModel(c.MonDb.Uri, c.MonDb.Db, "videos_reply", c.Mysql.DbCache),
	}
}

func (s *ServiceContext) InitLabel() {
	ctx := context.Background()
	titles := []string{"日常碎片", "生活记录", "美食", "美食分享", "吃货日常", "旅行", "旅游", "打卡景点", "时尚",
		"穿搭", "潮流", "帅哥", "运动", "健身", "足球赛事", "体育竞技", "瑜伽练习", "热门视频", "精彩瞬间", "必看", "娱乐时刻",
		"趣味视频", "轻松一刻", "美好时光", "户外", "美女", "运动达人", "搞笑", "幽默", "风趣", "欢乐无限"}
	for i := 6; i < 7; i++ {
		for _, title := range titles {
			s.VideosLabelModel.Insert(ctx, &model.Videos_label{
				Uid:   7488106568,
				Vid:   int64(i),
				Title: title,
			})
		}
	}
}
