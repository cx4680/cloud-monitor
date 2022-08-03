package util

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/errors"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
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
	userId := c.GetString(global.UserId)
	if strutil.IsBlank(userId) {
		return "", errors.NewBusinessError("获取用户ID失败")
	}
	return userId, nil
}

func GetTenantIdAndUserId(c *gin.Context) (string, string, error) {
	tenantId := c.GetString(global.TenantId)
	if strutil.IsBlank(tenantId) {
		return "", "", errors.NewBusinessError("获取租户ID失败")
	}
	userId := c.GetString(global.UserId)
	if strutil.IsBlank(tenantId) {
		return "", "", errors.NewBusinessError("获取用户ID失败")
	}
	return tenantId, userId, nil
}
