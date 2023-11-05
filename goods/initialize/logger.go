package initialize

import (
	"github.com/Yifangmo/micro-shop-web/goods/global"
	"go.uber.org/zap"
)

func InitLogger() {
	logger, _ := zap.NewProduction()
	if global.IsDebug {
		logger, _ = zap.NewDevelopment()
	}
	zap.ReplaceGlobals(logger)
}
