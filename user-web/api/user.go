package api

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"mxshop_api/user-web/forms"
	"mxshop_api/user-web/global"
	"mxshop_api/user-web/global/reponse"
	"mxshop_api/user-web/middlewares"
	"mxshop_api/user-web/models"
	"mxshop_api/user-web/proto"
)

/*
*

		去掉validator的返回的错误信息key的包相关的结构
	    @param ctx
	    @param c
	    @author csm
	    @date 2024/8/18 15:36
*/
func removeTopStruct(fileds map[string]string) map[string]string {
	//初始化一个返回的map
	rsp := map[string]string{}
	for filed, err := range fileds {
		//将validator的返回的错误信息进行处理，去掉key前面的package内容作为新的key
		rsp[filed[strings.Index(filed, ".")+1:]] = err
	}
	//返回处理后的map
	return rsp
}

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

func HandleValidatorError(err error, ctx *gin.Context) {
	//怎么返回错误信息
	//1.先看能不能转成校验错误
	errs, ok := err.(validator.ValidationErrors)
	if !ok { //如果不是校验错误那就返回原有的错误
		ctx.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
	}
	//否则将返回参数错误，并将错误信息处理后返回
	ctx.JSON(http.StatusBadRequest, gin.H{
		"error": removeTopStruct(errs.Translate(global.Trans)),
	})
	return
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
	//打印当前访问用户的Id
	claims, _ := ctx.Get("claims")               //返回的是接口类型，至于为什么能拿到claims，因为在jwt登录校验中间件中，如果jwt有效，我们将claims和id设置到了ctx中
	currentUser := claims.(*models.CustomClaims) //转换成自定义声明的类型
	zap.S().Infof("当前访问[GetUserList]的用户id为：【%d】", currentUser.ID)
	pn := ctx.DefaultQuery("pn", "0")
	pnInt, _ := strconv.Atoi(pn)
	pSize := ctx.DefaultQuery("pSize", "10")
	pSizeInt, _ := strconv.Atoi(pSize)
	//生成grpc的client并调用接口
	userSrvClient := proto.NewUserClient(userConn)
	rsp, err := userSrvClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    uint32(pnInt),
		PSize: uint32(pSizeInt),
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
func PassWordLoginForm(ctx *gin.Context) {
	//1.先校验前端传过来的参数
	passWordLoginForm := forms.PassWordLoginForm{}
	if err := ctx.ShouldBindJSON(&passWordLoginForm); err != nil {
		HandleValidatorError(err, ctx)
		return
	}
	//2.连接user_srv_grpc
	ip := global.ServerConfig.UserSrvInfo.Host
	port := global.ServerConfig.UserSrvInfo.Port
	//grpc.WithInsecure()过时了，得用grpc.WithTransportCredentials(insecure.NewCredentials())方法代替
	userConn, err := grpc.NewClient(fmt.Sprintf("%s:%d", ip, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.S().Errorw("[PassWordLoginForm] 连接【用户服务失败】",
			"msg", err.Error())
	}
	//生成grpc的client并调用相关的接口
	userSrvClient := proto.NewUserClient(userConn)

	//登录逻辑
	//1.先查询用户是否存在，如果不存在给出提示
	if rsp, err := userSrvClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: passWordLoginForm.Mobile,
	}); err != nil { //注意这个err是grpc返回的err有自己的一套状态码
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound: //grpc没有找到
				ctx.JSON(http.StatusBadRequest, map[string]string{ //http状态码返回400，请求参数不正确
					"mobile": "用户不存在", //给出具体字段的错误描述
				})
			default: //其他错误，就提示一个登录失败
				ctx.JSON(http.StatusInternalServerError, map[string]string{
					"mobile": "登录失败",
				})
				//后台输出看看
				zap.S().Infof("grpc [GetUserByMobile] 返回错误:%s", err.Error())
			}
			//有错误就终止这个逻辑
			return
		}
	} else { //2.如果存在校验密码的正确性(grpc服务查询用户没有出错)
		//调用grpc服务[CheckPassword]将表单中输入的密码和数据库中查到用户的密码进行比对
		if passRsp, pasErr := userSrvClient.CheckPassword(context.Background(), &proto.PasswordCheckInfo{
			Password:          passWordLoginForm.PassWord,
			EncryptedPassword: rsp.Password,
		}); pasErr != nil { //grpc调用返回的错误
			ctx.JSON(http.StatusInternalServerError, map[string]string{
				"password": "登录失败",
			})
		} else { //grpc调用过程中没有出现错误
			if passRsp.Success { //密码校验成功
				//生成token
				j := middlewares.NewJWT()
				claims := models.CustomClaims{
					ID:          uint(rsp.Id),
					NickName:    rsp.Nickname,
					AuthorityId: uint(rsp.Role),
					StandardClaims: jwt.StandardClaims{
						NotBefore: time.Now().Unix(),                          //签名生效时间
						ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(), //签名过期时间
						Issuer:    "csm",                                      //签名的发行者
					},
				}
				token, err := j.CreateToken(claims)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, map[string]string{
						"token": "token生成失败",
					})
					return
				}
				ctx.JSON(http.StatusOK, gin.H{
					"id":         rsp.Id,
					"nick_name":  rsp.Nickname,
					"token":      token,
					"expired_at": claims.ExpiresAt * 1000,
				})
				zap.S().Infof("手机号为【%s】的用户登录成功", passWordLoginForm.Mobile)
			} else {
				ctx.JSON(http.StatusBadRequest, map[string]string{
					"password": "密码错误",
				})
			}
		}
	}

}
