package iam

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	openapi2 "code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global/openapi"
	"net/http"

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
			if openapi2.OpenApiRouter(c) {
				c.JSON(http.StatusOK, openapi2.NewRespError(openapi2.AuthorizedNoPermission, c))
			} else {
				c.JSON(http.StatusOK, global.NewError("对不起，您没有操作权限"))
			}
			c.Abort()
		}
	}
}
