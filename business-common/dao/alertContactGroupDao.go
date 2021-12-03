package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/constants"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/errors"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/utils/snowflake"
	"gorm.io/gorm"
	"strconv"
)

type AlertContactGroupDao struct {
}

var AlertContactGroup = new(AlertContactGroupDao)

func (d *AlertContactGroupDao) SelectAlertContactGroup(db *gorm.DB, param forms.AlertContactParam) *[]forms.AlertContactGroupForm {
	var model = &[]forms.AlertContactGroupForm{}
	db.Raw(SelectAlterContactGroup, param.TenantId, param.GroupName).Find(model)
	return model
}

func (d *AlertContactGroupDao) SelectAlertGroupContact(db *gorm.DB, param forms.AlertContactParam) *[]forms.AlertContactForm {
	var model = &[]forms.AlertContactForm{}
	db.Raw(SelectAlterGroupContact, param.TenantId, param.GroupId).Find(model)
	return model
}

func (d *AlertContactGroupDao) Insert(db *gorm.DB, entity *models.AlertContactGroup) {
	db.Create(entity)
}

func (d *AlertContactGroupDao) Update(db *gorm.DB, entity *models.AlertContactGroup) {
	db.Updates(entity)
}

func (d *AlertContactGroupDao) Delete(db *gorm.DB, entity *models.AlertContactGroup) {
	db.Where("tenant_id = ? AND id = ?", entity.TenantId, entity.Id).Delete(models.AlertContactGroup{})
}

func (d *AlertContactGroupDao) InsertGroupRelBatch(db *gorm.DB, list []*models.AlertContactGroupRel) {
	now := tools.GetNowStr()
	for _, rel := range list {
		rel.Id = strconv.FormatInt(snowflake.GetWorker().NextId(), 10)
		rel.UpdateTime = now
		rel.CreateTime = now
	}
	db.Create(list)
}

func (d *AlertContactGroupDao) InsertAlertContactGroup(param forms.AlertContactGroupParam) error {
	var tx = global.DB.Begin()

	var count int64
	global.DB.Model(&models.AlertContactGroup{}).Where("tenant_id = ?", param.TenantId).Count(&count)
	if count >= constants.MaxGroupNum {
		return errors.NewBusinessError("联系组限制创建" + strconv.Itoa(constants.MaxGroupNum) + "个")
	}
	global.DB.Model(&models.AlertContactGroup{}).Where("tenant_id = ?", param.TenantId).Where("name = ?", param.GroupName).Count(&count)
	if count >= 1 {
		return errors.NewBusinessError("联系组名重复")
	}
	currentTime := tools.GetNowStr()
	groupId := strconv.FormatInt(snowflake.GetWorker().NextId(), 10)
	param.GroupId = groupId
	var alertContactGroup = models.AlertContactGroup{
		Id:          groupId,
		TenantId:    param.TenantId,
		Name:        param.GroupName,
		Description: param.Description,
		CreateUser:  param.CreateUser,
		CreateTime:  currentTime,
		UpdateTime:  currentTime,
	}
	tx = tx.Create(alertContactGroup)

	// 添加联系人组关联
	err := d.insertAlertContactGroupRel(param, currentTime)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (d *AlertContactGroupDao) UpdateAlertContactGroup(param forms.AlertContactGroupParam) error {
	var tx = global.DB.Begin()

	var count int64
	tx.Model(&models.AlertContactGroup{}).Where("tenant_id = ?", param.TenantId).Where("name = ?", param.GroupName).Count(&count)
	if count >= 1 {
		return errors.NewBusinessError("联系组名重复")
	}
	currentTime := tools.GetNowStr()
	var alertContactGroup = &models.AlertContactGroup{
		Id:          param.GroupId,
		TenantId:    param.TenantId,
		Name:        param.GroupName,
		Description: param.Description,
		UpdateTime:  currentTime,
	}
	tx.Model(alertContactGroup).Updates(alertContactGroup)

	//更新联系人关联
	err := d.updateAlertContactGroupRel(param, currentTime)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (d *AlertContactGroupDao) DeleteAlertContactGroup(param forms.AlertContactGroupParam) error {
	var tx = global.DB.Begin()

	var model models.AlertContactGroup
	tx.Delete(&model, param.GroupId)
	//删除联系人关联
	err := d.deleteAlertContactGroupRel(param.GroupId)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

//新增联系人关联
func (d *AlertContactGroupDao) insertAlertContactGroupRel(param forms.AlertContactGroupParam, currentTime string) error {
	//return errors.NewBusinessError("事务test")
	if len(param.ContactIdList) == 0 {
		return nil
	}
	var count int64
	global.DB.Model(&models.AlertContactGroupRel{}).Where("tenant_id = ?", param.TenantId).Where("group_id", param.GroupId).Count(&count)
	if count >= constants.MaxContactNum {
		return errors.NewBusinessError("每组联系人限制" + strconv.Itoa(constants.MaxContactNum) + "个")
	}
	for _, contactId := range param.ContactIdList {
		var alertContactGroupRel = &models.AlertContactGroupRel{
			Id:         strconv.FormatInt(snowflake.GetWorker().NextId(), 10),
			TenantId:   param.TenantId,
			ContactId:  contactId,
			GroupId:    param.GroupId,
			CreateUser: param.CreateUser,
			CreateTime: currentTime,
			UpdateTime: currentTime,
		}
		db := global.DB.Create(alertContactGroupRel)
		if db.Error != nil {
			return errors.NewBusinessError("添加失败")
		}
	}
	return nil
}

//更新联系人关联
func (d *AlertContactGroupDao) updateAlertContactGroupRel(param forms.AlertContactGroupParam, currentTime string) error {
	//清除旧联系人组关联
	deleteErr := d.deleteAlertContactGroupRel(param.GroupId)
	//添加新联系人组关联
	insertErr := d.insertAlertContactGroupRel(param, currentTime)
	if deleteErr != nil && insertErr != nil {
		return errors.NewBusinessError("修改失败")
	}
	return nil
}

//删除联系人关联
func (d *AlertContactGroupDao) deleteAlertContactGroupRel(groupId string) error {
	var tx = global.DB.Begin()
	db := tx.Where("group_id = ?", groupId).Delete(models.AlertContactGroupRel{})
	if db.Error != nil {
		return errors.NewBusinessError("删除失败")
	}
	return nil
}
