package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/pageUtils"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/vo"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/utils/snowflake"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

type AlarmRuleDao struct {
}

var AlarmRule = new(AlarmRuleDao)

func (dao *AlarmRuleDao) SaveRule(tx *gorm.DB, ruleReqDTO *forms.AlarmRuleAddReqDTO) string {
	rule := buildAlarmRule(ruleReqDTO)
	rule.ID = strconv.FormatInt(snowflake.GetWorker().NextId(), 10)
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
	tx.Model(&rule).Update("enabled", GetAlarmStatusTextInt(ruleReqDTO.Status))
}

func (dao *AlarmRuleDao) SelectRulePageList(param *forms.AlarmPageReqParam) *vo.PageVO {
	var model []forms.AlarmRulePageDTO
	selectList := &strings.Builder{}
	var sqlParam = []interface{}{param.TenantId}
	selectList.WriteString("SELECT name as name,monitor_type, product_type, trigger_condition,  status,  metric_name,  ruleId,  count(instance) as instanceNum, update_time       FROM (  SELECT NAME,   monitor_type,   product_type,  metric_name,  trigger_condition,    enabled AS 'status',      id     AS ruleId,    t2.instance_id AS instance,   t1.update_time   FROM t_alarm_rule t1    LEFT JOIN t_alarm_instance t2 ON t2.alarm_rule_id = t1.id  WHERE t1.tenant_id = ?    AND t1.deleted = 0")
	if len(param.Status) != 0 {
		selectList.WriteString("t1.enabled = ?")
		sqlParam = append(sqlParam, param.Status)
	}
	if len(param.RuleName) != 0 {
		selectList.WriteString("t1.name like concat('%',?,'%')")
		sqlParam = append(sqlParam, param.RuleName)
	}
	selectList.WriteString(") t group by t.ruleId order by t.update_time  desc ")
	return pageUtils.Paginate(param.PageSize, param.Current, selectList.String(), sqlParam, &model)

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
)

var AlarmStatusText = map[string]int{
	ENABLE:  1,
	DISABLE: 0,
}

func GetAlarmStatusTextInt(code string) int {
	return AlarmStatusText[code]
}
