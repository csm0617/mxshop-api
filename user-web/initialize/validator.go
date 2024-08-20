package initialize

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"

	"mxshop_api/user-web/global"
)

/**
* 初始化国际化翻译器
     * @param local 翻译器所需的语言环境
     * @return error 翻译器初始化失败时返回的错误信息，nil表示翻译器初始化成功
     *
     * 注意：gin.Validator.Engine()返回的是validator.Validate类型，
     * 而binding.Validator.Engine()返回的是validator.Validate类型，
     * 所以两者在功能上并无不同，但是binding.Validator.Engine()可以修改gin.Validator.Engine()的属性
     * 实现定制功能，在本例中，修改了gin.Validator.Engine()的validator.Validate类型，
     * 使之可以直接使用validator.Validate()来进行验证。
     * 而gin.Validator.Engine()返回的是gin.Validator类型，
     * 所以在gin.Validator.Engine()中无法直接使用validator.Validate()来进行验证。
*/
func InitTrans(local string) error {
	//修改gin框架中的Validator中的引擎属性，实现定制（被修改成validator.Validate）
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		//注册一个获取json的tag的自定义方法
		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		zhT := zh.New() //中文翻译器
		enT := en.New() //英文翻译器
		//第一个参数是备用的语言环境，后面的参数是应该支持的语言环境
		uti := ut.New(zhT, zhT, enT)
		var ok bool
		//将翻译器设置为全局变量
		global.Trans, ok = uti.GetTranslator(local)
		if !ok {
			return fmt.Errorf("uti.GetTranslator(%s)", local)
		}
		var err error
		switch local {
		case "en":
			err = en_translations.RegisterDefaultTranslations(v, global.Trans)
		case "zh":
			err = zh_translations.RegisterDefaultTranslations(v, global.Trans)
		default:
			err = en_translations.RegisterDefaultTranslations(v, global.Trans)
		}
		if err != nil {
			return err
		}
	}
	return nil
}
