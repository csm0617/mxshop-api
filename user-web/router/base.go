package router

import (
	"github.com/gin-gonic/gin"

	"mxshop_api/user-web/api"
)

/*
*
基础服务路由
*/
func InitBaseRouter(router *gin.RouterGroup) {
	BaseRouter := router.Group("base")
	{
		BaseRouter.GET("captcha", api.GetCaptcha) //图形验证码接口
		BaseRouter.POST("send_sms", api.SendSms)  //发送短信验证码接口
	}
}
