package initialize

import (
	"github.com/gin-gonic/gin"

	userRouter "mxshop_api/user-web/router"
)

func Routers() *gin.Engine {

	Router := gin.Default()
	//定义全局的路由分组（同时也是定义api版本）
	ApiGroup := Router.Group("v1")
	userRouter.InitUserRouter(ApiGroup)
	return Router
}
