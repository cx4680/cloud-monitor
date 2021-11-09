package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/pageUtils"
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
	var sqlParam = []interface{}{param.InstanceId}
	return pageUtils.Paginate(param.PageSize, param.Current, "select t2.id,t2.name,t2.metric_name as monitorItem,t2.trigger_condition  as ruleCondition,product_type,monitor_type ,t1.create_time from t_alarm_instance  t1        JOIN t_alarm_rule t2  on t2.id=t1.alarm_rule_id       where t1.instance_id=? and t2.deleted=0  ORDER BY create_time desc  , name ASC", sqlParam, &model, db)
}

func (dao *InstanceDao) UnbindInstance(param *forms.UnBindRuleParam) {
	model := &models.AlarmInstance{InstanceID: param.InstanceId, AlarmRuleID: param.RulId}
	dao.db.Where("instance_id", model.InstanceID).Where("alarm_rule_id", model.AlarmRuleID).Delete(model)
}
func (dao *InstanceDao) BindInstance(param *forms.InstanceBindRuleDTO) {
	model := &models.AlarmInstance{InstanceID: param.InstanceId}
	dao.db.Delete(model)
	if len(param.RuleIdList) != 0 {
		list := make([]*models.AlarmInstance, len(param.RuleIdList))
		for index, ruleId := range param.RuleIdList {
			list[index] = &models.AlarmInstance{
				AlarmRuleID:  ruleId,
				Ip:           param.Ip,
				RegionCode:   param.RegionCode,
				RegionName:   param.RegionName,
				ZoneCode:     param.ZoneCode,
				ZoneName:     param.ZoneName,
				InstanceName: param.InstanceName,
				InstanceID:   param.InstanceId,
				TenantID:     param.TenantId,
			}
		}
		dao.db.Create(&list)
	}
}

func (dao *InstanceDao) GetRuleListByProductType(param *forms.ProductRuleParam) *forms.ProductRuleListDTO {
	var unbindList []forms.InstanceRuleDTO
	dao.db.Raw("SELECT id,`name`,trigger_condition AS ruleCondition,product_type,monitor_type FROM`t_alarm_rule` WHERE product_type =? AND monitor_type =? AND tenant_id =? AND deleted = 0 AND id NOT IN ( SELECT  t2.id FROM  t_alarm_instance t1 JOIN t_alarm_rule t2 ON t2.id = t1.alarm_rule_id WHERE  t1.instance_id =? AND t2.deleted = 0 )", param.ProductType, param.MonitorType, param.TenantId, param.InstanceId).Scan(&unbindList)

	var instanceRuleList []forms.InstanceRuleDTO
	db := dao.db
	db.Raw(" select t2.id,t2.`name`,t2.trigger_condition  as ruleCondition,product_type,monitor_type ,t1.create_time from t_alarm_instance  t1        JOIN t_alarm_rule t2  on t2.id=t1.alarm_rule_id       where t1.instance_id=? and t2.deleted=0  ORDER BY create_time desc  , name ASC", param.InstanceId).Scan(&instanceRuleList)
	return &forms.ProductRuleListDTO{
		BindRuleList:   instanceRuleList,
		UnbindRuleList: unbindList,
	}
}
