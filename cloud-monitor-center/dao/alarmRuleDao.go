package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor-center/database"
	"code.cestc.cn/ccos-ops/cloud-monitor-center/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor-center/global"
	"code.cestc.cn/ccos-ops/cloud-monitor-center/models"
	"code.cestc.cn/ccos-ops/cloud-monitor-center/utils/snowflake"
	"code.cestc.cn/ccos-ops/cloud-monitor-center/vo"
	"github.com/jinzhu/gorm"
	"strconv"
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
	rule.Id = strconv.FormatInt(snowflake.GetWorker().NextId(), 10)
	rule.MonitorType = ruleReqDTO.MonitorType
	rule.ProductType = ruleReqDTO.ProductType
	dao.db.Create(rule)
	return rule.Id
}

func buildAlarmRule(ruleReqDTO *forms.AlarmRuleAddReqDTO) *models.AlarmRule {
	return &models.AlarmRule{TenantId: ruleReqDTO.TenantId,
		ProductType:   ruleReqDTO.ProductType,
		Dimensions:    global.GetResourceScopeInt(ruleReqDTO.Scope),
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
	list := make([]*models.AlarmNotice, len(ruleReqDTO.InstanceList))
	for index, group := range ruleReqDTO.GroupList {
		list[index] = &models.AlarmNotice{
			AlarmRuleId:     ruleId,
			ContractGroupId: group,
		}
	}
	// todo 批量插入 dao.db
}

func (dao *AlarmRuleDao) saveAlarmInstances(ruleReqDTO *forms.AlarmRuleAddReqDTO, ruleId string) {
	list := make([]*models.AlarmInstance, len(ruleReqDTO.InstanceList))
	for index, info := range ruleReqDTO.InstanceList {
		instance := &models.AlarmInstance{
			AlarmRuleId: ruleId,
			Ip:          info.Ip,
		}
		list[index] = instance
	}
}

func (dao *AlarmRuleDao) UpdateRule(ruleReqDTO *forms.AlarmRuleAddReqDTO) {

}

func (dao *AlarmRuleDao) DeleteRule(ruleReqDTO *forms.AlarmRuleAddReqDTO) {

}

func (dao *AlarmRuleDao) SelectRulePageList(param *forms.AlarmPageReqParam) interface{} {
	var model []forms.AlarmRulePageDTO
	db := dao.db
	db.Raw(" SELECT NAME,monitor_type, product_type, trigger_condition as ruleCondition,  status,  metric_name,  ruleId,  count(instance) as instanceNum, update_time       FROM (  SELECT NAME,   monitor_type,   product_type,  metric_name,  trigger_condition,    enabled AS 'status',      id     AS ruleId,    t2.instance_id AS instance,   t1.update_time   FROM t_alarm_rule t1    LEFT JOIN t_alarm_instance t2 ON t2.alarm_rule_id = t1.id  WHERE t1.tenant_id = ?    AND t1.deleted = 0", param.TenantId)
	if param.Status != "" {
		db.Where("t1.enabled = ?", param.Status)
	}
	if param.RuleName != "" {
		db.Where("t1.name like concat('%',?,'%')", param.RuleName)
	}
	db.Group(" t.ruleId ").Order(" t.update_time  desc ")
	db.Find(model)
	total := len(model)
	db.Limit(param.PageSize).Offset((param.Current - 1) * param.PageSize).Find(model)
	var page = &vo.PageVO{
		Records: model,
		Current: param.Current,
		Size:    param.PageSize,
		Total:   total,
	}
	return page
}

func (dao *AlarmRuleDao) GetDetail(id string, tenantId string) *forms.AlarmRuleDetailDTO {
	model := &forms.AlarmRuleDetailDTO{}
	dao.db.Raw(database.SelectRuleDetail, id, tenantId).Scan(model)
	return model
}

func (dao *AlarmRuleDao) GetRuleCondition(id string, tenantId string) {

}

func (dao *AlarmRuleDao) SelectInstanceRulePage(instanceId string) *forms.InstanceRuleDTO {
	return nil
}

func (dao *AlarmRuleDao) UnbindInstance(instanceId string, rulId string) *forms.InstanceRuleDTO {
	return nil
}
func (dao *AlarmRuleDao) BindInstance(bindRuleParam *forms.InstanceBindRuleDTO) *forms.InstanceRuleDTO {
	return nil
}

func (dao *AlarmRuleDao) GetRuleListByProductType(bindRuleParam *forms.InstanceBindRuleDTO) *forms.InstanceRuleDTO {
	return nil
}

func (dao *AlarmRuleDao) GetMonitorItem(metricName string) *models.MonitorItem {
	model := &models.MonitorItem{}
	dao.db.Raw("SELECT metrics_linux,metrics_windows,metric_name,unit,name  ,labels FROM monitor_item ").Where("where metric_name=?", metricName).Scan(model)
	return model
}

////todo 查询通知方式
func getNotifyChannel(notifyChannel string) int {
	return 1
}
