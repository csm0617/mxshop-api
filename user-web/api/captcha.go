package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha" //base64图形验证码
	"go.uber.org/zap"
	"net/http"
)

var store = base64Captcha.DefaultMemStore //拿到默认存储
/*
*
生成数字验证码
*/
func GetCaptcha(ctx *gin.Context) {
	driverDigit := base64Captcha.NewDriverDigit(80, 240, 4, 0.7, 80) //new一个数字验证码驱动实例
	cp := base64Captcha.NewCaptcha(driverDigit, store)               //new一个验证码实例
	id, b64s, answer, err := cp.Generate()                           //生成验证码
	if err != nil {
		zap.S().Errorf("验证码生成失败%s", err)
		ctx.JSON(http.StatusInternalServerError, map[string]string{
			"msg": "验证码生成失败",
		})
		//失败了，直接返回
		return
	}
	zap.S().Debugf("验证码生成成功:【%s】", answer)
	//将验证码的id和base64格式数据返回
	ctx.JSON(http.StatusOK, gin.H{
		"captchaId": id,
		"picPath":   b64s,
	})

}
