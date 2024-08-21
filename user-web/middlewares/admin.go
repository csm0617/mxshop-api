package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"mxshop_api/user-web/models"
)

/*
*
管理员权限校验中间件
*/
func IsAdminAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, _ := ctx.Get("claims")
		if authorityId := claims.(*models.CustomClaims).AuthorityId; authorityId != 2 {
			ctx.JSON(http.StatusForbidden, map[string]string{
				"msg": "无权限",
			})
			//无权限，不执行下一个中间件了
			ctx.Abort()
			return
		}
		//继续向下执行下一个中间
		ctx.Next()
	}
}
