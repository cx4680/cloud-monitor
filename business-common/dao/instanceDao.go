package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/vo"
	"gorm.io/gorm"
)

type InstanceDao struct {
	db *gorm.DB
}

func NewInstanceDao(db *gorm.DB) *InstanceDao {
	return &InstanceDao{db: db}
}

func (dao *InstanceDao) SelectInstanceRulePage(param *forms.InstanceRulePageReqParam) interface{} {
	var model []forms.InstanceRuleDTO
	db := dao.db
	db.Raw(" select t2.id,t2.`name`,t2.trigger_condition  as ruleCondition,product_type,monitor_type ,t1.create_time from t_alarm_instance  t1        JOIN t_alarm_rule t2  on t2.id=t1.alarm_rule_id       where t1.instance_id=? and t2.deleted=0  ORDER BY create_time desc  , name ASC", param.InstanceId)
	db.Find(&model)
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

func (dao *InstanceDao) UnbindInstance(param *forms.UnBindRuleParam) {
	model := &models.AlarmInstance{InstanceId: param.InstanceId, AlarmRuleId: param.RulId}
	dao.db.Delete(model)
}
func (dao *InstanceDao) BindInstance(param *forms.InstanceBindRuleDTO) {
	model := &models.AlarmInstance{InstanceId: param.InstanceId}
	dao.db.Delete(model)
	if len(param.RuleIdList) != 0 {
		list := make([]*models.AlarmInstance, len(param.RuleIdList))
		for index, ruleId := range param.RuleIdList {
			list[index] = &models.AlarmInstance{
				AlarmRuleId:  ruleId,
				Ip:           param.Ip,
				RegionCode:   param.RegionCode,
				RegionName:   param.RegionName,
				ZoneCode:     param.ZoneCode,
				ZoneName:     param.ZoneName,
				InstanceName: param.InstanceName,
				InstanceId:   param.InstanceId,
				TenantId:     param.TenantId,
			}
		}
		dao.db.Create(&list)
	}
}

func (dao *InstanceDao) GetRuleListByProductType(param *forms.ProductRuleParam) *forms.ProductRuleListDTO {
	var unbindList []forms.InstanceRuleDTO
	dao.db.Raw("SELECT\n\tid,\n\t`name`,\n\ttrigger_condition AS ruleCondition,\n\tproduct_type,\n\tmonitor_type\nFROM\n\t`t_alarm_rule`\nWHERE\n\tproduct_type =?\nAND monitor_type =?\nAND tenant_id =?\nAND deleted = 0\nAND id NOT IN (\n\tSELECT\n\t\tt2.id\n\tFROM\n\t\tt_alarm_instance t1\n\tJOIN t_alarm_rule t2 ON t2.id = t1.alarm_rule_id\n\tWHERE\n\t\tt1.instance_id =?\n\tAND t2.deleted = 0\n)", param.ProductType, param.MonitorType, param.TenantId, param.InstanceId)
	dao.db.Find(&unbindList)

	var instanceRuleList []forms.InstanceRuleDTO
	db := dao.db
	db.Raw(" select t2.id,t2.`name`,t2.trigger_condition  as ruleCondition,product_type,monitor_type ,t1.create_time from t_alarm_instance  t1        JOIN t_alarm_rule t2  on t2.id=t1.alarm_rule_id       where t1.instance_id=? and t2.deleted=0  ORDER BY create_time desc  , name ASC", param.InstanceId)
	db.Find(&instanceRuleList)

	return &forms.ProductRuleListDTO{
		BindRuleList:   instanceRuleList,
		UnbindRuleList: unbindList,
	}
}
