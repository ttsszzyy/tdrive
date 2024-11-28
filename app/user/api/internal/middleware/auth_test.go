package middleware

import (
	initdata "github.com/telegram-mini-apps/init-data-golang"
	"strings"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	authParts := strings.Split("tma user=%7B%22id%22%3A7488106568%2C%22first_name%22%3A%22s%22%2C%22last_name%22%3A%22monkey%22%2C%22language_code%22%3A%22zh-hans%22%2C%22allows_write_to_pm%22%3Atrue%7D&chat_instance=5743260181132096554&chat_type=sender&auth_date=1723796311&hash=672cac0cd98d7d9ca6593c8bba87f9cca15dea0fd37dc9ad154b61476596ee57", " ")
	authData := authParts[1]
	if err := initdata.Validate(string(authData), "7441175021:AAGYI1bXK-j-NTyYh03JRy2349sn8cxSatE", time.Hour); err != nil {
		t.Log(err)
		return
	}
}
