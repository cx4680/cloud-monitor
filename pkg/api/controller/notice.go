package controller

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/constant"
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
		Total: constant.MaxSmsNum,
	}
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", noticeUsageVO))
}

func (nc *NoticeCtl) GetCenterUsage(c *gin.Context) {
	tenantId := c.Query("tenantId")
	if strutil.IsBlank(tenantId) {
		c.JSON(http.StatusOK, global.NewError("获取租户ID失败"))
		return
	}
	num, err := nc.service.GetTenantCurrentMonthSmsUsedNum(tenantId)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", num))
}

type NoticeUsageVO struct {
	Total int `json:"total"` //总数
	Used  int `json:"used"`  //已用数量
}

type NoticeChannelVO struct {
	ChannelList []string `json:"channelList"`
}
