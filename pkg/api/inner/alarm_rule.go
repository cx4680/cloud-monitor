package inner

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/api/controller"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/enum/source_type"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/form"
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
