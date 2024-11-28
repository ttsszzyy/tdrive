package user

import (
	"T-driver/app/user/api/internal/svc"
	"T-driver/app/user/api/internal/types"
	"T-driver/app/user/model"
	"T-driver/app/user/rpc/pb"
	"T-driver/common/db"
	"T-driver/common/errors"
	"T-driver/common/utils"
	"T-driver/common/utils/rand"
	"context"
	"encoding/json"
	"fmt"
	storage "github.com/utopiosphe/titan-storage-sdk"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"strconv"
	"time"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	resp = &types.LoginResp{}
	//获取用户信息
	userData, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(fmt.Sprintf("%s", l.ctx.Value("language")))
	}
	//添加分布式锁，防止重复
	err = l.svcCtx.Redis.Setex("user_login_"+strconv.FormatInt(userData.User.ID, 10), "1", 2)
	if err != nil {
		return nil, errors.CustomError("Please try again later")
	}
	defer l.svcCtx.Redis.Del("user_login_" + strconv.FormatInt(userData.User.ID, 10))

	obj, err := l.svcCtx.Rpc.FindOneByUid(l.ctx, &pb.UidReq{Uid: userData.User.ID})
	if err != nil {
		return nil, err
	}
	resp.IsNew = false
	username := userData.User.Username
	if username == "" {
		username = userData.User.FirstName + " " + userData.User.LastName
	}
	//添加用户
	if obj.Id == 0 {
		//生成推荐码
	CreateCode:
		code, err := utils.GenerateReferralCode()
		if err != nil {
			return nil, errors.DbError()
		}
		//校验推荐码是否存在，如果存在重新生成
		byCode, err := l.svcCtx.Rpc.FindUserByCode(l.ctx, &pb.UserCodeReq{Code: code})
		if err != nil {
			return nil, errors.DbError()
		}
		if byCode != nil && byCode.Id > 0 {
			goto CreateCode
		}

		//获取用户来源
		var source int64
		pid := int64(0)
		if req.IsShare == 1 {
			source = model.SourceShare
		} else if userData.StartParam != "" { //查询邀请码
			Puser, err := l.svcCtx.Rpc.FindUserByCode(l.ctx, &pb.UserCodeReq{Code: userData.StartParam})
			if err != nil {
				return nil, errors.DbError()

			}
			pid = Puser.Uid
			source = model.SourceInvitation
		}

		obj = &pb.User{
			Uid:           userData.User.ID,
			Name:          username,
			Avatar:        userData.User.PhotoURL,
			Source:        source,
			RecommendCode: code,
			Pid:           pid,
			IsDisable:     model.State_Enable,
			CreatedTime:   time.Now().Unix(),
		}
		_, err = l.svcCtx.Rpc.SaveUser(l.ctx, obj)
		if err != nil {
			return nil, errors.DbError()
		}
	}
	if obj.IsDisable == model.State_Disable {
		return nil, errors.CustomError("Account Is Disabled")
	}
	//校验是否已经邀请用户了，但是还没有领取奖励
	if obj.Pid > 0 {
		go addReward(context.Background(), l.svcCtx, userData.User.ID, obj.Pid, obj.Source, userData.User.IsPremium)
	}
	//更新用户名
	if username != obj.Name || obj.Avatar != userData.User.PhotoURL {
		l.svcCtx.Rpc.SaveUser(l.ctx, &pb.User{
			Id:     obj.Id,
			Name:   username,
			Avatar: userData.User.PhotoURL,
		})
	}

	//查詢是否領取過獎勵
	if obj.IsReceive == 0 {
		resp.IsNew = true
	}
	getCtx, err := l.svcCtx.Redis.GetCtx(l.ctx, fmt.Sprintf(model.GetStorageUser+"%v", userData.User.ID))
	if err != nil && err != redis.Nil {
		return nil, err
	}
	one := &pb.UserTitanToken{}
	if err == redis.Nil || getCtx == "" || getCtx == "{}" {
		login, err := l.svcCtx.Tenant.SSOLogin(l.ctx, storage.SubUserInfo{
			EntryUUID: strconv.FormatInt(userData.User.ID, 10),
			Username:  username,
			Avatar:    userData.User.PhotoURL,
		})
		if err != nil {
			return nil, errors.DbError()
		}
		l.svcCtx.Rpc.UpdateUserTitanToken(l.ctx, &pb.UpdateUserTitanTokenReq{
			Uid:    userData.User.ID,
			Token:  login.Token,
			Expire: login.Exp,
		})
		one.Token = login.Token
		one.Expire = login.Exp
	} else {
		err = json.Unmarshal([]byte(getCtx), one)
		if err != nil {
			return nil, err
		}
	}

	resp.Exp = one.Expire
	resp.TitanToken = one.Token
	return resp, nil
}

// 添加推荐奖励
func addReward(ctx context.Context, svcCtx *svc.ServiceContext, uid, pid, source int64, isPremium bool) {
	//如果是引荐获取10w积分和5w随机积分
	if source == model.SourceInvitation {
		//校验是否邀请过
		invite, err := svcCtx.Rpc.FindUserInvite(ctx, &pb.FindUserInviteReq{Uid: []int64{uid}, Pid: pid})
		if err != nil {
			return
		}
		if len(invite.List) == 0 {
			dict, err := svcCtx.Rpc.FindDictByName(ctx, &pb.FindDictByNameReq{Code: []string{model.RecommendDictCode, model.RecommendVipDictCode}})
			if err != nil || len(dict.DictList) != 2 {
				logx.Error("获取注册奖励配置失败：", err)
				return
			}
			var points int64
			var pointsVip int64
			points, _ = strconv.ParseInt(dict.DictList[0].Value, 10, 64)
			//获得随机积分
			i := rand.RandomInt64(1, dict.DictList[0].BackupValue)
			points += i
			//是否是tg会员
			if isPremium {
				vip, _ := strconv.ParseInt(dict.DictList[1].Value, 10, 64)
				pointsVip += vip
			}
			svcCtx.Rpc.SaveUserReward(ctx, &pb.SaveUserRewardReq{
				Uid:       uid,
				Pid:       pid,
				Points:    points,
				PointsVip: pointsVip,
			})
		}
	}
}
