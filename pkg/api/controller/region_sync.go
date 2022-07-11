package controller

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/validator/translate"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RegionSyncCtl struct {
	service *service.RegionSyncService
}

func NewRegionSyncController() *RegionSyncCtl {
	return &RegionSyncCtl{service.NewRegionSyncService()}
}

func (ctl *RegionSyncCtl) GetContactSyncData(c *gin.Context) {
	time := c.Query("time")
	if strutil.IsBlank(time) {
		c.JSON(http.StatusBadRequest, global.NewError("更新时间不能为空"))
		return
	}
	data, err := ctl.service.GetContactSyncData(time)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError("查询失败"))
		return
	}
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", data))
}

//func (ctl RegionSyncCtl) ContactSync(c *gin.Context) {
//	err := ctl.service.ContactSync()
//	if err != nil {
//		c.JSON(http.StatusOK, global.NewError("同步失败"))
//		return
//	} else {
//		c.JSON(http.StatusOK, global.NewSuccess("同步成功", true))
//	}
//}

func (ctl *RegionSyncCtl) GetAlarmRuleSyncData(c *gin.Context) {
	time := c.Query("time")
	if strutil.IsBlank(time) {
		c.JSON(http.StatusBadRequest, global.NewError("更新时间不能为空"))
		return
	}
	data, err := ctl.service.GetAlarmRuleSyncData(time)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError("查询失败"))
		return
	}
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", data))
}

//func (ctl RegionSyncCtl) AlarmRuleSync(c *gin.Context) {
//	err := ctl.service.AlarmRuleSync()
//	if err != nil {
//		c.JSON(http.StatusOK, global.NewError("同步失败"))
//		return
//	} else {
//		c.JSON(http.StatusOK, global.NewSuccess("同步成功", true))
//	}
//}

func (ctl RegionSyncCtl) GetAlarmRecordSyncData(c *gin.Context) {
	time := c.Query("time")
	if strutil.IsBlank(time) {
		c.JSON(http.StatusBadRequest, global.NewError("更新时间不能为空"))
		return
	}
	data, err := ctl.service.GetAlarmRecordSyncData(time)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError("查询失败"))
		return
	}
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", data))
}

func (ctl *RegionSyncCtl) PullAlarmRecordSyncData(c *gin.Context) {
	var param form.AlarmRecordSync
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	contactSync := ctl.service.PullAlarmRecordSyncData(param)
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", contactSync))
}

//func (ctl RegionSyncCtl) AlarmRecordSync(c *gin.Context) {
//	err := ctl.service.AlarmRecordSync()
//	if err != nil {
//		c.JSON(http.StatusOK, global.NewError(err.Error()))
//		return
//	} else {
//		c.JSON(http.StatusOK, global.NewSuccess("同步成功", true))
//	}
//}
