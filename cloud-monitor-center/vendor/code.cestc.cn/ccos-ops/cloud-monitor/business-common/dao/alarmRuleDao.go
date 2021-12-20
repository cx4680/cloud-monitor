package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/pageUtils"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/vo"
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

type AlarmRuleDao struct {
}

var AlarmRule = new(AlarmRuleDao)

const (
	selectRule = "SELECT t.name as name, " +
		"t.monitor_type, " +
		"t.product_type, " +
		"t.trigger_condition, " +
		"t.status, " +
		"t.metric_name,  " +
		"t.ruleId, " +
		"count(instance) as instanceNum, " +
		"t.update_time " +
		"FROM ( SELECT t1.name, " +
		"t1.monitor_type, " +
		"t1.product_type, " +
		"t1.metric_name, " +
		"t1.trigger_condition, " +
		"t1.enabled AS 'status', " +
		"t1.id AS ruleId, " +
		"t2.instance_id AS instance, " +
		"t1.update_time " +
		"FROM t_alarm_rule t1 " +
		"LEFT JOIN t_alarm_instance t2 ON t2.alarm_rule_id = t1.id " +
		"WHERE t1.tenant_id = ? " +
		"AND t1.deleted = 0 "
)

func (dao *AlarmRuleDao) SaveRule(tx *gorm.DB, ruleReqDTO *forms.AlarmRuleAddReqDTO) string {
	rule := buildAlarmRule(ruleReqDTO)
	rule.MonitorType = ruleReqDTO.MonitorType
	rule.ProductType = ruleReqDTO.ProductType
	tx.Create(rule)
	dao.saveRuleOthers(tx, ruleReqDTO, rule.ID)
	return rule.ID
}
func (dao *AlarmRuleDao) UpdateRule(tx *gorm.DB, ruleReqDTO *forms.AlarmRuleAddReqDTO) {
	dao.deleteOthers(tx, ruleReqDTO.Id)
	rule := buildAlarmRule(ruleReqDTO)
	tx.Model(&rule).Updates(rule)
	dao.saveRuleOthers(tx, ruleReqDTO, ruleReqDTO.Id)
}

func (dao *AlarmRuleDao) DeleteRule(tx *gorm.DB, ruleReqDTO *forms.RuleReqDTO) {
	rule := models.AlarmRule{
		TenantID: ruleReqDTO.TenantId,
		ID:       ruleReqDTO.Id,
	}
	tx.Delete(&rule)
	dao.deleteOthers(tx, ruleReqDTO.Id)
}

func (dao *AlarmRuleDao) UpdateRuleState(tx *gorm.DB, ruleReqDTO *forms.RuleReqDTO) {
	rule := models.AlarmRule{ID: ruleReqDTO.Id}
	tx.Model(&rule).Update("enabled", getAlarmStatusTextInt(ruleReqDTO.Status))
}

func (dao *AlarmRuleDao) SelectRulePageList(param *forms.AlarmPageReqParam) *vo.PageVO {
	var modelList []forms.AlarmRulePageDTO
	selectRuleBuilder := &strings.Builder{}
	var sqlParam = []interface{}{param.TenantId}
	selectRuleBuilder.WriteString(selectRule)
	if len(param.Status) != 0 {
		selectRuleBuilder.WriteString(" AND t1.enabled = ? ")
		sqlParam = append(sqlParam, getAlarmStatusTextInt(param.Status))
	}
	if len(param.RuleName) != 0 {
		selectRuleBuilder.WriteString(" AND t1.name like concat('%',?,'%') ")
		sqlParam = append(sqlParam, param.RuleName)
	}
	selectRuleBuilder.WriteString(") t group by t.ruleId order by t.update_time  desc ")
	page := pageUtils.Paginate(param.PageSize, param.Current, selectRuleBuilder.String(), sqlParam, &modelList)
	for i, v := range modelList {
		modelList[i].MonitorItem = v.RuleCondition.MonitorItemName
		modelList[i].Express = getExpress(v.RuleCondition)
		modelList[i].Status = getAlarmStatusSqlText(v.Status)
	}
	return &vo.PageVO{
		Records: modelList,
		Current: page.Current,
		Size:    page.Size,
		Total:   page.Total,
		Pages:   page.Pages,
	}

}

func (dao *AlarmRuleDao) GetDetail(id string, tenantId string) *forms.AlarmRuleDetailDTO {
	model := &forms.AlarmRuleDetailDTO{}
	global.DB.Raw("SELECT    id ,    NAME  as ruleName,  enabled as status,  product_type,  monitor_type,   level as alarmLevel,  dimensions as scope,  trigger_condition as ruleCondition ,  silences_time,   effective_start,  effective_end,  notify_channel as noticeChannel        FROM t_alarm_rule        WHERE id = ?          AND deleted = 0  and tenant_id=?", id, tenantId).Scan(model)
	model.NoticeGroups = dao.GetNoticeGroupList(id)
	model.InstanceList = dao.GetInstanceList(id)
	return model
}

func (dao *AlarmRuleDao) GetInstanceList(ruleId string) []*forms.InstanceInfo {
	var model []*forms.InstanceInfo
	global.DB.Raw("select instance_id, region_code, zone_code, zone_name, region_name, ip,  instance_name  from t_alarm_instance  where alarm_rule_id =?", ruleId).Scan(&model)
	return model
}

func (dao *AlarmRuleDao) GetNoticeGroupList(ruleId string) []*forms.NoticeGroup {
	var model []*forms.NoticeGroup
	global.DB.Raw("SELECT t1.contract_group_id as id, t2.`name` as name  FROM t_alarm_notice t1,  alert_contact_group t2   WHERE t1.alarm_rule_id = ?   and t1.contract_group_id = t2.id  ORDER BY name", ruleId).Scan(&model)
	for _, group := range model {
		group.UserList = dao.GetUserList(group.Id)
	}
	return model
}

func (dao *AlarmRuleDao) GetUserList(groupId string) []*forms.UserInfo {
	var model []*forms.UserInfo
	global.DB.Raw("select t2.`name` as userName  ,t2.id as id, GROUP_CONCAT(CASE t3.type WHEN 1 THEN t3.no  END) as phone, GROUP_CONCAT(CASE t3.type WHEN 2 THEN t3.no  END) as email from alert_contact_group_rel  t   LEFT JOIN alert_contact t2 on t2.id = t.contact_id   LEFT JOIN alert_contact_information t3 on (t3.contact_id = t2.id and t3.is_certify=1)  where t.group_id=?  and t2.`status`=1  GROUP BY id  order by userName", groupId).Scan(&model)
	return model
}

func (dao *AlarmRuleDao) GetMonitorItem(metricName string) *models.MonitorItem {
	model := &models.MonitorItem{}
	global.DB.Raw("SELECT metrics_linux,metrics_windows,metric_name,unit,name  ,labels FROM monitor_item  where metric_name=? ", metricName).Scan(model)
	return model
}

func buildAlarmRule(ruleReqDTO *forms.AlarmRuleAddReqDTO) *models.AlarmRule {
	return &models.AlarmRule{TenantID: ruleReqDTO.TenantId,
		ID:            ruleReqDTO.Id,
		ProductType:   ruleReqDTO.ProductType,
		Dimensions:    GetResourceScopeInt(ruleReqDTO.Scope),
		Name:          ruleReqDTO.RuleName,
		MetricName:    ruleReqDTO.RuleCondition.MetricName,
		RuleCondition: ruleReqDTO.RuleCondition,
		SilencesTime:  ruleReqDTO.SilencesTime,
		Level:         ruleReqDTO.AlarmLevel,
		NotifyChannel: getNotifyChannel(ruleReqDTO.NoticeChannel),
		CreateUser:    ruleReqDTO.UserId,
	}
}

func (dao *AlarmRuleDao) saveRuleOthers(tx *gorm.DB, ruleReqDTO *forms.AlarmRuleAddReqDTO, ruleId string) {
	dao.saveAlarmNotice(tx, ruleReqDTO, ruleId)
	dao.saveAlarmInstances(tx, ruleReqDTO, ruleId)
}

func (dao *AlarmRuleDao) saveAlarmNotice(tx *gorm.DB, ruleReqDTO *forms.AlarmRuleAddReqDTO, ruleId string) {
	list := make([]models.AlarmNotice, len(ruleReqDTO.GroupList))
	for index, group := range ruleReqDTO.GroupList {
		list[index] = models.AlarmNotice{
			AlarmRuleID:     ruleId,
			ContractGroupID: group,
		}
	}
	tx.Create(&list)
}

func (dao *AlarmRuleDao) saveAlarmInstances(tx *gorm.DB, ruleReqDTO *forms.AlarmRuleAddReqDTO, ruleId string) {
	if len(ruleReqDTO.InstanceList) == 0 {
		return
	}
	list := make([]models.AlarmInstance, len(ruleReqDTO.InstanceList))
	for index, info := range ruleReqDTO.InstanceList {
		list[index] = models.AlarmInstance{
			AlarmRuleID:  ruleId,
			Ip:           info.Ip,
			InstanceID:   info.InstanceId,
			RegionCode:   info.RegionCode,
			ZoneCode:     info.ZoneCode,
			ZoneName:     info.ZoneName,
			RegionName:   info.RegionName,
			InstanceName: info.InstanceName,
			TenantID:     ruleReqDTO.TenantId,
		}
	}
	tx.Create(&list)
}

func (dao *AlarmRuleDao) deleteOthers(tx *gorm.DB, ruleId string) {
	notice := models.AlarmNotice{
		AlarmRuleID: ruleId,
	}
	tx.Where("alarm_rule_id=?", ruleId).Delete(&notice)
	instance := models.AlarmInstance{AlarmRuleID: ruleId}
	tx.Where("alarm_rule_id=?", ruleId).Delete(&instance)
}

func getNotifyChannel(notifyChannel string) int {
	notify, _ := strconv.Atoi(ConfigItem.GetConfigItem(nil, "33", notifyChannel).Code)
	return notify
}

const (
	ALL      = "ALL"
	INSTANCE = "INSTANCE"
)

var ResourceScopeText = map[string]int{
	ALL:      1,
	INSTANCE: 2,
}

func GetResourceScopeInt(code string) int {
	return ResourceScopeText[code]
}

const (
	ENABLE  = "enabled"
	DISABLE = "disabled"

	sqlEnabled  = "1"
	sqlDisabled = "2"

	Maximum = "Maximum"
	Minimum = "Minimum"
	Average = "Average"

	Greater        = "greater"
	GreaterOrEqual = "greaterOrEqual"
	Less           = "less"
	lessOrEqual    = "lessOrEqual"
	Equal          = "equal"
	NotEqual       = "notEqual"
)

var alarmStatusText = map[string]int{
	ENABLE:  1,
	DISABLE: 0,
}

func getAlarmStatusTextInt(code string) int {
	return alarmStatusText[code]
}

var alarmStatusSqlText = map[string]string{
	sqlEnabled:  "enabled",
	sqlDisabled: "disabled",
}

func getAlarmStatusSqlText(code string) string {
	return alarmStatusSqlText[code]
}

var alarmStatisticsText = map[string]string{
	Maximum: "最大值",
	Minimum: "最小值",
	Average: "平均值",
}

func getAlarmStatisticsText(s string) string {
	return alarmStatisticsText[s]
}

var comparisonOperatorText = map[string]string{
	Greater:        ">",
	GreaterOrEqual: ">=",
	Less:           "<",
	lessOrEqual:    "<=",
	Equal:          "==",
	NotEqual:       "!=",
}

func getComparisonOperator(s string) string {
	return comparisonOperatorText[s]
}

func getExpress(form *forms.RuleCondition) string {
	return fmt.Sprintf("%s%s%s%s%s 统计周期%s分钟 持续%s个周期", form.MonitorItemName, getAlarmStatisticsText(form.Statistics), getComparisonOperator(form.ComparisonOperator), strconv.FormatFloat(form.Threshold, 'g', 5, 32), form.Unit, strconv.Itoa(form.Period/60), strconv.Itoa(form.Times))
}
