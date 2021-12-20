package controllers

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/constants"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type NoticeCtl struct {
	service service.MessageService
}

func NewNoticeCtl(service service.MessageService) *NoticeCtl {
	return &NoticeCtl{service}
}

func (nc *NoticeCtl) GetUsage(c *gin.Context) {
	tenantId, exists := c.Get(global.TenantId)
	if !exists {
		c.JSON(http.StatusOK, global.NewError("获取租户ID失败"))
		return
	}
	num, err := nc.service.GetTenantCurrentMonthSmsUsedNum(tenantId.(string))
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
		return
	}
	var noticeUsageVO = &NoticeUsageVO{
		Used:  num,
		Total: constants.MaxSmsNum,
	}
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", noticeUsageVO))
}

type NoticeUsageVO struct {
	Total int `json:"total"` //总数
	Used  int `json:"used"`  //已用数量
}
