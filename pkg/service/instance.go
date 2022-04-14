package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/enum"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global/sys_component/sys_rocketmq"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/mq"
	"gorm.io/gorm"
)

func UnbindInstance(tx *gorm.DB, param interface{}) error {
	instanceDao := dao.Instance
	dto := param.(*form.UnBindRuleParam)
	err := instanceDao.UnbindInstance(tx, dto)
	if err != nil {
		return err
	}
	return mq.SendMsg(sys_rocketmq.RuleTopic, enum.UnbindRule, param)
}

func BindInstance(tx *gorm.DB, param interface{}) error {
	instanceDao := dao.Instance
	dto := param.(*form.InstanceBindRuleDTO)
	err := instanceDao.BindInstance(tx, dto)
	if err != nil {
		return err
	}
	return mq.SendMsg(sys_rocketmq.RuleTopic, enum.BindRule, param)
}
