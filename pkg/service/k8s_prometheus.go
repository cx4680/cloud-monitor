package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	commonConstant "code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/constant"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/dao"
	dtos2 "code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/dto"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/enum/calc_mode"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/enum/source_type"
	errors2 "code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/errors"
	form2 "code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global/sys_component/sys_redis"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/constant"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/dto"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/form"
	k8s2 "code.cestc.cn/ccos-ops/cloud-monitor/pkg/k8s"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/util"
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
	key := fmt.Sprintf(commonConstant.TenantRuleKey, tenantId)
	err := sys_redis.Lock(ctxLock, key, sys_redis.DefaultLease, true)
	if err != nil {
		log.Printf("获取 rule lock error  %+v", err)
		return
	}
	defer func(ctx context.Context, key string) {
		err := sys_redis.Unlock(ctx, key)
		if err != nil {
			logger.Logger().Errorf("unlock errorL%+v, lock:%s", err, key)
		}
	}(ctxLock, key)
	alertRuleDTO, router, err := service.buildPrometheusRule("", "", tenantId)
	if err != nil {
		service.deleteK8sRule(tenantId, err, router)
		return
	}

	err = k8s2.ApplyAlertRule(alertRuleDTO)
	if err != nil {
		logger.Logger().Infof("调用rule api apply 规格失败 %+v", err)
		return
	}
	if len(router.Router) == 0 {
		err = k8s2.DeleteAlertManagerConfig(router.Name)
		if err != nil {
			logger.Logger().Errorf("调用alertmanager api delete 规格失败 %+v", err)
		}
	} else {
		err = k8s2.ApplyAlertManagerConfig(*router)
		if err != nil {
			logger.Logger().Infof("调用alertmanager api apply 规格失败 %+v", err)
		}
	}
}

func (service *K8sPrometheusService) deleteK8sRule(tenantId string, err error, router *k8s2.AlertManagerConfig) {
	log.Printf(err.Error())
	businessError := err.(*errors2.BusinessError)
	if businessError != nil && businessError.Code == errors2.NoResource {
		err := k8s2.DeleteAlertRule(tenantId)
		if err != nil {
			logger.Logger().Errorf("调用rule api delete 规格失败 %+v", err)
		}
		err = k8s2.DeleteAlertManagerConfig(router.Name)
		if err != nil {
			logger.Logger().Errorf("调用alertmanager api delete 规格失败 %+v", err)
		}
	}
}

func (service *K8sPrometheusService) buildPrometheusRule(region string, zone string, tenantId string) (*form.AlertRuleDTO, *k8s2.AlertManagerConfig, error) {
	ruleDto := &form.AlertRuleDTO{Region: region, Zone: zone, TenantId: tenantId}
	var alertList []*form.AlertDTO
	var waitGroup = &sync.WaitGroup{}
	waitGroup.Add(2)
	alertListChan := make(chan []*form.AlertDTO, 5)
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
		return nil, router, errors2.NewBusinessErrorCode(errors2.NoResource, "instanceList 为空")
	}
	group := &form.SpecGroups{Name: tenantId, AlertList: alertList}
	var groups []*form.SpecGroups
	specGroups := append(groups, group)
	ruleDto.SpecGroupsList = specGroups
	return ruleDto, router, nil
}

func (service *K8sPrometheusService) buildAlertRuleListByResource(wg *sync.WaitGroup, resultChan chan []*form.AlertDTO, tenantId string) {
	defer wg.Done()
	var resRuleList []*dto.RuleExpress
	var alertList []*form.AlertDTO
	global.DB.Raw("SELECT   t1.name as ruleName ,t1.`level`, t1.trigger_condition as ruleCondition, t1.biz_id as ruleId,t1.product_name as product_type, t1.monitor_type ,t2.resource_id,t1.silences_time,t1.source_type FROM  t_alarm_rule t1,  t_alarm_rule_resource_rel t2   WHERE  t2.alarm_rule_id = t1.biz_id   AND t2.tenant_id = ?   AND t1.deleted = 0   AND t1.enabled = 1", tenantId).Scan(&resRuleList)
	for _, ruleExpress := range resRuleList {
		ruleExpress.NoticeGroupIds = dao.AlarmRule.GetNoticeGroupList(global.DB, ruleExpress.RuleId)
		ruleExpress.TenantId = tenantId
		rule, err := service.buildAlertRule(ruleExpress, ruleExpress.ResourceId)
		if err != nil {
			logger.Logger().Errorf("build rule err %+v", err)
			continue
		}
		alertList = append(alertList, rule)
	}
	resultChan <- alertList
}

func (service *K8sPrometheusService) buildAlertRuleListByResourceGroup(wg *sync.WaitGroup, resultChan chan []*form.AlertDTO, tenantId string) {
	defer wg.Done()
	var groupRuleList []*dto.RuleExpress
	var alertList []*form.AlertDTO
	global.DB.Raw("SELECT   t1.name as ruleName ,t1.`source` ,t1.`level`, t1.trigger_condition as ruleCondition, t1.biz_id as ruleId,t1.product_name as  product_type, t1.monitor_type ,t2.resource_group_id,t2.calc_mode ,t1.silences_time ,t1.source_type FROM  t_alarm_rule t1,  t_alarm_rule_group_rel t2   WHERE  t2.alarm_rule_id = t1.biz_id   AND t2.tenant_id = ?   AND t1.deleted = 0   AND t1.enabled = 1", tenantId).Scan(&groupRuleList)
	for _, ruleExpress := range groupRuleList {
		ruleExpress.NoticeGroupIds = dao.AlarmRule.GetNoticeGroupList(global.DB, ruleExpress.RuleId)
		instanceList := dao.AlarmRule.GetResourceListByGroup(global.DB, ruleExpress.ResGroupId)
		ruleExpress.TenantId = tenantId
		if calc_mode.Resource != ruleExpress.CalcMode {
			rule, err := service.buildAlertRule(ruleExpress, service.joinResourceId(instanceList, "|"))
			if err != nil {
				logger.Logger().Errorf("build rule err %+v", err)
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

func (service *K8sPrometheusService) buildAlertRuleList(instanceList []*form2.InstanceInfo, ruleExpress *dto.RuleExpress) []*form.AlertDTO {
	var alertList []*form.AlertDTO
	for _, instance := range instanceList {
		rule, err := service.buildAlertRule(ruleExpress, instance.InstanceId)
		if err != nil {
			continue
		}
		alertList = append(alertList, rule)
	}
	return alertList
}

func (service *K8sPrometheusService) buildAlertRule(ruleExpress *dto.RuleExpress, instanceId string) (*form.AlertDTO, error) {
	alert := &form.AlertDTO{}
	conditionId, err := util.MD5(ruleExpress.RuleCondition)
	if err != nil {
		return nil, err
	}
	resourceId := ""
	if len(ruleExpress.ResGroupId) == 0 || calc_mode.Resource == ruleExpress.CalcMode {
		alert.Alert = fmt.Sprintf("%s#%s#%s", ruleExpress.RuleId, instanceId, conditionId)
		resourceId = instanceId
	} else {
		alert.Alert = fmt.Sprintf("%s#group-%s#%s", ruleExpress.RuleId, ruleExpress.ResGroupId, conditionId)
	}
	alert.Expr = service.generateExpr(ruleExpress.RuleCondition, instanceId, ruleExpress.CalcMode)
	alert.RuleType = "alert"
	alert.ForTime = util.SecToTime(ruleExpress.RuleCondition.Times * ruleExpress.RuleCondition.Period)
	alert.Summary = service.getTemplateLabels(ruleExpress.RuleCondition.Labels, ruleExpress.CalcMode)
	labelMaps := map[string]interface{}{}
	labelMaps["severity"] = dao.ConfigItem.GetConfigItem(ruleExpress.Level, dao.AlarmLevel, "").Name
	labelMaps["app"] = ProductLabel
	source := "front"
	if source_type.AutoScaling == ruleExpress.SourceType {
		source = "autoScaling"
	}
	labelMaps["source"] = source
	labelMaps["namespace"] = ProductNamespaceLabel
	alert.Labels = labelMaps
	silenceTime, err := strconv.Atoi(ruleExpress.SilencesTime)
	if err != nil {
		alert.SilencesTime = "2h59s"
	} else {
		alert.SilencesTime = util.SecToTime(silenceTime - 1)
	}
	alert.SourceType = ruleExpress.SourceType
	noticeGroupIds := make([]string, len(ruleExpress.NoticeGroupIds))
	for index, noticeGroup := range ruleExpress.NoticeGroupIds {
		noticeGroupIds[index] = noticeGroup.Id
	}
	ruleDesc := dtos2.RuleDesc{
		RuleName:           ruleExpress.RuleName,
		Product:            ruleExpress.ProductType,
		MetricName:         ruleExpress.RuleCondition.MetricName,
		ComparisonOperator: ruleExpress.RuleCondition.ComparisonOperator,
		TargetValue:        ruleExpress.RuleCondition.Threshold,
		Time:               ruleExpress.RuleCondition.Times,
		Period:             ruleExpress.RuleCondition.Period,
		Unit:               ruleExpress.RuleCondition.Unit,
		Express:            dao.GetExpress(ruleExpress.RuleCondition),
		Level:              ruleExpress.Level,
		MonitorItem:        ruleExpress.RuleCondition.MonitorItemName,
		MonitorType:        ruleExpress.MonitorType,
		RuleId:             ruleExpress.RuleId,
		TenantId:           ruleExpress.TenantId,
		Statistic:          ruleExpress.RuleCondition.Statistics,
		GroupList:          noticeGroupIds,
		ResourceId:         resourceId,
		ResourceGroupId:    ruleExpress.ResGroupId,
		Source:             ruleExpress.Source,
		SourceType:         ruleExpress.SourceType,
	}
	desc, err := json.Marshal(ruleDesc)
	alert.Description = string(desc)
	return alert, nil
}

func (service *K8sPrometheusService) generateExpr(ruleCondition *form2.RuleCondition, instanceId string, mode int) string {
	monitorItem := dao.MonitorItem.GetMonitorItemByName(ruleCondition.MetricName)
	metric := strings.ReplaceAll(monitorItem.MetricsLinux, constant.MetricLabel, service.getLabels(instanceId, monitorItem.Labels))
	if calc_mode.ResourceGroup == mode {
		expr := fmt.Sprintf("%s_over_time((%s)[%s:1m])", dao.ConfigItem.GetConfigItem(ruleCondition.Statistics, dao.StatisticalMethodsPid, "").Data,
			metric, util.SecToTime(ruleCondition.Period))
		return fmt.Sprintf("%s(%s)%s%v", dao.ConfigItem.GetConfigItem(ruleCondition.Statistics, dao.StatisticalMethodsPid, "").Data, expr, dao.ConfigItem.GetConfigItem(ruleCondition.ComparisonOperator, dao.ComparisonMethodPid, "").Data, ruleCondition.Threshold)
	}
	return fmt.Sprintf("%s_over_time((%s)[%s:1m])%s%v", dao.ConfigItem.GetConfigItem(ruleCondition.Statistics, dao.StatisticalMethodsPid, "").Data,
		metric, util.SecToTime(ruleCondition.Period), dao.ConfigItem.GetConfigItem(ruleCondition.ComparisonOperator, dao.ComparisonMethodPid, "").Data,
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

func (service *K8sPrometheusService) getTemplateLabels(labelStr string, mode int) string {
	builder := strings.Builder{}
	builder.WriteString("currentValue={{$value}},")
	if calc_mode.ResourceGroup != mode {
		labels := strings.Split(labelStr, ",")
		for _, label := range labels {
			builder.WriteString(label)
			builder.WriteString("={{$labels.")
			builder.WriteString(label)
			builder.WriteString("}}")
			builder.WriteString(",")
		}
	}
	s := builder.String()
	return s[0:strings.LastIndex(s, ",")]
}

func (service *K8sPrometheusService) joinResourceId(elems []*form2.InstanceInfo, sep string) string {
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

func buildAlertManagerRouter(alertList []*form.AlertDTO, tenantId string) *k8s2.AlertManagerConfig {
	var router []k8s2.Router
	for _, alertDto := range alertList {
		if source_type.Front == alertDto.SourceType {
			continue
		}
		router = append(router, k8s2.Router{
			Matchers:       map[string]string{"alertname": alertDto.Alert},
			RepeatInterval: alertDto.SilencesTime,
		})
	}
	return &k8s2.AlertManagerConfig{
		Name:   "tenant-" + tenantId,
		Router: router,
	}
}
