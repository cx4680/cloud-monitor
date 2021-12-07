package iam

import (
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/pkg/middleware"
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/pkg/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthIdentify(identity *models.Identity) gin.HandlerFunc {
	return func(c *gin.Context) {
		// IAM鉴权接口
		err := middleware.AuthIdentify(c, identity, "")
		if err != nil {
			c.JSON(http.StatusOK, err)
			c.Abort()
		}
	}
}