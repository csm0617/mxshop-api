package global

import (
	"mxshop_api/user-web/config"
)

// 啥也不干，专门用来定义全局变量
var (
	ServerConfig *config.ServerConfig = &config.ServerConfig{} //因为经常要变化，所以采用指针并且初始化
)
