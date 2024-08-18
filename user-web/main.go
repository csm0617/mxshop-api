package main

import (
	"fmt"

	"mxshop_api/user-web/initialize"

	"go.uber.org/zap"
)

func main() {
	port := 8021
	//初始化全局Logger
	initialize.InitLogger()
	//初始化Router
	Routers := initialize.Routers()
	//拿到zap的全局sugar
	zap.S().Debugf("启动服务器，端口：%d", port)
	if err := Routers.Run(fmt.Sprintf(":%d", port)); err != nil {
		zap.S().Panic("启动失败：", err.Error())
	}

}
