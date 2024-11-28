package migrate

import (
	"T-driver/app/user/model"
	"T-driver/common/locker"
	"T-driver/common/utils"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"time"
)

var (
	Dict_Lock_key     = "lock:dict"     // 字典锁
	Admin_Lock_key    = "lock:admin"    // admin锁
	TaskPool_Lock_key = "lock:taskPool" // 任务池锁
)

func SeedDict(dictModel model.DictModel, r *redis.Redis) {
	locker := locker.NewRedisLocker(r, Dict_Lock_key)
	defer locker.Unlock()
	if err := locker.Lock(); err != nil {
		logx.Error(err)
		return
	}
	// 初始化固定标签
	for _, v := range model.FixedDicts {
		_, err := dictModel.FindOneByCodeDeletedTime(context.Background(), v.Code, 0)
		if err != nil {
			if err == model.ErrNotFound {
				if _, err := dictModel.Insert(context.Background(), &v); err != nil {
					logx.Error(err)
					panic(err)
				}
			} else {
				logx.Errorf("[Check Dict Failed]: %s", err.Error())
				panic(err)
			}
		}
	}
}

func SeedAdmin(adminModel model.AdminModel, r *redis.Redis) {
	locker := locker.NewRedisLocker(r, Admin_Lock_key)
	defer locker.Unlock()
	if err := locker.Lock(); err != nil {
		logx.Error(err)
		return
	}
	// 初始化固定标签
	_, err := adminModel.FindOneByAccountDeletedTime(context.Background(), "admin", 0)
	if err != nil {
		if err == model.ErrNotFound {
			md5 := utils.MD5("adminTdrive")
			if _, err := adminModel.Insert(context.Background(), &model.Admin{
				Account:     "admin",
				Password:    utils.GenPassword(md5, model.PassSalt),
				IsDisable:   2,
				LastTime:    0,
				CreatedTime: time.Now().Unix(),
			}); err != nil {
				logx.Error(err)
				panic(err)
			}
		} else {
			logx.Errorf("[Check admin Failed]: %s", err.Error())
			panic(err)
		}
	}
}

func SeedTaskPool(taskPoolModel model.TaskPoolModel, r *redis.Redis) {
	locker := locker.NewRedisLocker(r, TaskPool_Lock_key)
	defer locker.Unlock()
	if err := locker.Lock(); err != nil {
		logx.Error(err)
		return
	}
	// 初始化固定标签
	_, err := taskPoolModel.FindOneByBuilder(context.Background())
	if err != nil {
		if err == model.ErrNotFound {
			taskPoolModel.Insert(context.Background(), &model.TaskPool{
				Id:          1,
				TaskType:    3,
				TaskName:    "引薦好友，贏取豐厚獎勵",
				TaskNameEn:  "Refer friends and get generous rewards.",
				Integral:    0,
				IsDisable:   2,
				Link:        "",
				Sort:        1,
				CreatedTime: time.Now().Unix(),
			})
			taskPoolModel.Insert(context.Background(), &model.TaskPool{
				Id:          2,
				TaskType:    0,
				TaskName:    "每日簽到，贏取豐厚獎勵",
				TaskNameEn:  "Sign in daily to win generous rewards.",
				Integral:    2000,
				IsDisable:   2,
				Link:        "",
				Sort:        2,
				CreatedTime: time.Now().Unix(),
			})
			taskPoolModel.Insert(context.Background(), &model.TaskPool{
				Id:          3,
				TaskType:    1,
				TaskName:    "鏈接你的錢包",
				TaskNameEn:  "Link your wallet.",
				Integral:    500000,
				IsDisable:   2,
				Link:        "",
				Sort:        3,
				CreatedTime: time.Now().Unix(),
			})
			taskPoolModel.Insert(context.Background(), &model.TaskPool{
				Id:          4,
				TaskType:    1,
				TaskName:    "轉推和分享給好友",
				TaskNameEn:  "Retweet and share with friends.",
				Integral:    100000,
				IsDisable:   2,
				Sort:        4,
				Link:        "https://x.com/tdrive_dao/status/1821886839533531648?t=mipe7CnAVRsEfuxOotEYOw&s=19",
				CreatedTime: time.Now().Unix(),
			})
			taskPoolModel.Insert(context.Background(), &model.TaskPool{
				Id:          5,
				TaskType:    2,
				TaskName:    "加入官方社群",
				TaskNameEn:  "Join the official community.",
				Integral:    150000,
				IsDisable:   2,
				Sort:        5,
				Link:        "https://t.me/+jfcrHclC6TswZDEx",
				CreatedTime: time.Now().Unix(),
			})
			taskPoolModel.Insert(context.Background(), &model.TaskPool{
				Id:          6,
				TaskType:    2,
				TaskName:    "關注官方推特",
				TaskNameEn:  "Follow the official Twitter.",
				Integral:    100000,
				IsDisable:   2,
				Sort:        6,
				Link:        "https://x.com/tdrive_dao?t=UHoVNRJmQjOO-u9F_6eiBg&s=09",
				CreatedTime: time.Now().Unix(),
			})
			taskPoolModel.Insert(context.Background(), &model.TaskPool{
				Id:          7,
				TaskType:    2,
				TaskName:    "關注官方INS",
				TaskNameEn:  "Follow the official Instagram.",
				Integral:    100000,
				IsDisable:   2,
				Sort:        7,
				Link:        "https://www.instagram.com/tdrive_dao?igsh=MTBwdmo0eTRxcTY1Zg==",
				CreatedTime: time.Now().Unix(),
			})
			taskPoolModel.Insert(context.Background(), &model.TaskPool{
				Id:          8,
				TaskType:    4,
				TaskName:    "加入Titan社群",
				TaskNameEn:  "Join the Titan community.",
				Integral:    50000,
				IsDisable:   2,
				Sort:        8,
				Link:        "https://t.me/titannet_dao",
				CreatedTime: time.Now().Unix(),
			})
			taskPoolModel.Insert(context.Background(), &model.TaskPool{
				Id:          9,
				TaskType:    4,
				TaskName:    "關注Titan推特",
				TaskNameEn:  "Follow Titan's Twitter.",
				Integral:    50000,
				IsDisable:   2,
				Sort:        9,
				Link:        "https://x.com/Titannet_dao?t=ZD5qy5F49hy2nbN0jCAPtw&s=09",
				CreatedTime: time.Now().Unix(),
			})
			taskPoolModel.Insert(context.Background(), &model.TaskPool{
				Id:          10,
				TaskType:    4,
				TaskName:    "關注Titan TikTok",
				TaskNameEn:  "Follow Titan's TikTok.",
				Integral:    50000,
				IsDisable:   2,
				Sort:        10,
				Link:        "https://www.tiktok.com/@titannet.dao",
				CreatedTime: time.Now().Unix(),
			})
			taskPoolModel.Insert(context.Background(), &model.TaskPool{
				Id:          11,
				TaskType:    4,
				TaskName:    "關注Titan Youtube",
				TaskNameEn:  "Follow Titan's Youtube",
				Integral:    50000,
				IsDisable:   2,
				Sort:        11,
				Link:        "https://youtube.com/@titan_dao?si=zMWvQv3YtdnaWQud",
				CreatedTime: time.Now().Unix(),
			})
			//清理查询缓存
			taskPoolModel.ClearCache()
		} else {
			logx.Errorf("[Check admin Failed]: %s", err.Error())
			panic(err)
		}
	}
}
