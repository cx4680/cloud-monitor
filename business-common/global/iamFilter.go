package global

import (
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/pkg/middleware"
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/pkg/models"
	"github.com/gin-gonic/gin"
)

func IamAuthIdentify(identity *models.Identity) gin.HandlerFunc {
	return func(c *gin.Context) {
		// IAM鉴权接口
		err := middleware.AuthIdentify(c, identity, "")
		if err != nil {
			c.JSON(200, err)
			c.Abort()
		}
	}
}