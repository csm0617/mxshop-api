package router

import (
	"github.com/gin-gonic/gin"

	"mxshop_api/user-web/api"
)

func InitUserRouter(router *gin.RouterGroup) {
	//创建user路由分组
	UserRouter := router.Group("user")
	//router下的api
	{

		UserRouter.GET("list", api.GetUserList)

	}

}
