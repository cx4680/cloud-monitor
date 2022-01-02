package inner

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/enum/source_type"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/api/controller"
	"github.com/gin-gonic/gin"
)

type AlarmRuleCtl struct {
}

func NewAlarmRuleCtl() *AlarmRuleCtl {
	return &AlarmRuleCtl{}
}

func (ctl *AlarmRuleCtl) CreateRule(c *gin.Context) {
	var param = form.AlarmRuleAddReqDTO{
		SourceType: source_type.AutoScaling,
	}
	controller.CreateRule(c, param)
}

func (ctl *AlarmRuleCtl) UpdateRule(c *gin.Context) {
	var param = form.AlarmRuleAddReqDTO{
		SourceType: source_type.AutoScaling,
	}
	controller.UpdateRule(c, param)
}
