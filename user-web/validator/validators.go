package validator

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

/*
*
自定义手机号校验规则
*/
func ValidateMobile(fl validator.FieldLevel) bool {
	//拿到要检验的字段
	mobile := fl.Field().String()
	//使用正则表达式库来判断mobile是手机号是否合法
	ok, _ := regexp.MatchString(`^1[3456789]\d{9}$`, mobile) //如果不想用\\来转译\，可以使用``
	if !ok {
		return false
	}
	return true
}
