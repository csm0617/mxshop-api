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
		//在访问用户列表之前
		//1.先加上登录的jwt验证的中间件
		//2.然后加上管理员权限校验的中间件
		UserRouter.GET("list", middlewares.JWTAuth(), middlewares.IsAdminAuth(), api.GetUserList)
		UserRouter.POST("login", api.PassWordLoginForm)

	}

}
