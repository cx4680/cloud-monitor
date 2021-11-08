package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/redis"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/dtos"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/errors"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/k8s"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/utils"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"context"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"strings"
	"time"
)

type PrometheusRuleDao struct {
	db *gorm.DB
	dao.AlarmRuleDao
}

func NewPrometheusRuleDao(db *gorm.DB) *PrometheusRuleDao {
	return &PrometheusRuleDao{db: db}
}

func (dao *PrometheusRuleDao) GenerateUserPrometheusRule(region string, zone string, tenantId string) {
	alertRuleDTO, err := dao.buildPrometheusRule("", "", tenantId)
	if err != nil {
		logger.Logger().Infof(err.Error())
		return
	}
	var ctx context.Context
	key := "hawkeye-rule-" + tenantId
	result := true
	for result {
		success, err := redis.GetClient().SetNX(ctx, key, 1, time.Minute).Result()
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
		k8s.UpdateAlertRule(alertRuleDTO)
	} else {
		k8s.CreateAlertRule(alertRuleDTO)
		dao.createUserPrometheus(alertRuleDTO, tenantId)
	}
	redis.GetClient().Del(ctx, key)
}

func (dao *PrometheusRuleDao) createUserPrometheus(alertRuleDTO *forms.AlertRuleDTO, tenantId string) {
	prometheus := &models.UserPrometheusID{
		PrometheusRuleID: alertRuleDTO.AlertRuleId,
		TenantID:         tenantId,
	}
	dao.db.Create(prometheus)
}

func (dao *PrometheusRuleDao) buildPrometheusRule(region string, zone string, tenantId string) (*forms.AlertRuleDTO, error) {
	ruleDto := &forms.AlertRuleDTO{Region: region, Zone: zone, TenantId: tenantId}
	var list []*dtos.RuleExpress
	dao.db.Raw("select t1.name as ruleName ,t1.`level`, t1.trigger_condition as ruleCondition, t1.id as ruleId,t1.product_type, t1.notify_channel as noticeChannel,t1.monitor_type        from t_alarm_rule t1        where t1.tenant_id = ?         and t1.enabled = 1         and t1.deleted = 0", tenantId).Scan(&list)
	var alertList []*forms.AlertDTO
	for _, ruleExpress := range list {
		ruleExpress.GroupIds = dao.GetNoticeGroupList(ruleExpress.RuleId)
		for _, instance := range ruleExpress.InstanceList {
			alert := &forms.AlertDTO{}
			conditionId, err := utils.MD5(ruleExpress.RuleCondition)
			if err != nil {
				continue
			}
			alert.Alert = fmt.Sprintf("%s#%s#%s", ruleExpress.RuleId, instance.InstanceId, conditionId)
			alert.Expr = fmt.Sprintf("%s_over_time(%s{%s}[%s])%s%v", ruleExpress.RuleCondition.Statistics, ruleExpress.RuleCondition.MetricName, dao.getLabels(instance.InstanceId, ruleExpress.RuleCondition.Labels), utils.SecToTime(ruleExpress.RuleCondition.Period), ruleExpress.RuleCondition.ComparisonOperator, ruleExpress.RuleCondition.Threshold)
			alert.RuleType = "alert"
			alert.ForTime = utils.SecToTime(ruleExpress.RuleCondition.Times * ruleExpress.RuleCondition.Period)
			alert.Summary = dao.getTemplateLabels(ruleExpress.RuleCondition.Labels)
			labelMaps := map[string]interface{}{}
			labelMaps["severity"] = ruleExpress.Level
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

func (dao *PrometheusRuleDao) getLabels(instanceId string, labelStr string) string {
	builder := strings.Builder{}
	labels := strings.Split(labelStr, ",")
	for _, label := range labels {
		if strings.EqualFold(label, "instance") {
			builder.WriteString(label)
			builder.WriteString("='")
			builder.WriteString(instanceId)
			builder.WriteString("='")
		}
	}
	return builder.String()
}

func (dao *PrometheusRuleDao) getTemplateLabels(labelStr string) string {
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

func (dao *PrometheusRuleDao) getUserPrometheusId(tenantId string) string {
	prometheus := &models.UserPrometheusID{
		TenantID: tenantId,
	}
	dao.db.Find(prometheus)
	return prometheus.PrometheusRuleID
}
