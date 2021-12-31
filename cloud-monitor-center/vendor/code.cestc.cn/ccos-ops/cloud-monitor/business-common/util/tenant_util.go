package util

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/errors"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"github.com/gin-gonic/gin"
)

func GetTenantId(c *gin.Context) (string, error) {
	tenantId := c.GetString(global.TenantId)
	if strutil.IsBlank(tenantId) {
		return "", errors.NewBusinessError("获取租户ID失败")
	}
	return tenantId, nil
}

func GetUserId(c *gin.Context) (string, error) {
	tenantId := c.GetString(global.UserId)
	if strutil.IsBlank(tenantId) {
		return "", errors.NewBusinessError("获取用户ID失败")
	}
	return tenantId, nil
}
