package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/model"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/constant"
	"gorm.io/gorm"
)

type AlarmRuleTemplateDao struct{}

var AlarmRuleTemplate = new(AlarmRuleTemplateDao)

type TemplateProduct struct {
	ProductBizId string `gorm:"column:productBizId" json:"productBizId"`
	ProductName  string `gorm:"column:productName" json:"name"`
	Description  string `gorm:"column:description" json:"description"`
	OpenStatus   uint64 `gorm:"column:openStatus" json:"openStatus"`
}

func (dao *AlarmRuleTemplateDao) QueryTemplateProductList(db *gorm.DB, tenantId string) []TemplateProduct {
	sql := "select p.`name` productName, p.description, p.biz_id productBizId,   IF(count(rel.id) > 0, 1, 0) openStatus\n " +
		"from t_alarm_rule_template t \n" +
		"join t_monitor_product p on t.product_biz_id=p.biz_id\n" +
		"left join t_tenant_alarm_template_rel rel on t.biz_id=rel.template_biz_id and rel.tenant_id=?\n" +
		"group by p.`name`, p.description, p.biz_id\n" +
		"order by p.biz_id"
	var list []TemplateProduct
	db.Raw(sql, tenantId).Scan(&list)
	return list
}

type AlarmRuleTemplateRule struct {
	RuleId         string   `json:"ruleId" gorm:"column:ruleId"`
	RuleTemplateId string   `json:"ruleTemplateId" gorm:"column:ruleTemplateId"`
	RuleName       string   `json:"ruleName" gorm:"column:ruleName"`
	Enabled        int      `json:"enabled" gorm:"column:enabled"`
	GroupNames     string   `json:"groupNames" gorm:"column:groupNames"`
	Conditions     []string `json:"conditions" gorm:"-"`
	Type           int      `json:"type" gorm:"type"`
}

func (dao *AlarmRuleTemplateDao) QueryRuleTemplateListByProduct(db *gorm.DB, tenantId, productBizId string) []AlarmRuleTemplateRule {
	templateBizIds := dao.QueryTemplateBizIdListByProductBizId(db, productBizId)
	var relList []model.TenantAlarmTemplateRel
	db.Model(&model.TenantAlarmTemplateRel{}).Where("tenant_id=? and template_biz_id in (?)", tenantId, templateBizIds).Scan(&relList)

	var params []interface{}
	sql := "select t.name ruleName, t.biz_id ruleTemplateId, "

	if len(relList) == 0 {
		sql += "0 type, t.biz_id ruleId, 0 enabled, '" + constant.DefaultContact + "' groupNames FROM t_alarm_rule_template t "
	} else {
		sql += "1 type, r.biz_id ruleId, r.enabled, GROUP_CONCAT(g.`name`) groupNames FROM t_alarm_rule_template t\n" +
			" join t_alarm_rule r on t.biz_id=r.template_biz_id and r.deleted=0 and r.tenant_id=?\n" +
			"left join t_alarm_notice n on r.biz_id=n.alarm_rule_id\n" +
			"left join t_contact_group g on n.contract_group_id=g.biz_id "

		params = append(params, tenantId)
	}
	sql += " where t.product_biz_id=?\ngroup by t.id "
	if len(relList) > 0 {
		sql += ",r.biz_id"
	}
	params = append(params, productBizId)

	var list []AlarmRuleTemplateRule
	db.Raw(sql, params...).Scan(&list)
	return list
}

func (dao *AlarmRuleTemplateDao) QueryTemplateBizIdListByProductBizId(db *gorm.DB, productBizId string) []string {
	var templateBizIds []string
	db.Raw("select biz_id from t_alarm_rule_template where product_biz_id=?", productBizId).Scan(&templateBizIds)
	return templateBizIds
}

func (dao *AlarmRuleTemplateDao) QueryByBizId(db *gorm.DB, bizId string) model.AlarmRuleTemplate {
	var record model.AlarmRuleTemplate
	db.Where("biz_id=?", bizId).Scan(&record)
	return record
}

func (dao *AlarmRuleTemplateDao) QueryCreateRuleInfo(db *gorm.DB, templateBizId string) []form.AlarmRuleAddReqDTO {
	sql := "select \nt.monitor_type,p.`name` product_type,t.product_biz_id product_id,'ALL' scope,t.`level`,t.`name` rule_name,t.silences_time,t.type,t.period,t.times,t.combination,t.effective_start,t.effective_end,t.metric_code, t.biz_id template_biz_id \n " +
		"FROM\nt_alarm_rule_template t \njoin t_monitor_product p on t.product_biz_id=p.biz_id\n where t.product_biz_id=?"
	var d []form.AlarmRuleAddReqDTO
	db.Raw(sql, templateBizId).Scan(&d)
	return d
}
