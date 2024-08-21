package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"golang.org/x/exp/rand"

	"mxshop_api/user-web/forms"
	"mxshop_api/user-web/global"
)

/*
*
生成width长度的验证码
*/
func GenerateSmsCode(width int) string {

	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(uint64(time.Now().UnixNano()))
	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}

func SendSms(ctx *gin.Context) {
	sendSmsForm := forms.SendSmsForm{}
	if err := ctx.ShouldBindJSON(&sendSmsForm); err != nil {
		HandleValidatorError(err, ctx)
		return
	}
	smsCode := GenerateSmsCode(6)
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	})
	rdb.Set(context.Background(), sendSmsForm.Mobile, smsCode, time.Duration(global.ServerConfig.RedisInfo.Expire)*time.Second)
	ctx.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("验证码发送成功:【%s】", smsCode),
	})

}
