package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/snowflake"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/enum"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global/sys_component/sys_rocketmq"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/mq"
	"gorm.io/gorm"
	"strconv"
)

func CreateRule(tx *gorm.DB, param interface{}) error {
	ruleDao := dao.AlarmRule
	dto := param.(*form.AlarmRuleAddReqDTO)
	dto.Id = strconv.FormatInt(snowflake.GetWorker().NextId(), 10)
	ruleDao.SaveRule(tx, dto)
	return mq.SendMsg(sys_rocketmq.RuleTopic, enum.CreateRule, dto)
}

func UpdateRule(tx *gorm.DB, param interface{}) error {
	ruleDao := dao.AlarmRule
	dto := param.(*form.AlarmRuleAddReqDTO)
	err := ruleDao.UpdateRule(tx, dto)
	if err != nil {
		return err
	}
	return mq.SendMsg(sys_rocketmq.RuleTopic, enum.UpdateRule, dto)
}

func DeleteRule(tx *gorm.DB, param interface{}) error {
	ruleDao := dao.AlarmRule
	dto := param.(*form.RuleReqDTO)
	err := ruleDao.DeleteRule(tx, dto)
	if err != nil {
		return err
	}
	return mq.SendMsg(sys_rocketmq.RuleTopic, enum.DeleteRule, param)
}

func ChangeRuleStatus(tx *gorm.DB, param interface{}) error {
	ruleDao := dao.AlarmRule
	dto := param.(*form.RuleReqDTO)
	err := ruleDao.UpdateRuleState(tx, dto)
	if err != nil {
		return err
	}
	return mq.SendMsg(sys_rocketmq.RuleTopic, enum.ChangeStatus, param)
}
