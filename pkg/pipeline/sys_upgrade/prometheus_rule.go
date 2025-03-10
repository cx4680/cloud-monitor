package sys_upgrade

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global/sys_component/sys_db"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/k8s"
)

func PrometheusRuleUpgrade() {
	if !sys_db.Update {
		return
	}
	var tenantIds []string
	global.DB.Raw("SELECT DISTINCT   t1.tenant_id    FROM   t_resource t1    WHERE   t1.tenant_id != ''").Scan(&tenantIds)
	for _, tenantId := range tenantIds {
		k8s.PrometheusRule.GenerateUserPrometheusRule(tenantId)
	}
}
