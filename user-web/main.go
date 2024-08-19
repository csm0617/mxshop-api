package main

import (
	"fmt"

	"mxshop_api/user-web/global"
	"mxshop_api/user-web/initialize"

	"go.uber.org/zap"
)

func main() {
	//初始化全局Logger
	initialize.InitLogger()
	//初始化全局Config
	initialize.InitConfig()
	//初始化Router
	Routers := initialize.Routers()
	//拿到zap的全局sugar
	zap.S().Debugf("启动服务器，端口：%d", global.ServerConfig.Port)
	if err := Routers.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
		zap.S().Panic("启动失败：", err.Error())
	}

}
