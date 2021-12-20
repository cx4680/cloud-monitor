package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/enums"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysComponent/sysRocketMq"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/mq"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/utils/snowflake"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

func CreateRule(tx *gorm.DB, param interface{}) error {
	ruleDao := dao.AlarmRule
	dto := param.(*forms.AlarmRuleAddReqDTO)
	dto.Id = strconv.FormatInt(snowflake.GetWorker().NextId(), 10)
	ruleDao.SaveRule(tx, dto)
	return mq.SendMsg(sysRocketMq.RuleTopic, enums.CreateRule, dto)
}

func UpdateRule(tx *gorm.DB, param interface{}) error {
	ruleDao := dao.AlarmRule
	dto := param.(*forms.AlarmRuleAddReqDTO)
	ruleDao.UpdateRule(tx, dto)
	return mq.SendMsg(sysRocketMq.RuleTopic, enums.UpdateRule, dto)
}

func DeleteRule(tx *gorm.DB, param interface{}) error {
	ruleDao := dao.AlarmRule
	dto := param.(*forms.RuleReqDTO)
	ruleDao.DeleteRule(tx, dto)
	return mq.SendMsg(sysRocketMq.RuleTopic, enums.DeleteRule, param)
}

func ChangeRuleStatus(tx *gorm.DB, param interface{}) error {
	ruleDao := dao.AlarmRule
	dto := param.(*forms.RuleReqDTO)
	ruleDao.UpdateRuleState(tx, dto)
	enum := enums.DisableRule
	if strings.EqualFold(dto.Status, dao.ENABLE) {
		enum = enums.EnableRule
	}
	return mq.SendMsg(sysRocketMq.RuleTopic, enum, param)
}
