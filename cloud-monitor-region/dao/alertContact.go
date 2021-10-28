package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/models"
	"gorm.io/gorm"
)

type AlertContactDao struct {
	db *gorm.DB
}

func NewAlertContact(db *gorm.DB) *AlertContactDao {
	return &AlertContactDao{db: db}
}

func (mpd *AlertContactDao) InsertAlertContact(model models.AlertContact) {
	mpd.db.Create(model)
}

func (mpd *AlertContactDao) UpdateAlertContact(model models.AlertContact) {
	mpd.db.Model(model).Updates(model)
}

func (mpd *AlertContactDao) DeleteAlertContact(contactId string) {
	var model models.AlertContact
	mpd.db.Delete(&model, contactId)
}

func (mpd *AlertContactDao) InsertAlertContactInformation(model models.AlertContactInformation) {
	mpd.db.Create(model)
}

func (mpd *AlertContactDao) DeleteAlertContactInformation(contactId string) {
	mpd.db.Where("contact_id = ?", contactId).Delete(models.AlertContactGroupRel{})
}

func (mpd *AlertContactDao) InsertAlertContactGroupRel(model models.AlertContactGroupRel) {
	mpd.db.Create(model)
}

func (mpd *AlertContactDao) DeleteAlertContactGroupRelByContactId(contactId string) {
	mpd.db.Where("contact_id = ?", contactId).Delete(models.AlertContactGroupRel{})
}

func (mpd *AlertContactDao) CertifyAlertContact(activeCode string) {
	var model = &models.AlertContactInformation{}
	mpd.db.Model(model).Where("active_code = ?", activeCode).Update("is_certify", 1)
}

func (mpd *AlertContactDao) InsertAlertContactGroup(model models.AlertContactGroup) {
	mpd.db.Create(model)
}

func (mpd *AlertContactDao) UpdateAlertContactGroup(model models.AlertContactGroup) {
	mpd.db.Model(model).Updates(model)
}

func (mpd *AlertContactDao) DeleteAlertContactGroup(groupId string) {
	var model models.AlertContactGroup
	mpd.db.Delete(&model, groupId)
}

func (mpd *AlertContactDao) DeleteAlertContactGroupRelByGroupId(groupId string) {
	mpd.db.Where("group_id = ?", groupId).Delete(models.AlertContactGroupRel{})
}
