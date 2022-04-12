package inner

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/validator/translate"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AlarmRecordCtl struct {
	alertRecordAddSvc *service.AlarmRecordAddService
}

func NewAlertRecordCtl(alertRecordAddSvc *service.AlarmRecordAddService) *AlarmRecordCtl {
	return &AlarmRecordCtl{
		alertRecordAddSvc: alertRecordAddSvc,
	}
}

// AddAlarmRecord 创建告警记录
func (ctl AlarmRecordCtl) AddAlarmRecord(c *gin.Context) {
	reqCtx := util.GenerateRequest(c)
	var f form.InnerAlarmRecordAddForm
	if err := c.ShouldBindJSON(&f); err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	logger.Logger().Info("requestId=", util.GetRequestId(reqCtx), ", receive alarm data=", jsonutil.ToString(f))
	if err := ctl.alertRecordAddSvc.Add(reqCtx, f); err != nil {
		c.JSON(http.StatusInternalServerError, global.NewError("创建告警记录失败"))
	}
	c.JSON(http.StatusOK, global.NewSuccess("创建告警成功", nil))

}
