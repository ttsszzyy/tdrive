package parameter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"T-driver/app/admin/api/internal/logic/parameter"
	"T-driver/app/admin/api/internal/svc"
	"T-driver/app/admin/api/internal/types"
	"T-driver/common/errors"
	"T-driver/common/response"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func SendPhotoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SendPhotoReq
		if err := httpx.Parse(r, &req); err != nil {
			logx.Error(err)
			response.Response(w, nil, errors.ParamsError(r.Header.Get("Language")))
			return
		}
		file, header, err := r.FormFile("file")
		if err != nil {
			logx.Error(err)
			response.Response(w, nil, errors.ParamsError(r.Header.Get("Language")))
			return
		}
		defer file.Close()
		req.ChatID = r.FormValue("chat_id")
		req.Caption = strings.Replace(r.FormValue("caption"), "\\n", "\n", -1)
		parseBool, _ := strconv.ParseBool(r.FormValue("is_pin_chat_message"))
		req.IsPinChatMessage = parseBool
		replyMarkup := r.FormValue("reply_markup")
		replyMarkups := make([][]types.Markup, 0)
		err = json.Unmarshal([]byte(replyMarkup), &replyMarkups)
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			return
		}

		req.FileSize = header.Size
		req.ReplyMarkup = replyMarkups
		bytes, err := io.ReadAll(file)
		if err != nil {
			logx.Error(err)
			response.Response(w, nil, errors.CustomError("读取文件失败"))
			return
		}
		req.Photo = bytes
		l := parameter.NewSendPhotoLogic(r.Context(), svcCtx)
		resp, err := l.SendPhoto(&req)
		response.Response(w, resp, err)
	}
}
