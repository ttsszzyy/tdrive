package bot

import (
	"T-driver/app/tgbot/api/internal/svc"
	"T-driver/app/user/model"
	"T-driver/common/lib/json"
	"context"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"strconv"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// 登录
func BotWebhookHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	// return func(w http.ResponseWriter, r *http.Request) {
	// 	l := bot.NewBotWebhookLogic(r.Context(), svcCtx)
	// 	err := l.BotWebhook()
	// 	if err != nil {
	// 		httpx.ErrorCtx(r.Context(), w, err)
	// 	} else {
	// 		httpx.Ok(w)
	// 	}

	// 	svcCtx.TgBot.WebhookHandler()
	// }
	return svcCtx.TgBot.WebhookHandler()
}

func RegisterBotProcessHandler(svcCtx *svc.ServiceContext) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		// 发往MQ
		payload, err := json.Marshal(update)
		if err != nil {
			logx.WithContext(ctx).Error(err)
		}
		t := asynq.NewTask(model.BotProcess, payload)
		//防止重复消费和重试机制
		//, asynq.TaskID(strconv.FormatInt(update.Message.Chat.ID, 10)), asynq.Retention(3)
		var taskID string
		switch {
		case update.Message != nil:
			taskID = strconv.Itoa(update.Message.ID)
		case update.CallbackQuery != nil:
			taskID = update.CallbackQuery.ID
		}

		_, err = svcCtx.AsynqClient.EnqueueContext(
			ctx, t,
			asynq.TaskID(taskID),
			asynq.Queue(model.QueueLow),
		)
		if err != nil {
			logx.WithContext(ctx).Error(err)
		}
	}
}
