package forms

type SendSmsForm struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required,mobileValidation"`
	Type   int    `form:"type" json:"type" binding:"required,oneof=1 2"`
}
