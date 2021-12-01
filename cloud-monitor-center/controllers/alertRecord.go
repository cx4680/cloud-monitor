package controllers

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/validator/translate"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AlertRecordController struct {
}

func NewAlertRecordController() *AlertRecordController {
	return &AlertRecordController{}
}

func (a *AlertRecordController) GetPageList(c *gin.Context) {
	var f forms.AlertRecordPageQueryForm
	if err := c.ShouldBindJSON(&f); err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	tenantId := c.GetString(global.TenantId)
	page := dao.AlertRecord.GetPageList(global.DB, tenantId, f)
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", page))
}

func (a *AlertRecordController) GetDetail(c *gin.Context) {
	recordId := c.Query("recordId")
	if tools.IsBlank(recordId) {
		c.JSON(http.StatusBadRequest, global.NewError("参数异常"))
		return
	}
	c.JSON(http.StatusOK, dao.AlertRecord.GetById(global.DB, recordId))
}
