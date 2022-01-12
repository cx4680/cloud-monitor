package inner

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/validator/translate"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type AlertRecordCtl struct {
	alertRecordAddSvc *service.AlertRecordAddService
}

func NewAlertRecordCtl(alertRecordAddSvc *service.AlertRecordAddService) *AlertRecordCtl {
	return &AlertRecordCtl{
		alertRecordAddSvc: alertRecordAddSvc,
	}
}

// AddAlertRecord 创建告警记录
func (ctl AlertRecordCtl) AddAlertRecord(c *gin.Context) {
	requestId := uuid.New().String()
	var f form.InnerAlertRecordAddForm
	if err := c.ShouldBindJSON(&f); err != nil {
		//TODO 添加requestId
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	logger.Logger().Info("requestId=", requestId, ", receive alarm data=", jsonutil.ToString(f))
	if err := ctl.alertRecordAddSvc.Add(requestId, f); err != nil {
		c.JSON(http.StatusInternalServerError, global.NewError("创建告警记录失败"))
	}
	c.JSON(http.StatusOK, global.NewSuccess("创建告警成功", nil))

}
