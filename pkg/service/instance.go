package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/form"
	"gorm.io/gorm"
)

func UnbindInstance(tx *gorm.DB, param interface{}) error {
	instanceDao := dao.Instance
	dto := param.(*form.UnBindRuleParam)
	err := instanceDao.UnbindInstance(tx, dto)
	if err != nil {
		return err
	}
	return nil
}

func BindInstance(tx *gorm.DB, param interface{}) error {
	instanceDao := dao.Instance
	dto := param.(*form.InstanceBindRuleDTO)
	err := instanceDao.BindInstance(tx, dto)
	if err != nil {
		return err
	}
	return nil
}
