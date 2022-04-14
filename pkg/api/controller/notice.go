package controller

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/constant"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/form"
	global2 "code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type NoticeCtl struct {
	service service.MessageService
}

func NewNoticeCtl(service service.MessageService) *NoticeCtl {
	return &NoticeCtl{service}
}

func (nc *NoticeCtl) GetUsage(c *gin.Context) {
	tenantId, exists := c.Get(global2.TenantId)
	if !exists {
		c.JSON(http.StatusOK, global2.NewError("获取租户ID失败"))
		return
	}
	num, err := nc.service.GetTenantCurrentMonthSmsUsedNum(tenantId.(string))
	if err != nil {
		c.JSON(http.StatusOK, global2.NewError(err.Error()))
		return
	}
	var noticeUsageVO = &NoticeUsageVO{
		Used:  num,
		Total: constant.MaxSmsNum,
	}
	c.JSON(http.StatusOK, global2.NewSuccess("查询成功", noticeUsageVO))
}

func (nc *NoticeCtl) GetCenterUsage(c *gin.Context) {
	tenantId := c.Query("tenantId")
	if strutil.IsBlank(tenantId) {
		c.JSON(http.StatusOK, global2.NewError("获取租户ID失败"))
		return
	}
	num, err := nc.service.GetTenantCurrentMonthSmsUsedNum(tenantId)
	if err != nil {
		c.JSON(http.StatusOK, global2.NewError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, global2.NewSuccess("查询成功", num))
}

func (nc *NoticeCtl) ChangeNoticeChannel(c *gin.Context) {
	if config.Cfg.Common.MsgIsOpen == config.MsgClose || strutil.IsBlank(config.Cfg.Common.MsgChannel) {
		c.JSON(http.StatusOK, global2.NewSuccess("该环境无告警渠道", false))
		return
	}
	var noticeChannelVO NoticeChannelVO
	err := c.ShouldBindJSON(&noticeChannelVO)
	if err != nil {
		c.JSON(http.StatusOK, global2.NewError(err.Error()))
		return
	}
	global2.NoticeChannelList = []form.NoticeChannel{}
	for _, v := range noticeChannelVO.ChannelList {
		if !strings.Contains(config.Cfg.Common.MsgChannel, v) {
			continue
		}
		switch v {
		case config.MsgChannelEmail:
			global2.NoticeChannelList = append(global2.NoticeChannelList, form.NoticeChannel{Name: "邮箱", Code: v, Data: 1})
		case config.MsgChannelSms:
			global2.NoticeChannelList = append(global2.NoticeChannelList, form.NoticeChannel{Name: "短信", Code: v, Data: 2})
		}
	}
	c.JSON(http.StatusOK, global2.NewSuccess("修改成功", true))
}

type NoticeUsageVO struct {
	Total int `json:"total"` //总数
	Used  int `json:"used"`  //已用数量
}

type NoticeChannelVO struct {
	ChannelList []string `json:"channelList"`
}
