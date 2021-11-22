package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/database"
)

type AlertContactDao struct {
}

var AlertContact = new(AlertContactDao)

func (mpd *AlertContactDao) InsertAlertContact(model models.AlertContact) {
	database.GetDb().Create(model)
}

func (mpd *AlertContactDao) UpdateAlertContact(model models.AlertContact) {
	database.GetDb().Model(model).Updates(model)
}

func (mpd *AlertContactDao) DeleteAlertContact(contactId string) {
	var model models.AlertContact
	database.GetDb().Delete(&model, contactId)
}

func (mpd *AlertContactDao) InsertAlertContactInformation(model models.AlertContactInformation) {
	database.GetDb().Create(model)
}

func (mpd *AlertContactDao) DeleteAlertContactInformation(contactId string) {
	database.GetDb().Where("contact_id = ?", contactId).Delete(models.AlertContactGroupRel{})
}

func (mpd *AlertContactDao) InsertAlertContactGroupRel(model models.AlertContactGroupRel) {
	database.GetDb().Create(model)
}

func (mpd *AlertContactDao) DeleteAlertContactGroupRelByContactId(contactId string) {
	database.GetDb().Where("contact_id = ?", contactId).Delete(models.AlertContactGroupRel{})
}

func (mpd *AlertContactDao) CertifyAlertContact(activeCode string) {
	var model = &models.AlertContactInformation{}
	database.GetDb().Model(model).Where("active_code = ?", activeCode).Update("is_certify", 1)
}

func (mpd *AlertContactDao) InsertAlertContactGroup(model models.AlertContactGroup) {
	database.GetDb().Create(model)
}

func (mpd *AlertContactDao) UpdateAlertContactGroup(model models.AlertContactGroup) {
	database.GetDb().Model(model).Updates(model)
}

func (mpd *AlertContactDao) DeleteAlertContactGroup(groupId string) {
	var model models.AlertContactGroup
	database.GetDb().Delete(&model, groupId)
}

func (mpd *AlertContactDao) DeleteAlertContactGroupRelByGroupId(groupId string) {
	database.GetDb().Where("group_id = ?", groupId).Delete(models.AlertContactGroupRel{})
}
