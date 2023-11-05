package pay

import (
	"github.com/gin-gonic/gin"
)

func InitGroup(group *gin.RouterGroup) {
	payG := group.Group("pay")
	{
		payG.POST("alipay/notify", Notify)
	}
}
