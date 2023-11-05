package initialize

import (
	"fmt"

	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/Yifangmo/micro-shop-services/user/proto"
	"github.com/Yifangmo/micro-shop-web/user/global"
)

func InitRPC() {
	consulConfig := global.ServerConfig.Consul
	userConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulConfig.Host, consulConfig.Port, global.ServerConfig.UserService.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【用户服务失败】")
	}

	global.UserSrvClient = proto.NewUserClient(userConn)
}
