package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"

	"github.com/Yifangmo/micro-shop-web/goods/global"
	"github.com/Yifangmo/micro-shop-web/goods/handlers"
	"github.com/Yifangmo/micro-shop-web/goods/initialize"
	"github.com/Yifangmo/micro-shop-web/goods/utils"
)

func main() {
	// 全局初始化
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitValidator()
	initialize.InitRPC()
	initialize.InitSentinel()
	zap.S().Infof("server port: %d", global.ServerConfig.Port)

	// 启动服务
	server := handlers.NewServer(fmt.Sprintf(":%d", global.ServerConfig.Port))
	done := make(chan struct{})
	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			zap.S().Error(err)
		}
		close(done)
	}()

	// consul 注册服务
	rc := utils.NewRegistryClient(global.ServerConfig.Consul.Host, global.ServerConfig.Consul.Port)
	serviceID := uuid.NewV4().String()
	err := rc.Register(global.ServerConfig.Host, global.ServerConfig.Port, global.ServerConfig.Name, global.ServerConfig.Tags, serviceID)
	if err != nil {
		zap.S().Panicf("Service register error: %v", err)
	}

	// 退出前注销服务
	defer func() {
		if err = rc.Deregister(serviceID); err != nil {
			zap.S().Error("Service deregister failed: ", err)
		} else {
			zap.S().Infof("Service[%s] deregister success", serviceID)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = server.Shutdown(ctx)
		if err != nil {
			zap.S().Error("Server shutdown error: ", err)
		}
	}()

	// 捕获退出信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-quit:
	case <-done:
	}
}
