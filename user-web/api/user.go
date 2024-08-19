package api

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"mxshop_api/user-web/global"
	"mxshop_api/user-web/global/reponse"
	"mxshop_api/user-web/proto"
)

/**
 * grpc返回的错误转为http返回
 * @param err
 * @param ctx
 * @return 无返回值，直接修改gin.Context的状态码和返回值
 * @author csm
 * @date 2024/8/18 16:23
 */
func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	//将grpc的code转换成http的状态码返回
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误", //内部错误不适合暴露给外界
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			case codes.Unavailable:
				c.JSON(http.StatusServiceUnavailable, gin.H{
					"msg": "用户服务不可用",
				})

			default: //其他的就不写了，可能有自定义的需要转换
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "其他错误" + e.Code().String(),
				})
			}
			return
		}
	}
}
func GetUserList(ctx *gin.Context) {
	ip := global.ServerConfig.UserSrvInfo.Host
	port := global.ServerConfig.UserSrvInfo.Port
	//grpc.WithInsecure()过时了，得用grpc.WithTransportCredentials(insecure.NewCredentials())方法代替
	userConn, err := grpc.NewClient(fmt.Sprintf("%s:%d", ip, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接【用户服务失败】",
			"msg", err.Error())
	}
	//生成grpc的client并调用接口
	userSrvClient := proto.NewUserClient(userConn)
	rsp, err := userSrvClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    1,
		PSize: 100,
	})
	if err != nil {
		zap.S().Errorw("[GetUserList] [查询用户列表]失败",
			"msg", err.Error())
		//错误处理
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		//data := make(map[string]interface{})

		//只向web层展示这五个字段，或者由前端决定
		user := reponse.UserResponse{
			Id:       value.Id,
			NickName: value.Nickname,
			Mobile:   value.Mobile,
			Gender:   value.Gender,
			Birthday: reponse.JsonTime(time.Unix(int64(value.BirthDay), 0)),
		}
		//data["id"] = value.Id
		//data["nickname"] = value.Nickname
		//data["mobile"] = value.Mobile
		//data["gender"] = value.Gender
		//data["birthday"] = value.BirthDay
		result = append(result, user)
	}
	ctx.JSON(http.StatusOK, result)
	zap.S().Debug("获取用户列表页")

}
