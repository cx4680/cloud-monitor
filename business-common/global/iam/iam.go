package iam

import (
	"net/http"

	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/pkg/middleware"
	"code.cestc.cn/yyptb-group_tech/iam-sdk-go/pkg/models"
	"github.com/gin-gonic/gin"
)

func AuthIdentify(identity *models.Identity) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("Product", identity.Product)
		c.Set("Action", identity.Action)
		c.Set("ResourceType", identity.ResourceType)
		c.Set("ResourceId", identity.ResourceId)
		// IAM鉴权接口
		err := middleware.AuthIdentify(c, identity, "")
		if err != nil {
			c.JSON(http.StatusInternalServerError, global.NewError("用户不允许进行该操作"))
			c.Abort()
		}
	}
}
