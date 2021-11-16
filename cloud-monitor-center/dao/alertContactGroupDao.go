package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor-center/constant"
	"code.cestc.cn/ccos-ops/cloud-monitor-center/database"
	"code.cestc.cn/ccos-ops/cloud-monitor-center/errors"
	"code.cestc.cn/ccos-ops/cloud-monitor-center/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor-center/models"
	"code.cestc.cn/ccos-ops/cloud-monitor-center/mq"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/enums"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/utils/snowflake"
	"gorm.io/gorm"
	"strconv"
)

type AlertContactGroupDao struct {
	db *gorm.DB
}

func NewAlertContactGroup(db *gorm.DB) *AlertContactGroupDao {
	return &AlertContactGroupDao{db: db}
}

func (mpd *AlertContactGroupDao) GetAlertContactGroup(tenantId string, groupName string) *[]forms.AlertContactGroupForm {
	var model = &[]forms.AlertContactGroupForm{}
	mpd.db.Raw(database.SelectAlterContactGroup, tenantId, groupName).Find(model)
	return model
}

func (mpd *AlertContactGroupDao) GetAlertGroupContact(tenantId string, groupId string) *[]forms.AlertContactForm {
	var model = &[]forms.AlertContactForm{}
	mpd.db.Raw(database.SelectAlterGroupContact, tenantId, groupId).Find(model)
	return model
}

func (mpd *AlertContactGroupDao) InsertAlertContactGroup(param forms.AlertContactGroupParam) error {
	var tx = mpd.db.Begin()

	var count int64
	mpd.db.Model(&models.AlertContactGroup{}).Where("tenant_id = ?", param.TenantId).Count(&count)
	if count >= constant.MAX_GROUP_NUM {
		return errors.NewBusinessError("联系组限制创建" + strconv.Itoa(constant.MAX_GROUP_NUM) + "个")
	}
	mpd.db.Model(&models.AlertContactGroup{}).Where("tenant_id = ?", param.TenantId).Where("name = ?", param.GroupName).Count(&count)
	if count >= 1 {
		return errors.NewBusinessError("联系组名重复")
	}
	currentTime := getCurrentTime()
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
	err := mpd.insertAlertContactGroupRel(param, currentTime)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	//同步region
	mq.SendMsg(config.GetRocketmqConfig().AlertContactTopic, enums.InsertAlertContactGroup, alertContactGroup)
	return nil
}

func (mpd *AlertContactGroupDao) UpdateAlertContactGroup(param forms.AlertContactGroupParam) error {
	var tx = mpd.db.Begin()

	var count int64
	tx.Model(&models.AlertContactGroup{}).Where("tenant_id = ?", param.TenantId).Where("name = ?", param.GroupName).Count(&count)
	if count >= 1 {
		return errors.NewBusinessError("联系组名重复")
	}
	currentTime := getCurrentTime()
	var alertContactGroup = &models.AlertContactGroup{
		Id:          param.GroupId,
		TenantId:    param.TenantId,
		Name:        param.GroupName,
		Description: param.Description,
		UpdateTime:  currentTime,
	}
	tx.Model(alertContactGroup).Updates(alertContactGroup)

	//更新联系人关联
	err := mpd.updateAlertContactGroupRel(param, currentTime)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	//同步region
	mq.SendMsg(config.GetRocketmqConfig().AlertContactTopic, enums.UpdateAlertContactGroup, alertContactGroup)
	return nil
}

func (mpd *AlertContactGroupDao) DeleteAlertContactGroup(param forms.AlertContactGroupParam) error {
	var tx = mpd.db.Begin()

	var model models.AlertContactGroup
	tx.Delete(&model, param.GroupId)
	//删除联系人关联
	err := mpd.deleteAlertContactGroupRel(param.GroupId)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	//同步region
	mq.SendMsg(config.GetRocketmqConfig().AlertContactTopic, enums.DeleteAlertContactGroup, param.GroupId)
	return nil
}

//新增联系人关联
func (mpd *AlertContactGroupDao) insertAlertContactGroupRel(param forms.AlertContactGroupParam, currentTime string) error {
	//return errors.NewBusinessError("事务test")
	if len(param.ContactIdList) == 0 {
		return nil
	}
	var count int64
	mpd.db.Model(&models.AlertContactGroupRel{}).Where("tenant_id = ?", param.TenantId).Where("group_id", param.GroupId).Count(&count)
	if count >= constant.MAX_CONTACT_NUM {
		return errors.NewBusinessError("每组联系人限制" + strconv.Itoa(constant.MAX_CONTACT_NUM) + "个")
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
		db := mpd.db.Create(alertContactGroupRel)
		if db.Error != nil {
			return errors.NewBusinessError("添加失败")
		}
		//同步region
		mq.SendMsg(config.GetRocketmqConfig().AlertContactTopic, enums.InsertAlertContactGroupRel, alertContactGroupRel)
	}
	return nil
}

//更新联系人关联
func (mpd *AlertContactGroupDao) updateAlertContactGroupRel(param forms.AlertContactGroupParam, currentTime string) error {
	//清除旧联系人组关联
	deleteErr := mpd.deleteAlertContactGroupRel(param.GroupId)
	//添加新联系人组关联
	insertErr := mpd.insertAlertContactGroupRel(param, currentTime)
	if deleteErr != nil && insertErr != nil {
		return errors.NewBusinessError("修改失败")
	}
	return nil
}

//删除联系人关联
func (mpd *AlertContactGroupDao) deleteAlertContactGroupRel(groupId string) error {
	var tx = mpd.db.Begin()
	db := tx.Where("group_id = ?", groupId).Delete(models.AlertContactGroupRel{})
	if db.Error != nil {
		return errors.NewBusinessError("删除失败")
	}
	//同步region
	mq.SendMsg(config.GetRocketmqConfig().AlertContactTopic, enums.DeleteAlertContactGroupRelByGroupId, groupId)
	return nil
}
