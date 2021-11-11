package inner

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/validator/translate"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AlertRecordCtl struct {
	alertRecordService *service.AlertRecordService
}

func NewAlertRecordCtl(alertRecordService *service.AlertRecordService) *AlertRecordCtl {
	return &AlertRecordCtl{
		alertRecordService: alertRecordService,
	}
}

// AddAlertRecord 创建告警记录
func (ctl AlertRecordCtl) AddAlertRecord(c *gin.Context) {
	var f forms.InnerAlertRecordAddForm
	if err := c.ShouldBindJSON(&f); err != nil {
		c.JSON(http.StatusBadRequest, translate.GetErrorMsg(err))
		return
	}
	err := ctl.alertRecordService.RecordAlertAndSendMessage(&f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "创建告警记录失败")
	}

}
