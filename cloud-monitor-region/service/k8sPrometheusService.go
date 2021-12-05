package service

import (
	dao2 "code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysComponent/sysRedis"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/dtos"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/errors"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/k8s"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/utils"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type K8sPrometheusService struct {
	*dao2.AlarmRuleDao
}

var PrometheusRule = &K8sPrometheusService{AlarmRuleDao: dao2.AlarmRule}

func (dao *K8sPrometheusService) GenerateUserPrometheusRule(tenantId string) {
	alertRuleDTO, err := dao.buildPrometheusRule("", "", tenantId)
	if err != nil {
		logger.Logger().Infof(err.Error())
		return
	}
	key := "hawkeye-rule-" + tenantId
	result := true
	ctx := context.Background()
	for result {
		success, err := sysRedis.GetClient().SetNX(ctx, key, 1, time.Minute).Result()
		if err != nil {
			return
		}
		if !success {
			time.Sleep(time.Second)
		}
		result = !success
	}
	id := dao.getUserPrometheusId(tenantId)
	if len(id) > 0 {
		alertRuleDTO.AlertRuleId = id
		err := k8s.UpdateAlertRule(alertRuleDTO)
		if err != nil {
			logger.Logger().Errorf("规则更新失败%+v", err)
		}
	} else {
		_, err := k8s.CreateAlertRule(alertRuleDTO)
		if err != nil {
			logger.Logger().Errorf("规则创建失败%+v", err)
		} else {
			dao.createUserPrometheus(alertRuleDTO, tenantId)
		}
	}
	sysRedis.GetClient().Del(ctx, key)
}

func (dao *K8sPrometheusService) createUserPrometheus(alertRuleDTO *forms.AlertRuleDTO, tenantId string) {
	prometheus := &models.UserPrometheusID{
		PrometheusRuleID: alertRuleDTO.AlertRuleId,
		TenantID:         tenantId,
	}
	global.DB.Create(prometheus)
}

func (dao *K8sPrometheusService) buildPrometheusRule(region string, zone string, tenantId string) (*forms.AlertRuleDTO, error) {
	ruleDto := &forms.AlertRuleDTO{Region: region, Zone: zone, TenantId: tenantId}
	var list []*dtos.RuleExpress
	global.DB.Raw("select t1.name as ruleName ,t1.`level`, t1.trigger_condition as ruleCondition, t1.id as ruleId,t1.product_type, t1.notify_channel as noticeChannel,t1.monitor_type        from t_alarm_rule t1        where t1.tenant_id = ?         and t1.enabled = 1         and t1.deleted = 0", tenantId).Scan(&list)
	var alertList []*forms.AlertDTO
	for _, ruleExpress := range list {
		ruleExpress.GroupIds = dao.GetNoticeGroupList(ruleExpress.RuleId)
		ruleExpress.InstanceList = dao.GetInstanceList(ruleExpress.RuleId)
		for _, instance := range ruleExpress.InstanceList {
			alert := &forms.AlertDTO{}
			conditionId, err := utils.MD5(ruleExpress.RuleCondition)
			if err != nil {
				continue
			}
			alert.Alert = fmt.Sprintf("%s#%s#%s", ruleExpress.RuleId, instance.InstanceId, conditionId)
			alert.Expr = fmt.Sprintf("%s_over_time(%s{%s}[%s])%s%v", dao2.GetConfigItem(ruleExpress.RuleCondition.Statistics, "3", "").Data, ruleExpress.RuleCondition.MetricName, dao.getLabels(instance.InstanceId, ruleExpress.RuleCondition.Labels), utils.SecToTime(ruleExpress.RuleCondition.Period), dao2.GetConfigItem(ruleExpress.RuleCondition.ComparisonOperator, "4", "").Data, ruleExpress.RuleCondition.Threshold)
			alert.RuleType = "alert"
			alert.ForTime = utils.SecToTime(ruleExpress.RuleCondition.Times * ruleExpress.RuleCondition.Period)
			alert.Summary = dao.getTemplateLabels(ruleExpress.RuleCondition.Labels)
			labelMaps := map[string]interface{}{}
			labelMaps["severity"] = dao2.GetConfigItem(ruleExpress.Level, "28", "").Name
			labelMaps["app"] = "hawkeye"
			labelMaps["namespace"] = "product-cec-hawkeye"
			alert.Labels = labelMaps
			desc, err := json.Marshal(ruleExpress.RuleCondition)
			alert.Description = string(desc)
			alertList = append(alertList, alert)
		}
	}
	if len(alertList) == 0 {
		return nil, errors.NewBussinessError(1, "instanceList 为空")
	}
	group := &forms.SpecGroups{Name: tenantId, AlertList: alertList}
	var groups []*forms.SpecGroups
	specGroups := append(groups, group)
	ruleDto.SpecGroupsList = specGroups
	return ruleDto, nil
}

func (dao *K8sPrometheusService) getLabels(instanceId string, labelStr string) string {
	builder := strings.Builder{}
	labels := strings.Split(labelStr, ",")
	for _, label := range labels {
		if strings.EqualFold(label, "instance") {
			builder.WriteString(label)
			builder.WriteString(fmt.Sprintf("='%s'", instanceId))
		}
	}
	return builder.String()
}

func (dao *K8sPrometheusService) getTemplateLabels(labelStr string) string {
	builder := strings.Builder{}
	builder.WriteString("currentValue={{$value}},")
	labels := strings.Split(labelStr, ",")
	for _, label := range labels {
		builder.WriteString(label)
		builder.WriteString("={{$labels.")
		builder.WriteString(label)
		builder.WriteString("}}")
		builder.WriteString(",")
	}
	s := builder.String()
	return s[0:strings.LastIndex(s, ",")]
}

func (dao *K8sPrometheusService) getUserPrometheusId(tenantId string) string {
	prometheus := &models.UserPrometheusID{
		TenantID: tenantId,
	}
	global.DB.Find(prometheus)
	return prometheus.PrometheusRuleID
}
