package router

import (
	"github.com/gin-gonic/gin"

	"mxshop_api/user-web/api"
	"mxshop_api/user-web/middlewares"
)

func InitUserRouter(router *gin.RouterGroup) {
	//创建user路由分组
	UserRouter := router.Group("user")
	//router下的api
	{
		//对访问用户列表加上登录的jwt验证
		UserRouter.GET("list", middlewares.JWTAuth(), api.GetUserList)
		UserRouter.POST("login", api.PassWordLoginForm)

	}

}
