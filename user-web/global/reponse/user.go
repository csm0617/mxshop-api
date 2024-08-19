package reponse

import (
	"fmt"
	"time"
)

// 自定义一个JsonTime的类型别名
type JsonTime time.Time

// 除了拥有原来的Time类型的方法以外,还绑定了另外一个转换格式的方法
func (j JsonTime) MarshalJSON() ([]byte, error) {
	stmp := fmt.Sprintf("\"%s\"", time.Time(j).Format(time.DateOnly))
	return []byte(stmp), nil
}

type UserResponse struct {
	Id       int32    `json:"id"`
	NickName string   `json:"name"`
	Birthday JsonTime `json:"birthday"`
	Gender   string   `json:"gender"`
	Mobile   string   `json:"mobile"`
}
