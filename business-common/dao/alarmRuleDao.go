package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/vo"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/utils"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/utils/snowflake"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

type AlarmRuleDao struct {
	db *gorm.DB
}

func NewAlarmRuleDao(db *gorm.DB) *AlarmRuleDao {
	return &AlarmRuleDao{
		db: db,
	}
}

func (dao *AlarmRuleDao) SaveRule(ruleReqDTO *forms.AlarmRuleAddReqDTO) string {
	rule := buildAlarmRule(ruleReqDTO)
	rule.ID = strconv.FormatInt(snowflake.GetWorker().NextId(), 10)
	rule.MonitorType = ruleReqDTO.MonitorType
	rule.ProductType = ruleReqDTO.ProductType
	dao.db.Create(rule)
	dao.saveRuleOthers(ruleReqDTO, rule.ID)
	return rule.ID
}
func (dao *AlarmRuleDao) UpdateRule(ruleReqDTO *forms.AlarmRuleAddReqDTO) {
	dao.deleteOthers(ruleReqDTO.Id)
	rule := buildAlarmRule(ruleReqDTO)
	//.db.Update(rule)
	dao.saveRuleOthers(ruleReqDTO, rule.ID)
}

func (dao *AlarmRuleDao) DeleteRule(ruleReqDTO *forms.RuleReqDTO) {
	rule := models.AlarmRule{
		TenantID: ruleReqDTO.TenantId,
		ID:       ruleReqDTO.Id,
	}
	dao.db.Delete(&rule)
	dao.deleteOthers(ruleReqDTO.Id)
}

func (dao *AlarmRuleDao) UpdateRuleState(ruleReqDTO *forms.RuleReqDTO) {
	rule := models.AlarmRule{ID: ruleReqDTO.Id, Enabled: GetAlarmStatusTextInt(ruleReqDTO.Status), TenantID: ruleReqDTO.TenantId}
	dao.db.Model(&rule).Updates(&rule)
}

func (dao *AlarmRuleDao) SelectRulePageList(param *forms.AlarmPageReqParam) interface{} {
	var model []forms.AlarmRulePageDTO
	db := dao.db
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
	var total int64
	db.Offset((param.Current-1)*param.PageSize).Limit(param.PageSize).Raw(selectList.String(), sqlParam).Scan(&model).Count(&total)
	var page = &vo.PageVO{
		Records: model,
		Current: param.Current,
		Size:    param.PageSize,
		Total:   utils.Int64ToInt(total),
	}
	return page
}

func (dao *AlarmRuleDao) GetDetail(id string, tenantId string) *forms.AlarmRuleDetailDTO {
	model := &forms.AlarmRuleDetailDTO{}
	dao.db.Raw("SELECT    id ,    NAME  as ruleName,  enabled as status,  product_type,  monitor_type,   level as alarmLevel,  dimensions as scope,  trigger_condition as ruleCondition ,  silences_time,   effective_start,  effective_end,  notify_channel as noticeChannel        FROM t_alarm_rule        WHERE id = ?          AND deleted = 0  and tenant_id=?", id, tenantId).Scan(model)
	return model
}

func (dao *AlarmRuleDao) GetRuleCondition(id string, tenantId string) {

}

func (dao *AlarmRuleDao) GetMonitorItem(metricName string) *models.MonitorItem {
	model := &models.MonitorItem{}
	dao.db.Raw("SELECT metrics_linux,metrics_windows,metric_name,unit,name  ,labels FROM monitor_item  where metric_name=? ", metricName).Scan(model)
	return model
}

func buildAlarmRule(ruleReqDTO *forms.AlarmRuleAddReqDTO) *models.AlarmRule {
	return &models.AlarmRule{TenantID: ruleReqDTO.TenantId,
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

func (dao *AlarmRuleDao) saveRuleOthers(ruleReqDTO *forms.AlarmRuleAddReqDTO, ruleId string) {
	dao.saveAlarmNotice(ruleReqDTO, ruleId)
	dao.saveAlarmInstances(ruleReqDTO, ruleId)
}

func (dao *AlarmRuleDao) saveAlarmNotice(ruleReqDTO *forms.AlarmRuleAddReqDTO, ruleId string) {
	list := make([]models.AlarmNotice, len(ruleReqDTO.GroupList))
	for index, group := range ruleReqDTO.GroupList {
		list[index] = models.AlarmNotice{
			AlarmRuleID:     ruleId,
			ContractGroupID: group,
		}
	}
	dao.db.Create(&list)
}

func (dao *AlarmRuleDao) saveAlarmInstances(ruleReqDTO *forms.AlarmRuleAddReqDTO, ruleId string) {
	if len(ruleReqDTO.InstanceList) == 0 {
		return
	}
	list := make([]models.AlarmInstance, len(ruleReqDTO.InstanceList))
	for index, info := range ruleReqDTO.InstanceList {
		list[index] = models.AlarmInstance{
			AlarmRuleID:  ruleId,
			IP:           info.Ip,
			InstanceID:   info.InstanceId,
			RegionCode:   info.RegionCode,
			ZoneCode:     info.ZoneCode,
			ZoneName:     info.ZoneName,
			RegionName:   info.RegionName,
			InstanceName: info.InstanceName,
			TenantID:     ruleReqDTO.TenantId,
		}
	}
	dao.db.Create(&list)
}

func (dao *AlarmRuleDao) deleteOthers(ruleId string) {
	notice := models.AlarmNotice{
		AlarmRuleID: ruleId,
	}
	dao.db.Delete(&notice)
	instance := models.AlarmInstance{AlarmRuleID: ruleId}
	dao.db.Delete(&instance)
}

////todo 查询通知方式
func getNotifyChannel(notifyChannel string) int {
	return 1
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
