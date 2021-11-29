package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/mq"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/enums"
	"gorm.io/gorm"
	"strings"
)

func CreateRule(tx *gorm.DB, param interface{}) error {
	ruleDao := dao.AlarmRule
	dto := param.(*forms.AlarmRuleAddReqDTO)
	id := ruleDao.SaveRule(tx, dto)
	dto.Id = id
	err := mq.SendMsg(config.GetRocketmqConfig().RuleTopic, enums.CreateRule, dto)
	return err
}

func UpdateRule(tx *gorm.DB, param interface{}) error {
	ruleDao := dao.AlarmRule
	dto := param.(*forms.AlarmRuleAddReqDTO)
	ruleDao.UpdateRule(tx, dto)
	err := mq.SendMsg(config.GetRocketmqConfig().RuleTopic, enums.UpdateRule, dto)
	return err
}

func DeleteRule(tx *gorm.DB, param interface{}) error {
	ruleDao := dao.AlarmRule
	dto := param.(*forms.RuleReqDTO)
	ruleDao.DeleteRule(tx, dto)
	err := mq.SendMsg(config.GetRocketmqConfig().RuleTopic, enums.DeleteRule, param)
	return err
}

func ChangeRuleStatus(tx *gorm.DB, param interface{}) error {
	ruleDao := dao.AlarmRule
	dto := param.(*forms.RuleReqDTO)
	ruleDao.UpdateRuleState(tx, dto)
	enum := enums.DisableRule
	if strings.EqualFold(dto.Status, dao.ENABLE) {
		enum = enums.EnableRule
	}
	err := mq.SendMsg(config.GetRocketmqConfig().RuleTopic, enum, param)
	return err
}
