package sys_upgrade

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/service"
)

func PrometheusRuleUpgrade() {
	var tenantIds []string
	global.DB.Raw("SELECT DISTINCT   t1.tenant_id    FROM   t_alarm_instance t1    WHERE   t1.tenant_id != ''").Scan(&tenantIds)
	for _, tenantId := range tenantIds {
		service.PrometheusRule.GenerateUserPrometheusRule(tenantId)
	}
}
