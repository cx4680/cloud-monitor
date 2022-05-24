package k8s

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/enum/calc_mode"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/enum/source_type"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/errors"
	form2 "code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global/sys_component/sys_redis"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/model"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/constant"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/util"
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

type K8sPrometheusService struct {
}

var PrometheusRule = &K8sPrometheusService{}

const ProductLabel = "hawkeye"
const ProductNamespaceLabel = "product-cec-hawkeye"
const defaultAlarmInterval = "2h59s"

func (service *K8sPrometheusService) GenerateUserPrometheusRule(tenantId string) {
	ctxLock := context.Background()
	key := fmt.Sprintf(constant.TenantRuleKey, tenantId)
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

	err = ApplyAlertRule(alertRuleDTO)
	if err != nil {
		logger.Logger().Infof("调用rule api apply 规格失败 %+v", err)
		return
	}
	if len(router.Router) == 0 {
		err = DeleteAlertManagerConfig(router.Name)
		if err != nil {
			logger.Logger().Errorf("调用alertmanager api delete 规格失败 %+v", err)
		}
	} else {
		err = ApplyAlertManagerConfig(*router)
		if err != nil {
			logger.Logger().Infof("调用alertmanager api apply 规格失败 %+v", err)
		}
	}
}

func (service *K8sPrometheusService) deleteK8sRule(tenantId string, err error, router *AlertManagerConfig) {
	log.Printf(err.Error())
	businessError := err.(*errors.BusinessError)
	if businessError != nil && businessError.Code == errors.NoResource {
		err := DeleteAlertRule(tenantId)
		if err != nil {
			logger.Logger().Errorf("调用rule api delete 规格失败 %+v", err)
		}
		err = DeleteAlertManagerConfig(router.Name)
		if err != nil {
			logger.Logger().Errorf("调用alertmanager api delete 规格失败 %+v", err)
		}
	}
}

func (service *K8sPrometheusService) buildPrometheusRule(region string, zone string, tenantId string) (*AlertRuleDTO, *AlertManagerConfig, error) {
	ruleDto := &AlertRuleDTO{Region: region, Zone: zone, TenantId: tenantId}
	var alertList []AlertDTO
	var waitGroup = &sync.WaitGroup{}
	waitGroup.Add(2)
	alertListChan := make(chan []AlertDTO, 5)
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
	group := &SpecGroups{Name: tenantId, AlertList: alertList}
	var groups []*SpecGroups
	specGroups := append(groups, group)
	ruleDto.SpecGroupsList = specGroups
	return ruleDto, router, nil
}

const getResourceRuleInfoSql = " SELECT\nt1.tenant_id,t1.NAME rule_name,t1.`level`,t1.biz_id rule_id,t1.product_name  product_type,t1.monitor_type,t1.silences_time,t1.source_type,t1.source,t1.type, t1.combination, t1.period,t1.times,t2.resource_id,GROUP_CONCAT( t3.contract_group_id ) group_ids\n " +
	"from t_alarm_rule t1\n " +
	"join t_alarm_rule_resource_rel t2 on t1.biz_id=t2.alarm_rule_id\n " +
	"left JOIN t_alarm_notice t3 on t1.biz_id=t3.alarm_rule_id\n " +
	"left join t_contact_group t4 on t3.contract_group_id=t4.biz_id\n " +
	"WHERE t2.tenant_id = ? AND t1.deleted = 0   AND t1.enabled = 1\n " +
	"group by t1.biz_id, t2.resource_id"

const getResourceGroupRuleInfoSql = "SELECT\nt1.tenant_id,t1.`name` rule_name,t1.`source`,t1.`level`,t1.biz_id  rule_id,t1.product_name product_type,t1.monitor_type,t1.silences_time,t1.source_type,t1.type,t1.combination,t1.period,t1.times,t2.resource_group_id,t2.calc_mode,GROUP_CONCAT(t3.contract_group_id) groupIds\n" +
	"FROM\nt_alarm_rule t1\n" +
	"join t_alarm_rule_group_rel t2 on t2.alarm_rule_id = t1.biz_id \n" +
	"left join t_alarm_notice t3  on t1.biz_id = t3.alarm_rule_id\n" +
	"WHERE\n\n\t t2.tenant_id = ? \n\tAND t1.deleted = 0 \n\tAND t1.enabled = 1\n\t\n\t" +
	"group by t1.biz_id, t2.resource_group_id"

type RuleInfo struct {
	RuleId       string
	TenantId     string
	RuleName     string
	ProductType  string
	MonitorType  string
	SilencesTime string
	Level        uint8
	SourceType   uint8
	Source       string

	Type        uint8
	Combination uint8
	Period      int
	Times       int

	ResourceId string

	GroupIds string
}

func (info *RuleInfo) Parse() model.AlarmRule {
	return model.AlarmRule{
		Id:             0,
		BizId:          info.RuleId,
		MonitorType:    info.MonitorType,
		ProductName:    info.ProductType,
		ProductBizId:   "",
		Dimensions:     0,
		Name:           info.RuleName,
		SilencesTime:   info.SilencesTime,
		EffectiveStart: "",
		EffectiveEnd:   "",
		Level:          info.Level,
		Enabled:        0,
		TenantID:       info.TenantId,
		CreateTime:     time.Time{},
		CreateUser:     "",
		Deleted:        0,
		UserName:       "",
		UpdateTime:     time.Time{},
		Source:         info.Source,
		SourceType:     info.SourceType,
		Type:           info.Type,
		Combination:    info.Combination,
		Period:         info.Period,
		Times:          info.Times,
	}
}

type GroupRuleInfo struct {
	RuleId       string
	TenantId     string
	RuleName     string
	ProductType  string
	MonitorType  string
	SilencesTime string
	Level        uint8
	SourceType   uint8
	Source       string

	Type        uint8
	Combination uint8
	Period      int
	Times       int

	ResourceGroupId string
	CalcMode        int

	GroupIds string
}

func (info *GroupRuleInfo) Parse() model.AlarmRule {
	return model.AlarmRule{
		Id:             0,
		BizId:          info.RuleId,
		MonitorType:    info.MonitorType,
		ProductName:    info.ProductType,
		ProductBizId:   "",
		Dimensions:     0,
		Name:           info.RuleName,
		SilencesTime:   info.SilencesTime,
		EffectiveStart: "",
		EffectiveEnd:   "",
		Level:          info.Level,
		Enabled:        0,
		TenantID:       info.TenantId,
		CreateTime:     time.Time{},
		CreateUser:     "",
		Deleted:        0,
		UserName:       "",
		UpdateTime:     time.Time{},
		Source:         info.Source,
		SourceType:     info.SourceType,
		Type:           info.Type,
		Combination:    info.Combination,
		Period:         info.Period,
		Times:          info.Times,
	}
}

func (service *K8sPrometheusService) buildAlertRuleListByResource(wg *sync.WaitGroup, resultChan chan []AlertDTO, tenantId string) {
	defer wg.Done()
	var ruleInfoList []RuleInfo
	global.DB.Raw(getResourceRuleInfoSql, tenantId).Scan(&ruleInfoList)

	var alertList []AlertDTO
	for _, info := range ruleInfoList {
		items := dao.AlarmItem.GetItemListByRuleBizId(global.DB, info.RuleId)
		alerts, err := service.buildAlarmInfo(info, items)
		if err != nil {
			logger.Logger().Errorf("build rule err %+v", err)
			continue
		}
		alertList = append(alertList, alerts...)
	}
	resultChan <- alertList
}

func (service *K8sPrometheusService) buildAlarmInfo(ruleInfo RuleInfo, alarmItems []model.AlarmItem) ([]AlertDTO, error) {
	conditionId, err := util.MD5(alarmItems)
	if err != nil {
		return nil, err
	}
	var alertName = fmt.Sprintf("%s#%s#%s", ruleInfo.RuleId, ruleInfo.ResourceId, conditionId)

	var desc = AlarmDescription{
		Rule:            ruleInfo.Parse(),
		ContactGroupIds: strings.Split(ruleInfo.GroupIds, ","),
		ResourceId:      ruleInfo.ResourceId,
	}

	if ruleInfo.Type == constant.AlarmRuleTypeMultipleMetric {
		//多指标
		finalExpr, exprDetail, err := generateGroupRuleExpression(alarmItems, ruleInfo.ResourceId, calc_mode.Resource, ruleInfo.Combination)
		if err != nil {
			return nil, errors.NewBusinessError("表达式生成失败")
		}
		desc.RuleItems = alarmItems
		desc.Expr = finalExpr
		desc.ExprDetail = exprDetail
		summary := service.getTemplateLabels(strings.Join(getRuleLabels(alarmItems), ","), calc_mode.Resource)

		return []AlertDTO{{
			RuleType:     "alert",
			Alert:        alertName,
			Record:       "",
			Expr:         finalExpr,
			ForTime:      util.SecToTime(ruleInfo.Times * ruleInfo.Period),
			Summary:      summary,
			Description:  jsonutil.ToString(desc),
			Labels:       getLabels(ruleInfo.RuleName, ruleInfo.SourceType, ruleInfo.Level),
			SilencesTime: getSilenceTime(ruleInfo.SilencesTime),
			SourceType:   ruleInfo.SourceType,
		}}, nil
	}

	//单指标
	list := make([]AlertDTO, len(alarmItems))
	for i, item := range alarmItems {
		expr, exprDetail, err := generateRuleExpression(item, ruleInfo.ResourceId, calc_mode.Resource)
		if err != nil {
			return nil, errors.NewBusinessError("表达式生成失败")
		}

		desc.RuleItems = []model.AlarmItem{item}
		desc.Expr = expr
		desc.ExprDetail = exprDetail

		list[i] = AlertDTO{
			RuleType:     "alert",
			Alert:        alertName,
			Record:       "",
			Expr:         expr,
			ForTime:      util.SecToTime(item.TriggerCondition.Times * item.TriggerCondition.Period),
			Summary:      service.getTemplateLabels(item.TriggerCondition.Labels, calc_mode.Resource),
			Description:  jsonutil.ToString(desc),
			Labels:       getLabels(ruleInfo.RuleName, ruleInfo.SourceType, item.Level),
			SilencesTime: getSilenceTime(item.SilencesTime),
			SourceType:   ruleInfo.SourceType,
		}
	}
	return list, nil

}

func getRuleLabels(alarmItems []model.AlarmItem) []string {
	var labels = make(map[string]string)
	for _, item := range alarmItems {
		for _, l := range strings.Split(item.TriggerCondition.Labels, ",") {
			labels[l] = "1"
		}
	}
	var ss []string
	for k, _ := range labels {
		ss = append(ss, k)
	}
	return ss
}

func getLabels(ruleName string, sourceType uint8, level uint8) map[string]interface{} {
	labelMaps := map[string]interface{}{}
	labelMaps["app"] = ProductLabel
	source := "front"
	if source_type.AutoScaling == sourceType {
		source = "autoScaling"
	}
	labelMaps["ruleName"] = ruleName
	labelMaps["source"] = source
	labelMaps["namespace"] = ProductNamespaceLabel
	labelMaps["severity"] = strconv.Itoa(int(level))
	return labelMaps
}

func getSilenceTime(silencesTime string) string {
	silenceTime, err := strconv.Atoi(silencesTime)
	if err != nil {
		return defaultAlarmInterval
	} else {
		return util.SecToTime(silenceTime - 1)
	}
}

func (service *K8sPrometheusService) buildAlertRuleListByResourceGroup(wg *sync.WaitGroup, resultChan chan []AlertDTO, tenantId string) {
	defer wg.Done()
	var groupRuleList []GroupRuleInfo
	global.DB.Raw(getResourceGroupRuleInfoSql, tenantId).Scan(&groupRuleList)
	var alertList []AlertDTO
	for _, info := range groupRuleList {
		items := dao.AlarmItem.GetItemListByRuleBizId(global.DB, info.RuleId)
		alerts, err := service.buildGroupAlarmInfo(info, items)
		if err != nil {
			logger.Logger().Errorf("build rule err %+v", err)
			continue
		}
		alertList = append(alertList, alerts...)
	}
	resultChan <- alertList
}

func generateRuleExpression(item model.AlarmItem, instanceStr string, calcMode int) (string, string, error) {
	monitorItem := dao.MonitorItem.GetMonitorItemByName(item.TriggerCondition.MetricCode)
	if monitorItem.Id == 0 {
		return "", "", errors.NewBusinessError("指标不存在")
	}
	var expr = ""
	funcName := fmt.Sprintf("%s_over_time", dao.ConfigItem.GetConfigItem(item.TriggerCondition.Statistics, dao.StatisticalMethodsPid, "").Data)
	innerExpr := strings.ReplaceAll(monitorItem.MetricsLinux, constant.MetricLabel, getInstanceLabels(instanceStr, monitorItem.Labels))

	comparison := dao.ConfigItem.GetConfigItem(item.TriggerCondition.ComparisonOperator, dao.ComparisonMethodPid, "").Data
	if calc_mode.Resource == calcMode {
		expr = fmt.Sprintf("%s((%s)[%s:1m])%s%v", funcName, innerExpr, util.SecToTime(item.TriggerCondition.Period), comparison, item.TriggerCondition.Threshold)
	} else {
		fun := dao.ConfigItem.GetConfigItem(item.TriggerCondition.Statistics, dao.StatisticalMethodsPid, "").Data
		expr = fmt.Sprintf("%s((%s)[%s:1m])", funcName, innerExpr, util.SecToTime(item.TriggerCondition.Period))
		expr = fmt.Sprintf("%s(%s)%s%v", fun, expr, comparison, item.TriggerCondition.Threshold)
	}
	return expr, dao.GetExpress2(*item.TriggerCondition), nil
}

func generateGroupRuleExpression(items []model.AlarmItem, instanceStr string, calcMode int, combination uint8) (string, string, error) {
	var exprs = make([]string, len(items))
	var exprDetails = make([]string, len(items))
	for i, item := range items {
		expression, detail, err := generateRuleExpression(item, instanceStr, calcMode)
		if err != nil {
			return "", "", errors.NewBusinessError("生成告警表达式失败")
		}
		exprs[i] = expression
		exprDetails[i] = detail
	}

	if constant.AlarmRuleCombinationAnd == combination {
		return strings.Join(exprs, " and on(instance) "), strings.Join(exprDetails, " 并且 "), nil
	}
	return strings.Join(exprs, " or on(instance) "), strings.Join(exprDetails, " 或者 "), nil
}

func (service *K8sPrometheusService) buildGroupAlarmInfo(ruleInfo GroupRuleInfo, alarmItems []model.AlarmItem) ([]AlertDTO, error) {
	conditionId, err := util.MD5(alarmItems)
	if err != nil {
		return nil, err
	}
	var alertName = fmt.Sprintf("%s#%s#%s", ruleInfo.RuleId, ruleInfo.ResourceGroupId, conditionId)

	var desc = AlarmDescription{
		Rule:            ruleInfo.Parse(),
		ContactGroupIds: strings.Split(ruleInfo.GroupIds, ","),
		ResourceGroupId: ruleInfo.ResourceGroupId,
	}

	instanceList := dao.AlarmRule.GetResourceListByGroup(global.DB, ruleInfo.ResourceGroupId)

	if ruleInfo.Type == constant.AlarmRuleTypeMultipleMetric {
		//多指标
		finalExpr, exprDetail, err := generateGroupRuleExpression(alarmItems, service.joinResourceId(instanceList, "|"), ruleInfo.CalcMode, ruleInfo.Combination)
		if err != nil {
			return nil, err
		}

		desc.RuleItems = alarmItems
		desc.Expr = finalExpr
		desc.ExprDetail = exprDetail

		summary := service.getTemplateLabels(strings.Join(getRuleLabels(alarmItems), ","), calc_mode.Resource)

		return []AlertDTO{{
			RuleType:     "alert",
			Alert:        alertName,
			Record:       "",
			Expr:         finalExpr,
			ForTime:      util.SecToTime(ruleInfo.Times * ruleInfo.Period),
			Summary:      summary,
			Description:  jsonutil.ToString(desc),
			Labels:       getLabels(ruleInfo.RuleName, ruleInfo.SourceType, ruleInfo.Level),
			SilencesTime: getSilenceTime(ruleInfo.SilencesTime),
			SourceType:   ruleInfo.SourceType,
		}}, nil
	}

	//单指标
	list := make([]AlertDTO, len(alarmItems))
	for i, item := range alarmItems {
		expr, exprDetail, err := generateRuleExpression(item, service.joinResourceId(instanceList, "|"), ruleInfo.CalcMode)
		if err != nil {
			return nil, err
		}
		desc.RuleItems = []model.AlarmItem{item}
		desc.Expr = expr
		desc.ExprDetail = exprDetail

		list[i] = AlertDTO{
			RuleType:     "alert",
			Alert:        alertName,
			Record:       "",
			Expr:         expr,
			ForTime:      util.SecToTime(item.TriggerCondition.Times * item.TriggerCondition.Period),
			Summary:      service.getTemplateLabels(item.TriggerCondition.Labels, calc_mode.Resource),
			Description:  jsonutil.ToString(desc),
			Labels:       getLabels(ruleInfo.RuleName, ruleInfo.SourceType, item.Level),
			SilencesTime: getSilenceTime(item.SilencesTime),
			SourceType:   ruleInfo.SourceType,
		}
	}
	return list, nil

}

func getInstanceLabels(instanceId string, labelStr string) string {
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

func buildAlertManagerRouter(alertList []AlertDTO, tenantId string) *AlertManagerConfig {
	var router []Router
	for _, alertDto := range alertList {
		if source_type.Front == alertDto.SourceType {
			continue
		}
		router = append(router, Router{
			Matchers:       map[string]string{"alertname": alertDto.Alert},
			RepeatInterval: alertDto.SilencesTime,
		})
	}
	return &AlertManagerConfig{
		Name:   "tenant-" + tenantId,
		Router: router,
	}
}
