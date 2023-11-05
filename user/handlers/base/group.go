package base

import (
	"github.com/gin-gonic/gin"
)

func InitGroup(g *gin.RouterGroup) {
	baseG := g.Group("base")
	{
		baseG.GET("captcha", GenerateCaptcha)
		baseG.POST("send_sms", SendSms)
	}
}
