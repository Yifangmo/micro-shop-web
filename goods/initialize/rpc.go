package initialize

import (
	"fmt"

	"github.com/Yifangmo/micro-shop-services/goods/proto"
	"github.com/Yifangmo/micro-shop-web/goods/global"

	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func InitRPC() {
	consulConfig := global.ServerConfig.Consul
	conn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulConfig.Host, consulConfig.Port, global.ServerConfig.GoodsRPC.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		//grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		zap.S().Fatalf("InitRPC error: %v", err)
	}

	global.GoodsSrvClient = proto.NewGoodsClient(conn)
}
