package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

/*
*
跨域请求处理中间件
*/
func Cors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		method := ctx.Request.Method
		//允许跨域的源
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token,Authorization,Token,x-token")
		//允许跨域的方法
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS,DELETE, PATCH")
		ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin,Access-Control-Allow-Headers,Content-Type,Access-Control-")
		ctx.Header("Access-Control-Allow-Credentials", "true")
		//如果是预请求，数据体无内容
		if method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}
	}
}
