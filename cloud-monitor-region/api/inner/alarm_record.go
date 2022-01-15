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
	requestId := uuid.New().String()
	var f form.InnerAlarmRecordAddForm
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
