package bot

import (
	"T-driver/app/tgbot/api/internal/svc"
	"T-driver/app/user/model"
	"T-driver/common/errors"
	"T-driver/common/response"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
)

func TitanCallbackHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		callback, err := svcCtx.Tenant.ValidateUploadCallback(ctx, svcCtx.Config.TiTan.APISecret, r)
		if err != nil {
			logx.WithContext(ctx).Error(err)
			response.ResponseBlob(w, "error", errors.NewErrCodeMsg(errors.ErrInvalidCallback, "ValidateUploadCallback error: "+err.Error()))
			return
		}

		// 发往MQ
		payload, err := json.Marshal(callback)
		if err != nil {
			logx.WithContext(ctx).Error(err)
		}
		t := asynq.NewTask(model.TitanProcess, payload)
		//防止重复消费和重试机制
		_, err = svcCtx.AsynqClient.EnqueueContext(
			ctx, t,
			asynq.TaskID(fmt.Sprintf("%s:%s", model.TitanProcess, callback.ExtraID)),
			asynq.Queue(model.QueueLow),
		)
		if err != nil {
			logx.WithContext(ctx).Error(err)
		}
		response.ResponseRaw(w, "success")
	}
}
