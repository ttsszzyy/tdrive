package user

import (
	"T-driver/app/user/api/internal/svc"
	"T-driver/app/user/api/internal/types"
	"T-driver/app/user/model"
	"T-driver/app/user/rpc/pb"
	"T-driver/common/db"
	"T-driver/common/errors"
	"context"
	"fmt"
	"sort"
	"strconv"

	"github.com/zeromicro/go-zero/core/mr"
	"github.com/zeromicro/go-zero/core/stores/redis"

	"github.com/zeromicro/go-zero/core/logx"
)

type RackingListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRackingListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RackingListLogic {
	return &RackingListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RackingListLogic) RackingList(req *types.RackingListReq) (resp *types.RackingListRes, err error) {
	resp = &types.RackingListRes{Rackings: make([]*types.Racking, 0)}
	//获取用户信息
	userData, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(fmt.Sprintf("%s", l.ctx.Value("language")))
	}
	if req.Size > 100 {
		req.Size = 100
	}
	// 获取有序集合中的所有用户和积分
	usersWithScores, err := l.svcCtx.Redis.ZrevrangeWithScores(model.UserIntegral, (req.Page-1)*req.Size, (req.Page-1)*req.Size+req.Size-1)
	if err != nil {
		return nil, err
	}
	generateFunc := func(source chan<- redis.Pair) {
		for _, pair := range usersWithScores {
			source <- pair
		}
	}

	mapperFunc := func(pair redis.Pair, writer mr.Writer[*types.Racking], cancel func(error)) {
		id, _ := strconv.ParseInt(pair.Key, 10, 64)
		user, err := l.svcCtx.Rpc.FindOneByUid(l.ctx, &pb.UidReq{Uid: id})
		if err != nil {
			cancel(err)
			return
		}
		ser, err := l.svcCtx.Redis.Zrevrank(model.UserIntegral, pair.Key)
		serial := ser + 1
		if err != nil {
			if err != redis.Nil {
				serial = 0
			} else {
				cancel(errors.UnauthError(fmt.Sprintf("%s", l.ctx.Value("language"))))
				return
			}
		}
		t := &types.Racking{
			Uid:      id,
			Name:     user.Name,
			Serial:   serial,
			Integral: pair.Score,
		}
		writer.Write(t)
	}

	reducerFunc := func(pipe <-chan *types.Racking, writer mr.Writer[*types.RackingListRes], cancel func(error)) {
		resp = &types.RackingListRes{Rackings: make([]*types.Racking, 0)}
		serial, err := l.svcCtx.Redis.Zrevrank(model.UserIntegral, strconv.FormatInt(userData.User.ID, 10))
		if err != nil && err != redis.Nil {
			cancel(errors.UnauthError(fmt.Sprintf("%s", l.ctx.Value("language"))))
			return
		}
		one, err := l.svcCtx.Rpc.FindOneByUid(l.ctx, &pb.UidReq{Uid: userData.User.ID})
		if err != nil {
			cancel(errors.DbError())
			return
		}
		p, err := l.svcCtx.Rpc.FindOneUserPoints(l.ctx, &pb.FindOneUserPointsReq{Uid: userData.User.ID})
		if err != nil {
			cancel(errors.DbError())
			return
		}
		resp.Racking = &types.Racking{
			Uid:      userData.User.ID,
			Name:     one.Name,
			Serial:   serial + 1,
			Integral: p.Points,
		}
		//获取总人数
		count, err := l.svcCtx.Redis.ZcardCtx(l.ctx, model.UserIntegral)
		if err != nil {
			cancel(err)
			return
		}
		resp.Total = int64(count)
		for v := range pipe {
			resp.Rackings = append(resp.Rackings, v)
		}
		//对resp.Rackings安装id大小排序
		sort.Slice(resp.Rackings, func(i, j int) bool {
			return resp.Rackings[i].Serial < resp.Rackings[j].Serial
		})
		writer.Write(resp)
	}

	result, err := mr.MapReduce(generateFunc, mapperFunc, reducerFunc, mr.WithContext(l.ctx), mr.WithWorkers(4))
	if err != nil {
		return nil, err
	}

	return result, nil
}
