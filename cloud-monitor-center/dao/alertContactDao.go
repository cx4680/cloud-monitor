package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor-center/config"
	"code.cestc.cn/ccos-ops/cloud-monitor-center/database"
	"code.cestc.cn/ccos-ops/cloud-monitor-center/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor-center/models"
	"code.cestc.cn/ccos-ops/cloud-monitor-center/mq"
	"code.cestc.cn/ccos-ops/cloud-monitor-center/utils/snowflake"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/enums"
	"encoding/json"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type AlertContactDao struct {
	db *gorm.DB
}

func NewAlertContact(db *gorm.DB) *AlertContactDao {
	return &AlertContactDao{db: db}
}

var cfg = config.GetConfig()

func (mpd *AlertContactDao) GetAlertContact(param forms.AlertContactParam) *forms.AlertContactFormPage {
	var model = &[]forms.AlertContactForm{}
	db := mpd.db
	if param.Phone != "" {
		db = mpd.db.Raw(database.SelectAlterContact+"AND ac.id = ANY(SELECT contact_id FROM alert_contact_information WHERE type = 1 AND no LIKE CONCAT('%',?,'%')) ", param.TenantId, param.ContactName, param.GroupName, param.Phone).Group("ac.id")
	} else if param.Email != "" {
		db = mpd.db.Raw(database.SelectAlterContact+"AND ac.id = ANY(SELECT contact_id FROM alert_contact_information WHERE type = 2 AND no LIKE CONCAT('%',?,'%')) ", param.TenantId, param.ContactName, param.GroupName, param.Email).Group("ac.id")
	} else {
		db = mpd.db.Raw(database.SelectAlterContact, param.TenantId, param.ContactName, param.GroupName).Group("ac.id")
	}
	db.Find(model)
	total := len(*model)
	db.Limit(param.PageSize).Offset((param.PageCurrent - 1) * param.PageSize).Find(model)
	var alertContactFormPage = &forms.AlertContactFormPage{
		Records: model,
		Current: param.PageCurrent,
		Size:    param.PageSize,
		Total:   total,
	}
	return alertContactFormPage
}

func (mpd *AlertContactDao) InsertAlertContact(param forms.AlertContactParam) {
	currentTime := getCurrentTime()
	contactId := strconv.FormatInt(snowflake.GetWorker().NextId(), 10)
	param.ContactId = contactId
	var alertContact = &models.AlertContact{
		Id:          contactId,
		TenantId:    param.TenantId,
		Name:        param.ContactName,
		Status:      1,
		Description: param.Description,
		CreateUser:  param.CreateUser,
		CreateTime:  currentTime,
		UpdateTime:  currentTime,
	}
	mpd.db.Create(alertContact)
	//同步region
	mq.SendMsg(cfg.Rocketmq.AlertContactTopic, enums.InsertAlertContact, alertContact)

	//添加联系方式
	mpd.insertAlertContactInformation(param, param.Phone, 1, currentTime)
	mpd.insertAlertContactInformation(param, param.Email, 2, currentTime)

	//将联系人添加到组
	mpd.insertAlertContactGroupRel(param, contactId, currentTime)
}

func (mpd *AlertContactDao) UpdateAlertContact(param forms.AlertContactParam) {
	currentTime := getCurrentTime()
	var alertContact = &models.AlertContact{
		Id:          param.ContactId,
		TenantId:    param.TenantId,
		Name:        param.ContactName,
		Status:      1,
		Description: param.Description,
		UpdateTime:  currentTime,
	}
	mpd.db.Model(alertContact).Updates(alertContact)
	//同步region
	mq.SendMsg(cfg.Rocketmq.AlertContactTopic, enums.UpdateAlertContact, alertContact)

	//更新联系方式
	mpd.updateAlertContactInformation(param, currentTime)

	//更新联系人组关联
	mpd.updateAlertContactGroupRel(param, currentTime)
}

func (mpd *AlertContactDao) DeleteAlertContact(param forms.AlertContactParam) {
	var model models.AlertContact
	mpd.db.Delete(&model, param.ContactId)
	//同步region
	mq.SendMsg(cfg.Rocketmq.AlertContactTopic, enums.DeleteAlertContact, param.ContactId)

	//删除联系方式
	mpd.deleteAlertContactInformation(param.ContactId)
	//删除联系人组关联
	mpd.deleteAlertContactGroupRel(param.ContactId)
}

func (mpd *AlertContactDao) CertifyAlertContact(activeCode string) string {
	var model = &models.AlertContactInformation{}
	mpd.db.Model(model).Where("active_code = ?", activeCode).Update("is_certify", 1)
	//同步region
	mq.SendMsg(cfg.Rocketmq.AlertContactTopic, enums.CertifyAlertContact, activeCode)
	return getTenantName(model.TenantId)
}

func getCurrentTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

//创建联系方式
func (mpd *AlertContactDao) insertAlertContactInformation(param forms.AlertContactParam, no string, noType int, currentTime string) {
	//空号码则不保存
	if no == "" {
		return
	}
	activeCode := uuid.NewV4().String()
	var isCertify int
	if config.GetConfig().HasNoticeModel {
		isCertify = 0
	} else {
		isCertify = 1
	}
	var alertContactInformation = &models.AlertContactInformation{
		Id:         strconv.FormatInt(snowflake.GetWorker().NextId(), 10),
		TenantId:   param.TenantId,
		ContactId:  param.ContactId,
		No:         no,
		Type:       noType,
		IsCertify:  isCertify,
		ActiveCode: activeCode,
		CreateUser: param.CreateUser,
		CreateTime: currentTime,
		UpdateTime: currentTime,
	}
	mpd.db.Create(alertContactInformation)
	// 同步region
	mq.SendMsg(cfg.Rocketmq.AlertContactTopic, enums.InsertAlertContactInformation, alertContactInformation)
	// TODO 发送验证消息
}

//创建联系人组关联
func (mpd *AlertContactDao) insertAlertContactGroupRel(param forms.AlertContactParam, contactId string, currentTime string) {
	for i := range param.GroupIdList {
		var alertContactGroupRel = &models.AlertContactGroupRel{
			Id:         strconv.FormatInt(snowflake.GetWorker().NextId(), 10),
			TenantId:   param.TenantId,
			ContactId:  contactId,
			GroupId:    param.GroupIdList[i],
			CreateUser: param.CreateUser,
			CreateTime: currentTime,
			UpdateTime: currentTime,
		}
		mpd.db.Create(alertContactGroupRel)
		// 同步region
		mq.SendMsg(cfg.Rocketmq.AlertContactTopic, enums.InsertAlertContactGroupRel, alertContactGroupRel)
	}
}

//更新联系方式
func (mpd *AlertContactDao) updateAlertContactInformation(param forms.AlertContactParam, currentTime string) {
	//删除旧的联系方式
	mpd.deleteAlertContactInformation(param.ContactId)
	//再新增新的联系方式
	mpd.insertAlertContactInformation(param, param.Phone, 1, currentTime)
	mpd.insertAlertContactInformation(param, param.Email, 2, currentTime)
}

//更新联系人组关联
func (mpd *AlertContactDao) updateAlertContactGroupRel(param forms.AlertContactParam, currentTime string) {
	//清除联系人组的关联
	mpd.deleteAlertContactGroupRel(param.ContactId)
	//再新增新的关联
	mpd.insertAlertContactGroupRel(param, param.ContactId, currentTime)
}

//删除联系方式
func (mpd *AlertContactDao) deleteAlertContactInformation(contactId string) {
	mpd.db.Where("contact_id = ?", contactId).Delete(models.AlertContactInformation{})
	//同步region
	mq.SendMsg(cfg.Rocketmq.AlertContactTopic, enums.DeleteAlertContactInformation, contactId)
}

//删除联系人组关联
func (mpd *AlertContactDao) deleteAlertContactGroupRel(contactId string) {
	mpd.db.Where("contact_id = ?", contactId).Delete(models.AlertContactGroupRel{})
	//同步region
	mq.SendMsg(cfg.Rocketmq.AlertContactTopic, enums.DeleteAlertContactGroupRelByContactId, contactId)
}

//获取租户名字
func getTenantName(tenantId string) string {
	var request = strings.NewReader("{\"loginId\":\"" + tenantId + "\"}")
	response, err := http.Post(config.GetConfig().TenantUrl, "application/json; charset=utf-8", request)
	if err != nil {
		return "未命名"
	}
	data, err := ioutil.ReadAll(response.Body)
	var result map[string]map[string]map[string]string
	err = json.Unmarshal(data, &result)
	return result["module"]["login"]["loginCode"]
}
