package middleware

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(global.TenantId, "210310011310200")
		c.Set(global.UserId, "210310011310200")
		c.Next()
	}
}
