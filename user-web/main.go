package main

import (
	"fmt"
	ut "github.com/go-playground/universal-translator"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"mxshop_api/user-web/global"
	"mxshop_api/user-web/initialize"
	myValidator "mxshop_api/user-web/validator"
)

func main() {
	//初始化全局Logger
	initialize.InitLogger()
	//初始化validator翻译器
	if err := initialize.InitTrans("zh"); err != nil {
		zap.S().Errorw("初始化翻译器失败：", err.Error())
		panic(err)
	}
	//初始化全局Config
	initialize.InitConfig()
	//初始化Router
	Routers := initialize.Routers()
	//注册验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		//把自己定义的tag和校验规则放入验证器中
		v.RegisterValidation("mobileValidation", myValidator.ValidateMobile)
		//自定义校验器翻译问题（不知道怎么给出校验错误提示，需要我们自定义）
		v.RegisterTranslation("mobileValidation", global.Trans, func(ut ut.Translator) error {
			return ut.Add("mobileValidation", "{0}格式不正确", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobileValidation", fe.Field(), fe.Tag())
			return t
		})
	}

	//拿到zap的全局sugar
	zap.S().Debugf("启动服务器，端口：%d", global.ServerConfig.Port)
	if err := Routers.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
		zap.S().Panic("启动失败：", err.Error())
	}

}
