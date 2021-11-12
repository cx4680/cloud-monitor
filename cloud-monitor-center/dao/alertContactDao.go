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
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
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

var rocketmqConfig = config.GetRocketmqConfig()

func (mpd *AlertContactDao) GetAlertContact(param forms.AlertContactParam) *forms.AlertContactFormPage {
	var model = &[]forms.AlertContactForm{}
	var sql string
	if param.ContactName != "" {
		sql = fmt.Sprintf(database.SelectAlterContact, param.TenantId, "AND ac.name LIKE CONCAT('%',"+param.ContactName+",'%')")
	}
	if param.Phone != "" {
		sql = fmt.Sprintf(database.SelectAlterContact, param.TenantId, "AND ac.id = ANY(SELECT contact_id FROM alert_contact_information WHERE type = 1 AND no LIKE CONCAT('%',"+param.Phone+",'%')) ")
	} else if param.Email != "" {
		sql = fmt.Sprintf(database.SelectAlterContact, param.TenantId, "AND ac.id = ANY(SELECT contact_id FROM alert_contact_information WHERE type = 1 AND no LIKE CONCAT('%',"+param.Email+",'%')) ")
	} else {
		sql = fmt.Sprintf(database.SelectAlterContact, param.TenantId, "")
	}
	var count int64
	mpd.db.Model(&models.AlertContact{}).Where("tenant_id = ?", param.TenantId).Count(&count)
	total := count
	sql += "LIMIT " + strconv.Itoa((param.PageCurrent-1)*param.PageSize) + "," + strconv.Itoa(param.PageSize)
	mpd.db.Raw(sql).Find(model)
	var alertContactFormPage = &forms.AlertContactFormPage{
		Records: model,
		Current: param.PageCurrent,
		Size:    param.PageSize,
		Total:   total,
	}
	return alertContactFormPage
}

func (mpd *AlertContactDao) InsertAlertContact(param forms.AlertContactParam) error {
	if param.ContactName == "" {
		return errors.NewError("联系人名字不能为空")
	}
	//每个账号限制创建100个联系人
	var count int64
	mpd.db.Model(&models.AlertContact{}).Where("tenant_id = ?", param.TenantId).Count(&count)
	if count >= constant.MAX_CONTACT_NUM {
		return errors.NewError("联系人限制创建" + strconv.Itoa(constant.MAX_CONTACT_NUM) + "个")
	}
	//每个联系人最多加入5个联系组
	if len(param.GroupIdList) >= constant.MAX_CONTACT_GROUP {
		return errors.NewError("每个联系人最多加入" + strconv.Itoa(constant.MAX_CONTACT_GROUP) + "个联系组")
	}

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

	//添加联系方式
	mpd.insertAlertContactInformation(param, param.Phone, 1, currentTime)
	mpd.insertAlertContactInformation(param, param.Email, 2, currentTime)

	//将联系人添加到组
	err := mpd.insertAlertContactGroupRel(param, contactId, currentTime)
	if err != nil {
		return err
	}
	//同步region
	mq.SendMsg(rocketmqConfig.AlertContactTopic, enums.InsertAlertContact, alertContact)
	return nil
}

func (mpd *AlertContactDao) UpdateAlertContact(param forms.AlertContactParam) error {
	if param.ContactName == "" {
		return errors.NewError("联系人名字不能为空")
	}
	//每个联系人最多加入5个联系组
	if len(param.GroupIdList) >= constant.MAX_CONTACT_GROUP {
		return errors.NewError("每个联系人最多加入" + strconv.Itoa(constant.MAX_CONTACT_GROUP) + "个联系组")
	}
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

	//更新联系方式
	mpd.updateAlertContactInformation(param, currentTime)

	//更新联系人组关联
	err := mpd.updateAlertContactGroupRel(param, currentTime)
	if err != nil {
		return err
	}
	//同步region
	mq.SendMsg(rocketmqConfig.AlertContactTopic, enums.UpdateAlertContact, alertContact)
	return nil
}

func (mpd *AlertContactDao) DeleteAlertContact(param forms.AlertContactParam) {
	var model models.AlertContact
	mpd.db.Delete(&model, param.ContactId)
	//删除联系方式
	mpd.deleteAlertContactInformation(param.ContactId)
	//删除联系人组关联
	mpd.deleteAlertContactGroupRel(param.ContactId)
	//同步region
	mq.SendMsg(rocketmqConfig.AlertContactTopic, enums.DeleteAlertContact, param.ContactId)
}

func (mpd *AlertContactDao) CertifyAlertContact(activeCode string) string {
	var model = &models.AlertContactInformation{}
	mpd.db.Model(model).Where("active_code = ?", activeCode).Update("is_certify", 1)
	//同步region
	mq.SendMsg(rocketmqConfig.AlertContactTopic, enums.CertifyAlertContact, activeCode)
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
	activeCode := strings.ReplaceAll(uuid.New().String(), "-", "")
	var isCertify int
	if config.GetCommonConfig().HasNoticeModel {
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
	// TODO 发送验证消息
	sendMsg(param.TenantId, param.ContactId, no, noType, activeCode)
	// 同步region
	mq.SendMsg(rocketmqConfig.AlertContactTopic, enums.InsertAlertContactInformation, alertContactInformation)
}

//创建联系人组关联
func (mpd *AlertContactDao) insertAlertContactGroupRel(param forms.AlertContactParam, contactId string, currentTime string) error {
	var count int64
	for _, v := range param.GroupIdList {
		mpd.db.Model(&models.AlertContactGroupRel{}).Where("tenant_id = ?", param.TenantId).Where("group_id = ?", v).Count(&count)
		if count >= constant.MAX_CONTACT_NUM {
			return errors.NewError("每组联系人限制创建" + strconv.Itoa(constant.MAX_CONTACT_NUM) + "个")
		}
		var alertContactGroupRel = &models.AlertContactGroupRel{
			Id:         strconv.FormatInt(snowflake.GetWorker().NextId(), 10),
			TenantId:   param.TenantId,
			ContactId:  contactId,
			GroupId:    v,
			CreateUser: param.CreateUser,
			CreateTime: currentTime,
			UpdateTime: currentTime,
		}
		mpd.db.Create(alertContactGroupRel)
		// 同步region
		mq.SendMsg(rocketmqConfig.AlertContactTopic, enums.InsertAlertContactGroupRel, alertContactGroupRel)
	}
	return nil
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
func (mpd *AlertContactDao) updateAlertContactGroupRel(param forms.AlertContactParam, currentTime string) error {
	//清除联系人组的关联
	mpd.deleteAlertContactGroupRel(param.ContactId)
	//再新增新的关联
	err := mpd.insertAlertContactGroupRel(param, param.ContactId, currentTime)
	if err != nil {
		return err
	}
	return nil
}

//删除联系方式
func (mpd *AlertContactDao) deleteAlertContactInformation(contactId string) {
	mpd.db.Where("contact_id = ?", contactId).Delete(models.AlertContactInformation{})
	//同步region
	mq.SendMsg(rocketmqConfig.AlertContactTopic, enums.DeleteAlertContactInformation, contactId)
}

//删除联系人组关联
func (mpd *AlertContactDao) deleteAlertContactGroupRel(contactId string) {
	mpd.db.Where("contact_id = ?", contactId).Delete(models.AlertContactGroupRel{})
	//同步region
	mq.SendMsg(rocketmqConfig.AlertContactTopic, enums.DeleteAlertContactGroupRelByContactId, contactId)
}

//获取租户名字
func getTenantName(tenantId string) string {
	var request = strings.NewReader("{\"loginId\":\"" + tenantId + "\"}")
	response, err := http.Post(config.GetCommonConfig().TenantUrl, "application/json; charset=utf-8", request)
	if err != nil {
		return "未命名"
	}
	data, err := ioutil.ReadAll(response.Body)
	var result map[string]map[string]map[string]string
	err = json.Unmarshal(data, &result)
	if result["module"] != nil {
		return result["module"]["login"]["loginCode"]
	} else {
		return result["result"]["login"]["loginCode"]
	}
}

func sendMsg(tenantId string, contactId string, no string, noType int, activeCode string) {
	if no == "" {
		return
	}
	params := make(map[string]string)
	params["userName"] = getTenantName(tenantId)
	params["verifyBtn"] = "/#/alarm/activation?code=" + activeCode
	params["activationlink"] = "code=" + activeCode

}
