package user

import (
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

	"T-driver/app/user/api/internal/svc"
	"T-driver/app/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RackingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRackingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RackingLogic {
	return &RackingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RackingLogic) Racking(req *types.Request) (resp *types.RackingListRes, err error) {
	//获取用户信息
	_, ok := db.CtxInitData(l.ctx)
	if !ok {
		return nil, errors.UnauthError(fmt.Sprintf("%s", l.ctx.Value("language")))
	}
	// 获取有序集合中的所有用户和积分
	usersWithScores, err := l.svcCtx.Redis.ZrevrangeWithScores(model.UserIntegral, 0, 2)
	if err != nil && err != redis.Nil {
		return nil, err
	}
	generateFunc := func(source chan<- redis.Pair) {
		for _, score := range usersWithScores {
			source <- score
		}
	}
	mapperFunc := func(item redis.Pair, writer mr.Writer[*types.Racking], cancel func(error)) {
		id, _ := strconv.ParseInt(item.Key, 10, 64)
		user, err := l.svcCtx.Rpc.FindOneByUid(l.ctx, &pb.UidReq{Uid: id})
		if err != nil {
			cancel(errors.DbError())
			return
		}
		serial, err := l.svcCtx.Redis.Zrevrank(model.UserIntegral, item.Key)
		if err != nil && err != redis.Nil {
			cancel(errors.UnauthError(fmt.Sprintf("%s", l.ctx.Value("language"))))
			return
		}
		writer.Write(&types.Racking{
			Uid:      id,
			Name:     user.Name,
			Serial:   serial + 1,
			Integral: item.Score,
		})
	}
	reducerFunc := func(pipe <-chan *types.Racking, writer mr.Writer[*types.RackingListRes], cancel func(error)) {
		resp = &types.RackingListRes{Rackings: make([]*types.Racking, 0)}
		for racking := range pipe {
			resp.Rackings = append(resp.Rackings, racking)
		}
		sort.Slice(resp.Rackings, func(i, j int) bool {
			return resp.Rackings[i].Serial < resp.Rackings[j].Serial
		})
		writer.Write(resp)
	}
	reduce, err := mr.MapReduce(generateFunc, mapperFunc, reducerFunc, mr.WithContext(l.ctx))
	if err != nil {
		return nil, err
	}

	return reduce, nil
}
