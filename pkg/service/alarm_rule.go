package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/snowflake"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/enum"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/errors"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global/sys_component/sys_rocketmq"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/constant"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/mq"
	"gorm.io/gorm"
	"strconv"
)

func CreateRule(tx *gorm.DB, param interface{}) error {
	dto := param.(*form.AlarmRuleAddReqDTO)
	if err := checkConditions(dto); err != nil {
		return errors.NewBusinessError(err.Error())
	}
	if dao.AlarmRule.CheckRuleName(dto.TenantId, dto.RuleName, "") {
		return errors.NewBusinessError("规则名称不能重复")
	}
	for i := range dto.GroupList {
		if dto.GroupList[i] == "-1" {
			contactService := NewContactService(NewContactGroupService(NewContactGroupRelService()), NewContactInformationService(nil), NewContactGroupRelService())
			groupBizId, err := contactService.CreateSysContact(dto.TenantId)
			if err != nil {
				return err
			}
			dto.GroupList[i] = groupBizId
		}
	}
	dto.Id = strconv.FormatInt(snowflake.GetWorker().NextId(), 10)
	dao.AlarmRule.SaveRule(tx, dto)
	//其他Region同步
	return mq.SendMsg(sys_rocketmq.RuleTopic, enum.CreateRule, dto)
}

func UpdateRule(tx *gorm.DB, param interface{}) error {
	dto := param.(*form.AlarmRuleAddReqDTO)
	if err := checkConditions(dto); err != nil {
		return errors.NewBusinessError(err.Error())
	}
	if dao.AlarmRule.CheckRuleName(dto.TenantId, dto.RuleName, dto.Id) {
		return errors.NewBusinessError("规则名称不能重复")
	}
	if !dao.AlarmRule.CheckRuleType(dto.TenantId, dto.Id, dto.Type) {
		return errors.NewBusinessError("规则类型不匹配")
	}
	for i := range dto.GroupList {
		if dto.GroupList[i] == "-1" {
			contactService := NewContactService(NewContactGroupService(NewContactGroupRelService()), NewContactInformationService(nil), NewContactGroupRelService())
			groupBizId, err := contactService.CreateSysContact(dto.TenantId)
			if err != nil {
				return err
			}
			dto.GroupList[i] = groupBizId
		}
	}
	err := dao.AlarmRule.UpdateRule(tx, dto)
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

func checkConditions(param *form.AlarmRuleAddReqDTO) error {
	if param.Type == constant.AlarmRuleTypeSingleMetric {
		if len(param.Conditions) == 0 {
			return errors.NewBusinessError("至少选择一个告警级别来设定触发条件")
		}
		if len(param.Conditions) > 4 {
			return errors.NewBusinessError("告警级别最多设置4条")
		}
		m := make(map[uint8]int)
		for i, v := range param.Conditions {
			if _, ok := m[v.Level]; ok {
				return errors.NewBusinessError("告警级别不能重复")
			} else {
				m[v.Level] = i
			}
		}
		return nil
	} else if param.Type == constant.AlarmRuleTypeMultipleMetric {
		if len(param.Conditions) == 0 || len(param.Conditions) > 10 {
			return errors.NewBusinessError("至少添加1条，最多可添加10条触发条件")
		}
	} else {
		return errors.NewBusinessError("规则类型错误")
	}
	return nil
}
