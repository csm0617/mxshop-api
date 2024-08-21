package initialize

import (
	"github.com/gin-gonic/gin"

	"mxshop_api/user-web/middlewares"
	router "mxshop_api/user-web/router"
)

func Routers() *gin.Engine {

	Router := gin.Default()
	Router.Use(middlewares.Cors()) //对全局的请求添加跨域处理的中间件
	//定义全局的路由分组（同时也是定义api版本）
	ApiGroup := Router.Group("v1")
	//注册用户服务的路由
	router.InitUserRouter(ApiGroup)
	//注册基础服务的路由
	router.InitBaseRouter(ApiGroup)
	return Router
}
