package dao

import (
	commonError "code.cestc.cn/ccos-ops/cloud-monitor/business-common/errors"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/model"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/vo"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type InstanceDao struct {
}

var Instance = new(InstanceDao)

func (dao *InstanceDao) SelectInstanceRulePage(param *form.InstanceRulePageReqParam) interface{} {
	var modelList []*form.InstanceRuleDTO
	var sqlParam = []interface{}{param.InstanceId}
	page := util.Paginate(param.PageSize, param.Current, "SELECT t2.biz_id as id, t2.`name`, t2.metric_name AS monitorItem, t2.trigger_condition AS ruleCondition, product_name, monitor_type, t1.create_time  FROM   t_alarm_rule_resource_rel t1   JOIN t_alarm_rule t2 ON t2.biz_id = t1.alarm_rule_id  WHERE t1.resource_id =?  AND t2.deleted = 0  ORDER BY create_time DESC, NAME ASC", sqlParam, &modelList)
	for _, instanceRuleDTO := range modelList {
		instanceRuleDTO.MonitorItem = instanceRuleDTO.RuleCondition.MonitorItemName
		instanceRuleDTO.Condition = GetExpress(instanceRuleDTO.RuleCondition)
	}
	return &vo.PageVO{
		Records: modelList,
		Current: page.Current,
		Size:    page.Size,
		Total:   page.Total,
		Pages:   page.Pages,
	}
}

func (dao *InstanceDao) UnbindInstance(tx *gorm.DB, param *form.UnBindRuleParam) error {
	exists := AlarmRule.CheckRuleExists(tx, param.RuleId, param.TenantId)
	if !exists {
		logger.Logger().Infof("%s %+v", commonError.RuleNotExist, param)
		return commonError.NewBusinessError(commonError.RuleNotExist)
	}
	tx.Where("resource_id=?", param.InstanceId).Where("alarm_rule_id=?", param.RuleId).Delete(&model.AlarmRuleResourceRel{})
	return nil
}
func (dao *InstanceDao) BindInstance(tx *gorm.DB, param *form.InstanceBindRuleDTO) error {
	tx.Where("resource_id=?", param.InstanceId).Delete(&model.AlarmRuleResourceRel{})
	instance := model.AlarmInstance{
		Ip:           param.Ip,
		RegionCode:   param.RegionCode,
		RegionName:   param.RegionName,
		ZoneCode:     param.ZoneCode,
		ZoneName:     param.ZoneName,
		InstanceName: param.InstanceName,
		InstanceID:   param.InstanceId,
		TenantID:     param.TenantId,
	}
	if len(param.RuleIdList) != 0 {
		var ruleRelList []*model.AlarmRuleResourceRel
		for _, ruleId := range param.RuleIdList {
			exists := AlarmRule.CheckRuleExists(tx, ruleId, param.TenantId)
			if !exists {
				continue
			}
			ruleRelList = append(ruleRelList, &model.AlarmRuleResourceRel{ResourceId: param.InstanceId, AlarmRuleId: ruleId, TenantId: param.TenantId})
		}
		tx.Create(&ruleRelList)
		tx.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "instance_id"}}, DoNothing: false}).Create(&instance)
	}
	return nil
}

func (dao *InstanceDao) GetRuleListByProductType(param *form.ProductRuleParam) *form.ProductRuleListDTO {
	var unbindList []*form.InstanceRuleDTO
	global.DB.Raw("SELECT   biz_id as id,   `name`,   trigger_condition AS ruleCondition,   product_name,   monitor_type   FROM   `t_alarm_rule`    WHERE   product_name =?   AND monitor_type =?   AND tenant_id =?   AND deleted = 0   and source_type=1   AND biz_id NOT IN (   SELECT   alarm_rule_id   FROM   t_alarm_rule_resource_rel    WHERE   resource_id =?   )", param.ProductType, param.MonitorType, param.TenantId, param.InstanceId).Scan(&unbindList)

	var instanceRuleList []*form.InstanceRuleDTO
	db := global.DB
	db.Raw("SELECT    t2.biz_id as id ,    t2.`name`,    t2.trigger_condition AS ruleCondition,    t2.product_name,    t2.monitor_type,    t1.create_time FROM    t_alarm_rule_resource_rel t1,    t_alarm_rule t2 WHERE    t1.resource_id = ? AND t1.alarm_rule_id = t2.biz_id AND t2.deleted = 0 ORDER BY    create_time DESC,    NAME ASC ", param.InstanceId).Scan(&instanceRuleList)
	list := &form.ProductRuleListDTO{
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
