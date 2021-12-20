package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/constants"
	dao2 "code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/enums/calcMode"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/errors"
	forms2 "code.cestc.cn/ccos-ops/cloud-monitor/business-common/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysComponent/sysRedis"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/dtos"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/k8s"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/utils"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
)

type K8sPrometheusService struct {
}

var PrometheusRule = &K8sPrometheusService{}

const ProductLabel = "hawkeye"
const ProductNamespaceLabel = "product-cec-hawkeye"

func (service *K8sPrometheusService) GenerateUserPrometheusRule(tenantId string) {
	ctxLock := context.Background()
	err := sysRedis.Lock(ctxLock, constants.TenantRuleKey, sysRedis.DefaultLease, true)
	if err != nil {
		log.Printf("获取 rule lock error  %+v", err)
		return
	}
	defer sysRedis.Unlock(ctxLock, constants.TenantRuleKey)
	alertRuleDTO, router, err := service.buildPrometheusRule("", "", tenantId)
	if err != nil {
		service.deleteK8sRule(tenantId, err, router)
		return
	}
	err = k8s.ApplyAlertRule(alertRuleDTO)
	if err != nil {
		logger.Logger().Infof("调用rule api apply 规格失败 %+v", err)
		return
	}
	err = k8s.ApplyAlertManagerConfig(*router)
	if err != nil {
		logger.Logger().Infof("调用alertmanager api apply 规格失败 %+v", err)
	}
}

func (service *K8sPrometheusService) deleteK8sRule(tenantId string, err error, router *k8s.AlertManagerConfig) {
	log.Printf(err.Error())
	businessError := err.(*errors.BusinessError)
	if businessError != nil && businessError.Code == errors.NoResource {
		err := k8s.DeleteAlertRule(tenantId)
		if err != nil {
			log.Printf("调用rule api delete 规格失败 %+v", err)
		}
		err = k8s.DeleteAlertManagerConfig(router.Name)
		if err != nil {
			log.Printf("调用alertmanager api delete 规格失败 %+v", err)
		}
	}
}

func (service *K8sPrometheusService) buildPrometheusRule(region string, zone string, tenantId string) (*forms.AlertRuleDTO, *k8s.AlertManagerConfig, error) {
	ruleDto := &forms.AlertRuleDTO{Region: region, Zone: zone, TenantId: tenantId}
	var alertList []*forms.AlertDTO
	var waitGroup = &sync.WaitGroup{}
	waitGroup.Add(2)
	alertListChan := make(chan []*forms.AlertDTO, 5)
	go service.buildAlertRuleListByResource(waitGroup, alertListChan, tenantId)
	go service.buildAlertRuleListByResourceGroup(waitGroup, alertListChan, tenantId)
	go func() {
		waitGroup.Wait()
		close(alertListChan)
	}()
	for list := range alertListChan {
		alertList = append(alertList, list...)
	}
	router := buildAlertManagerRouter(alertList, tenantId)
	if len(alertList) == 0 {
		return nil, router, errors.NewBusinessErrorCode(errors.NoResource, "instanceList 为空")
	}
	group := &forms.SpecGroups{Name: tenantId, AlertList: alertList}
	var groups []*forms.SpecGroups
	specGroups := append(groups, group)
	ruleDto.SpecGroupsList = specGroups
	return ruleDto, router, nil
}

func (service *K8sPrometheusService) buildAlertRuleListByResource(wg *sync.WaitGroup, resultChan chan []*forms.AlertDTO, tenantId string) {
	defer wg.Done()
	var resRuleList []*dtos.RuleExpress
	var alertList []*forms.AlertDTO
	global.DB.Raw("SELECT   t1.name as ruleName ,t1.`level`, t1.trigger_condition as ruleCondition, t1.id as ruleId,t1.product_type, t1.monitor_type ,t2.resource_id,t1.silences_time FROM  t_alarm_rule t1,  t_alarm_rule_resource_rel t2   WHERE  t2.alarm_rule_id = t1.id   AND t2.tenant_id = ?   AND t1.deleted = 0   AND t1.enabled = 1", tenantId).Scan(&resRuleList)
	for _, ruleExpress := range resRuleList {
		ruleExpress.NoticeGroupIds = dao2.AlarmRule.GetNoticeGroupList(global.DB, ruleExpress.RuleId)
		rule, err := service.buildAlertRule(ruleExpress, ruleExpress.ResourceId)
		if err != nil {
			fmt.Printf("build rule err %+v", err)
			continue
		}
		alertList = append(alertList, rule)
	}
	resultChan <- alertList
}

func (service *K8sPrometheusService) buildAlertRuleListByResourceGroup(wg *sync.WaitGroup, resultChan chan []*forms.AlertDTO, tenantId string) {
	defer wg.Done()
	var groupRuleList []*dtos.RuleExpress
	var alertList []*forms.AlertDTO
	global.DB.Raw("SELECT   t1.name as ruleName ,t1.`level`, t1.trigger_condition as ruleCondition, t1.id as ruleId,t1.product_type, t1.monitor_type ,t2.resource_group_id,t2.calc_mode ,t1.silences_time FROM  t_alarm_rule t1,  t_alarm_rule_group_rel t2   WHERE  t2.alarm_rule_id = t1.id   AND t2.tenant_id = ?   AND t1.deleted = 0   AND t1.enabled = 1", tenantId).Scan(&groupRuleList)
	for _, ruleExpress := range groupRuleList {
		ruleExpress.NoticeGroupIds = dao2.AlarmRule.GetNoticeGroupList(global.DB, ruleExpress.RuleId)
		instanceList := dao2.AlarmRule.GetResourceListByGroup(global.DB, ruleExpress.ResGroupId)
		if calcMode.ResourceGroup == ruleExpress.CalcMode {
			rule, err := service.buildAlertRule(ruleExpress, service.joinResourceId(instanceList, "|"))
			if err != nil {
				fmt.Printf("build rule err %+v", err)
				continue
			}
			alertList = append(alertList, rule)
		} else {
			alertRuleList := service.buildAlertRuleList(instanceList, ruleExpress)
			alertList = append(alertList, alertRuleList...)
		}
	}
	resultChan <- alertList
}

func (service *K8sPrometheusService) buildAlertRuleList(instanceList []*forms2.InstanceInfo, ruleExpress *dtos.RuleExpress) []*forms.AlertDTO {
	var alertList []*forms.AlertDTO
	for _, instance := range instanceList {
		rule, err := service.buildAlertRule(ruleExpress, instance.InstanceId)
		if err != nil {
			continue
		}
		alertList = append(alertList, rule)
	}
	return alertList
}

func (service *K8sPrometheusService) buildAlertRule(ruleExpress *dtos.RuleExpress, instanceId string) (*forms.AlertDTO, error) {
	alert := &forms.AlertDTO{}
	conditionId, err := utils.MD5(ruleExpress.RuleCondition)
	if err != nil {
		return nil, err
	}
	if len(ruleExpress.ResGroupId) == 0 {
		alert.Alert = fmt.Sprintf("%s#%s#%s", ruleExpress.RuleId, instanceId, conditionId)
	} else {
		alert.Alert = fmt.Sprintf("%s#group-%s#%s", ruleExpress.RuleId, ruleExpress.ResGroupId, conditionId)
	}
	alert.Expr = service.generateExpr(ruleExpress.RuleCondition, instanceId)
	alert.RuleType = "alert"
	alert.ForTime = utils.SecToTime(ruleExpress.RuleCondition.Times * ruleExpress.RuleCondition.Period)
	alert.Summary = service.getTemplateLabels(ruleExpress.RuleCondition.Labels)
	labelMaps := map[string]interface{}{}
	labelMaps["severity"] = dao2.ConfigItem.GetConfigItem(ruleExpress.Level, dao2.AlarmLevel, "").Name
	labelMaps["app"] = ProductLabel
	labelMaps["namespace"] = ProductNamespaceLabel
	alert.Labels = labelMaps
	silenceTime, err := strconv.Atoi(ruleExpress.SilencesTime)
	if err != nil {
		alert.SilencesTime = "3h"
	} else {
		alert.SilencesTime = utils.SecToTime(silenceTime)
	}
	desc, err := json.Marshal(ruleExpress.RuleCondition)
	alert.Description = string(desc)
	return alert, nil
}

func (service *K8sPrometheusService) generateExpr(ruleCondition *forms2.RuleCondition, instanceId string) string {
	return fmt.Sprintf("%s_over_time(%s{%s}[%s])%s%v", dao2.ConfigItem.GetConfigItem(ruleCondition.Statistics, dao2.StatisticalMethodsPid, "").Data,
		ruleCondition.MetricName, service.getLabels(instanceId, ruleCondition.Labels),
		utils.SecToTime(ruleCondition.Period), dao2.ConfigItem.GetConfigItem(ruleCondition.ComparisonOperator, dao2.ComparisonMethodPid, "").Data,
		ruleCondition.Threshold)
}

func (service *K8sPrometheusService) getLabels(instanceId string, labelStr string) string {
	builder := strings.Builder{}
	labels := strings.Split(labelStr, ",")
	for _, label := range labels {
		if strings.EqualFold(label, "instance") {
			builder.WriteString(label)
			builder.WriteString(fmt.Sprintf("=~'%s'", instanceId))
		}
	}
	return builder.String()
}

func (service *K8sPrometheusService) getTemplateLabels(labelStr string) string {
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

func (service *K8sPrometheusService) joinResourceId(elems []*forms2.InstanceInfo, sep string) string {
	size := len(elems)
	switch size {
	case 0:
		return ""
	case 1:
		return elems[0].InstanceId
	}
	n := len(sep) * (size - 1)
	for i := 0; i < size; i++ {
		n += len(elems[i].InstanceId)
	}

	var b strings.Builder
	b.Grow(n)
	b.WriteString(elems[0].InstanceId)
	for _, s := range elems[1:] {
		b.WriteString(sep)
		b.WriteString(s.InstanceId)
	}
	return b.String()
}

func buildAlertManagerRouter(alertList []*forms.AlertDTO, tenantId string) *k8s.AlertManagerConfig {
	router := make([]k8s.Router, len(alertList))
	for index, alertDto := range alertList {
		router[index] = k8s.Router{
			Matchers:       map[string]string{"ruleId": alertDto.Alert},
			RepeatInterval: alertDto.SilencesTime,
		}
	}
	return &k8s.AlertManagerConfig{
		Name:   "tenant-" + tenantId,
		Router: router,
	}
}
