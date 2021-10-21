package middleware

import (
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("tenantId", "210310011310200")
		c.Next()
	}
}
