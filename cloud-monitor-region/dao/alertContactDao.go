package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/models"
)

type AlertContactDao struct {
}

var AlertContact = new(AlertContactDao)

func (mpd *AlertContactDao) InsertAlertContact(model models.AlertContact) {
	global.DB.Create(model)
}

func (mpd *AlertContactDao) UpdateAlertContact(model models.AlertContact) {
	global.DB.Model(model).Updates(model)
}

func (mpd *AlertContactDao) DeleteAlertContact(contactId string) {
	var model models.AlertContact
	global.DB.Delete(&model, contactId)
}

func (mpd *AlertContactDao) InsertAlertContactInformation(model models.AlertContactInformation) {
	global.DB.Create(model)
}

func (mpd *AlertContactDao) DeleteAlertContactInformation(contactId string) {
	global.DB.Where("contact_id = ?", contactId).Delete(models.AlertContactGroupRel{})
}

func (mpd *AlertContactDao) InsertAlertContactGroupRel(model models.AlertContactGroupRel) {
	global.DB.Create(model)
}

func (mpd *AlertContactDao) DeleteAlertContactGroupRelByContactId(contactId string) {
	global.DB.Where("contact_id = ?", contactId).Delete(models.AlertContactGroupRel{})
}

func (mpd *AlertContactDao) CertifyAlertContact(activeCode string) {
	var model = &models.AlertContactInformation{}
	global.DB.Model(model).Where("active_code = ?", activeCode).Update("is_certify", 1)
}

func (mpd *AlertContactDao) InsertAlertContactGroup(model models.AlertContactGroup) {
	global.DB.Create(model)
}

func (mpd *AlertContactDao) UpdateAlertContactGroup(model models.AlertContactGroup) {
	global.DB.Model(model).Updates(model)
}

func (mpd *AlertContactDao) DeleteAlertContactGroup(groupId string) {
	var model models.AlertContactGroup
	global.DB.Delete(&model, groupId)
}

func (mpd *AlertContactDao) DeleteAlertContactGroupRelByGroupId(groupId string) {
	global.DB.Where("group_id = ?", groupId).Delete(models.AlertContactGroupRel{})
}
