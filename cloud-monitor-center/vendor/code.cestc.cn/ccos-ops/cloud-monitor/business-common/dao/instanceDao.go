package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/pageUtils"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/vo"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type InstanceDao struct {
}

var Instance = new(InstanceDao)

func (dao *InstanceDao) SelectInstanceRulePage(param *forms.InstanceRulePageReqParam) interface{} {
	var modelList []*forms.InstanceRuleDTO
	var sqlParam = []interface{}{param.InstanceId}
	page := pageUtils.Paginate(param.PageSize, param.Current, "SELECT t2.id, t2.`name`, t2.metric_name AS monitorItem, t2.trigger_condition AS ruleCondition, product_type, monitor_type, t1.create_time  FROM   t_alarm_rule_resource_rel t1   JOIN t_alarm_rule t2 ON t2.id = t1.alarm_rule_id  WHERE t1.resource_id =?  AND t2.deleted = 0  ORDER BY create_time DESC, NAME ASC", sqlParam, &modelList)
	for _, model := range modelList {
		model.MonitorItem = model.RuleCondition.MonitorItemName
		model.Condition = GetExpress(model.RuleCondition)
	}
	return &vo.PageVO{
		Records: modelList,
		Current: page.Current,
		Size:    page.Size,
		Total:   page.Total,
		Pages:   page.Pages,
	}
}

func (dao *InstanceDao) UnbindInstance(tx *gorm.DB, param *forms.UnBindRuleParam) {
	tx.Where("resource_id=?", param.InstanceId).Where("alarm_rule_id=?", param.RuleId).Delete(&models.AlarmRuleResourceRel{})
}
func (dao *InstanceDao) BindInstance(tx *gorm.DB, param *forms.InstanceBindRuleDTO) {
	tx.Where("resource_id=?", param.InstanceId).Delete(&models.AlarmRuleResourceRel{})
	if len(param.RuleIdList) != 0 {
		instanceList := make([]*models.AlarmInstance, len(param.RuleIdList))
		ruleRelList := make([]*models.AlarmRuleResourceRel, len(param.RuleIdList))
		for index, ruleId := range param.RuleIdList {
			instanceList[index] = &models.AlarmInstance{
				Ip:           param.Ip,
				RegionCode:   param.RegionCode,
				RegionName:   param.RegionName,
				ZoneCode:     param.ZoneCode,
				ZoneName:     param.ZoneName,
				InstanceName: param.InstanceName,
				InstanceID:   param.InstanceId,
				TenantID:     param.TenantId,
			}
			ruleRelList[index] = &models.AlarmRuleResourceRel{ResourceId: param.InstanceId, AlarmRuleId: ruleId, TenantId: param.TenantId}
		}
		tx.Create(&ruleRelList)
		tx.Clauses(clause.OnConflict{DoNothing: false}).Create(&instanceList)
	}
}

func (dao *InstanceDao) GetRuleListByProductType(param *forms.ProductRuleParam) *forms.ProductRuleListDTO {
	var unbindList []*forms.InstanceRuleDTO
	global.DB.Raw("SELECT   id,   `name`,   trigger_condition AS ruleCondition,   product_type,   monitor_type   FROM   `t_alarm_rule`    WHERE   product_type =?   AND monitor_type =?   AND tenant_id =?   AND deleted = 0   and source_type=1   AND id NOT IN (   SELECT   alarm_rule_id   FROM   t_alarm_rule_resource_rel    WHERE   resource_id =?   )", param.ProductType, param.MonitorType, param.TenantId, param.InstanceId).Scan(&unbindList)

	var instanceRuleList []*forms.InstanceRuleDTO
	db := global.DB
	db.Raw("SELECT    t2.id,    t2.`name`,    t2.trigger_condition AS ruleCondition,    t2.product_type,    t2.monitor_type,    t1.create_time FROM    t_alarm_rule_resource_rel t1,    t_alarm_rule t2 WHERE    t1.resource_id = ? AND t1.alarm_rule_id = t2.id AND t2.deleted = 0 ORDER BY    create_time DESC,    NAME ASC ", param.InstanceId).Scan(&instanceRuleList)
	list := &forms.ProductRuleListDTO{
		BindRuleList:   instanceRuleList,
		UnbindRuleList: unbindList,
	}
	for _, dto := range list.BindRuleList {
		dto.MonitorItem = dto.RuleCondition.MonitorItemName
		dto.Condition = GetExpress(dto.RuleCondition)
	}

	for _, dto := range list.UnbindRuleList {
		dto.MonitorItem = dto.RuleCondition.MonitorItemName
		dto.Condition = GetExpress(dto.RuleCondition)
	}
	return list
}
