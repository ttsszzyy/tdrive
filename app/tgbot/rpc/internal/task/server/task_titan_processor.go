package server

import (
	"T-driver/app/tgbot/rpc/internal/svc"
	"T-driver/app/user/model"
	"T-driver/app/user/rpc/pb"
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	storage "github.com/utopiosphe/titan-storage-sdk"
	"github.com/zeromicro/go-zero/core/logx"
)

type TaskTitanProcessor struct {
	svcCtx *svc.ServiceContext
	ctx    context.Context
}

func NewTaskTitanProcess(svcCtx *svc.ServiceContext, ctx context.Context) *TaskTitanProcessor {
	return &TaskTitanProcessor{
		svcCtx: svcCtx,
		ctx:    ctx,
	}
}

func (p *TaskTitanProcessor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	logx.Errorf("------payload name:%v---", t.Type())

	callback := &storage.AssetUploadNotifyCallback{}
	err := json.Unmarshal(t.Payload(), callback)
	if err != nil {
		logx.Error("json.Unmarshal:", err)
		return err
	}
	assetFile, err := p.svcCtx.Rpc.FindOneAssetFile(p.ctx, &pb.FindOneAssetFileReq{Id: callback.ExtraID})
	if err != nil {
		logx.Error("FindOneAssetFile:", err)
		return err
	}

	// 如果对应的文件记录已存在，则直接退出
	if assetFile.AssetId > 0 {
		return nil
	}

	//todo 文件夹特殊处理

	logx.Errorf("------callback.ExtraID:%v callback.AssetCID:%v---", callback.ExtraID, callback.AssetCID)

	//文件处理
	_, err = p.svcCtx.Rpc.UpdateAssetFile(p.ctx, &pb.UpdateAssetFileReq{
		Id:        assetFile.Id,
		Cid:       callback.AssetCID,
		Status:    model.AssetStatusEnable,
		AssetSize: callback.AssetSize,
		Link:      callback.AssetDirectUrl,
	})
	if err != nil {
		logx.Error("UpdateAssetFile:", err)
		return err
	}
	return nil
}
