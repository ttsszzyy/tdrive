package errors

var errMsgEn = map[int32]string{
	ErrGetRSAPublickey:      "Obtain RSA public key",
	ErrGetRequestIP:         "Error in obtaining request IP address",
	ErrNotActionDistance:    "Not within the activity address range",
	ErrReceivedActionPoints: "You have already received the activity reward, please do not claim it again",
	ErrSameIPForAction:      "This IP has already received rewards",
}

var errMsgTw = map[int32]string{
	ErrGetRSAPublickey:      "獲取rsa公鑰失敗",
	ErrGetRequestIP:         "獲取請求ip錯誤",
	ErrNotActionDistance:    "不在活動地址範圍",
	ErrReceivedActionPoints: "已經領取過活動獎勵，請勿重複領取",
	ErrSameIPForAction:      "該ip已領取過獎勵",
}
