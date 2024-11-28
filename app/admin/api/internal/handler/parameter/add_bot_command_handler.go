package parameter

import (
	"T-driver/app/admin/api/internal/logic/parameter"
	"T-driver/app/admin/api/internal/svc"
	"T-driver/app/admin/api/internal/types"
	"T-driver/common/errors"
	"T-driver/common/response"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func AddBotCommandHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddBotCommandReq
		if err := httpx.Parse(r, &req); err != nil {
			logx.Error(err)
			response.Response(w, nil, errors.ParamsError(r.Header.Get("Language")))
			return
		}
		file, header, err := r.FormFile("file")
		/*if err != nil {
			logx.Error(err)
			response.Response(w, nil, errors.ParamsError(r.Header.Get("Language")))
			return
		}*/
		if err == nil {
			defer file.Close()
			if header.Size > 10*1024*1024*1024 {
				response.Response(w, nil, errors.CustomError(fmt.Sprintf("文件大小不能超过%d MB", 10)))
				return
			}
			bytes, err := io.ReadAll(file)
			if err != nil {
				logx.Error(err)
				response.Response(w, nil, errors.CustomError("读取文件失败"))
				return
			}
			req.Photo = bytes
		}

		req.BotCommand = r.FormValue("bot_command")
		req.LanguageCode = r.FormValue("language_code")
		req.Description = r.FormValue("description")
		req.Text = strings.Replace(r.FormValue("text"), "\\n", "\n", -1)
		req.SendType, _ = strconv.ParseInt(r.FormValue("send_type"), 10, 64)
		replyMarkup := r.FormValue("reply_markup")
		replyMarkups := make([][]types.Markup, 0)
		err = json.Unmarshal([]byte(replyMarkup), &replyMarkups)
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			return
		}
		req.ReplyMarkup = replyMarkups

		l := parameter.NewAddBotCommandLogic(r.Context(), svcCtx)
		resp, err := l.AddBotCommand(&req)
		response.Response(w, resp, err)
	}
}
