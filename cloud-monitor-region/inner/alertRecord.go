package inner

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/validator/translate"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AlertRecordCtl struct {
}

func NewAlertRecordCtl() *AlertRecordCtl {
	return &AlertRecordCtl{}
}

// AddAlertRecord 创建告警记录
func AddAlertRecord(c *gin.Context) {
	var f forms.InnerAlertRecordAddForm
	if err := c.ShouldBindJSON(&f); err != nil {
		c.JSON(http.StatusBadRequest, translate.GetErrorMsg(err))
		return
	}

}
