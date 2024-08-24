package initialize

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"mxshop_api/user-web/global"
	"mxshop_api/user-web/proto"
)

func InitSrvClient() {
	zap.S().Info("正在连接user-grpc服务，并初始化UserSrvClient...")
	//先设置配置信息
	config := api.DefaultConfig()
	consulInfo := global.ServerConfig.ConsulInfo
	config.Address = fmt.Sprintf("%s:%d", consulInfo.Host, consulInfo.Port)
	//根据配置信息拿到consul客户端实例
	client, err := api.NewClient(config)
	if err != nil {
		zap.S().Fatal("consul client init failed")
		panic(err)
	}
	userSrvName := global.ServerConfig.UserSrvInfo.Name
	filterStr := fmt.Sprintf(`Service == "%s"`, userSrvName)
	userSrvs, err := client.Agent().ServicesWithFilter(filterStr)
	if err != nil {
		zap.S().Infof("服务发现获取" + userSrvName + "服务失败：" + err.Error())
		return
	}
	var userSrvHost string
	var userSrvPort int
	for _, service := range userSrvs {
		userSrvHost = service.Address
		userSrvPort = service.Port
	}
	if userSrvHost == "" || userSrvPort == 0 {
		zap.S().Fatalf(`[InitSrvClient] 获取【%s服务的ip和port失败】`, userSrvName)
		return
	}
	//grpc.WithInsecure()过时了，得用grpc.WithTransportCredentials(insecure.NewCredentials())方法代替
	userConn, err := grpc.NewClient(fmt.Sprintf("%s:%d", userSrvHost, userSrvPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.S().Fatalf("[InitSrvClient] 连接【用户服务失败】",
			"msg", err.Error())
		return
	}
	//生成grpc的client，todo 1.服务下线，或者改了ip和端口怎么办？ 2.多个go携程并发操作一个Conn的性能问题
	global.UserSrvClient = proto.NewUserClient(userConn)
	zap.S().Info("初始化UserSrvClient成功！")
}
