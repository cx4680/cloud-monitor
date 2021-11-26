package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dtos"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/constant"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/errors"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/mq"
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
)

type AlertContactDao struct {
}

var AlertContact = new(AlertContactDao)

const (
	SelectAlterContact = "SELECT " +
		"ac.id AS contact_id, " +
		"ac.name AS contact_name, " +
		"acg.group_id AS group_id, " +
		"acg.group_name AS group_name, " +
		"GROUP_CONCAT( CASE aci.type WHEN 1 THEN aci.NO END ) AS phone, " +
		"GROUP_CONCAT( CASE aci.type WHEN 1 THEN aci.is_certify END ) AS phone_certify, " +
		"GROUP_CONCAT( CASE aci.type WHEN 2 THEN aci.NO END ) AS email, " +
		"GROUP_CONCAT( CASE aci.type WHEN 2 THEN aci.is_certify END ) AS email_certify, " +
		"GROUP_CONCAT( CASE aci.type WHEN 3 THEN aci.NO END ) AS lanxin, " +
		"GROUP_CONCAT( CASE aci.type WHEN 3 THEN aci.is_certify END ) AS lanxin_certify, " +
		"ac.description AS description " +
		"FROM " +
		"alert_contact AS ac " +
		"LEFT JOIN alert_contact_information AS aci ON ac.id = aci.contact_id " +
		"LEFT JOIN ( " +
		"SELECT " +
		"acgr.contact_id AS contact_id, " +
		"GROUP_CONCAT( acg.id ) AS group_id, " +
		"GROUP_CONCAT( acg.name ) AS group_name " +
		"FROM " +
		"alert_contact_group AS acg " +
		"LEFT JOIN alert_contact_group_rel AS acgr ON acg.id = acgr.group_id  " +
		"GROUP BY " +
		"acgr.contact_id ) " +
		"AS acg ON acg.contact_id = ac.id " +
		"WHERE " +
		"ac.status = 1 " +
		"AND ac.tenant_id = %s " +
		"%s" +
		"GROUP BY " +
		"ac.id " +
		"ORDER BY " +
		"ac.create_time DESC  "

	SelectAlterContactGroup = "SELECT " +
		"acg.id AS group_id, " +
		"acg.name AS group_name, " +
		"acg.description AS description, " +
		"acg.create_time AS create_time, " +
		"acg.update_time AS update_time, " +
		"COUNT( acgr.group_id ) AS contact_count " +
		"FROM " +
		"alert_contact_group AS acg " +
		"LEFT JOIN alert_contact_group_rel AS acgr ON acg.id = acgr.group_id " +
		"WHERE " +
		"acg.tenant_id = ? " +
		"AND acg.name LIKE CONCAT('%',?,'%') " +
		"GROUP BY " +
		"acg.id " +
		"ORDER BY " +
		"acg.create_time DESC "

	SelectAlterGroupContact = "SELECT " +
		"ac.id AS contact_id, " +
		"ac.name AS contact_name, " +
		"acg.group_id AS group_id, " +
		"acg.group_name AS group_name, " +
		"GROUP_CONCAT( CASE aci.type WHEN 1 THEN aci.NO END ) AS phone, " +
		"GROUP_CONCAT( CASE aci.type WHEN 1 THEN aci.is_certify END ) AS phone_certify, " +
		"GROUP_CONCAT( CASE aci.type WHEN 2 THEN aci.NO END ) AS email, " +
		"GROUP_CONCAT( CASE aci.type WHEN 2 THEN aci.is_certify END ) AS email_certify, " +
		"GROUP_CONCAT( CASE aci.type WHEN 3 THEN aci.NO END ) AS lanxin, " +
		"GROUP_CONCAT( CASE aci.type WHEN 3 THEN aci.is_certify END ) AS lanxin_certify, " +
		"ac.description AS description " +
		"FROM " +
		"alert_contact AS ac " +
		"LEFT JOIN alert_contact_information AS aci ON ac.id = aci.contact_id " +
		"LEFT JOIN ( " +
		"SELECT " +
		"acgr.contact_id AS contact_id, " +
		"GROUP_CONCAT( acg.id ) AS group_id, " +
		"GROUP_CONCAT( acg.name ) AS group_name " +
		"FROM " +
		"alert_contact_group AS acg " +
		"LEFT JOIN alert_contact_group_rel AS acgr ON acg.id = acgr.group_id  " +
		"GROUP BY " +
		"acgr.contact_id ) " +
		"AS acg ON acg.contact_id = ac.id " +
		"WHERE " +
		"ac.status = 1 " +
		"AND ac.tenant_id = ? " +
		"AND acg.group_id = ? " +
		"GROUP BY " +
		"ac.id " +
		"ORDER BY " +
		"ac.create_time DESC "
)

func (mpd *AlertContactDao) GetAlertContact(param forms.AlertContactParam) *forms.AlertContactFormPage {
	var model = &[]forms.AlertContactForm{}
	var sql string
	if param.ContactName != "" {
		sql = fmt.Sprintf(SelectAlterContact, param.TenantId, "AND ac.name LIKE CONCAT('%',"+param.ContactName+",'%')")
	}
	if param.Phone != "" {
		sql = fmt.Sprintf(SelectAlterContact, param.TenantId, "AND ac.id = ANY(SELECT contact_id FROM alert_contact_information WHERE type = 1 AND no LIKE CONCAT('%',"+param.Phone+",'%')) ")
	} else if param.Email != "" {
		sql = fmt.Sprintf(SelectAlterContact, param.TenantId, "AND ac.id = ANY(SELECT contact_id FROM alert_contact_information WHERE type = 1 AND no LIKE CONCAT('%',"+param.Email+",'%')) ")
	} else {
		sql = fmt.Sprintf(SelectAlterContact, param.TenantId, "")
	}
	var count int64
	global.DB.Model(&models.AlertContact{}).Where("tenant_id = ?", param.TenantId).Count(&count)
	total := count
	sql += "LIMIT " + strconv.Itoa((param.PageCurrent-1)*param.PageSize) + "," + strconv.Itoa(param.PageSize)
	global.DB.Raw(sql).Find(model)
	var alertContactFormPage = &forms.AlertContactFormPage{
		Records: model,
		Current: param.PageCurrent,
		Size:    param.PageSize,
		Total:   total,
	}
	return alertContactFormPage
}

func (acd *AlertContactDao) Insert(db *gorm.DB, entity *models.AlertContact) {
	currentTime := tools.GetNowStr()
	contactId := strconv.FormatInt(snowflake.GetWorker().NextId(), 10)

	entity.Id = contactId
	entity.CreateTime = currentTime
	entity.UpdateTime = currentTime
	entity.Status = 1
	db.Create(entity)
}
func (mpd *AlertContactDao) InsertAlertContact(param forms.AlertContactParam) error {
	if param.ContactName == "" {
		return errors.NewBusinessError("联系人名字不能为空")
	}
	//每个账号限制创建100个联系人
	var count int64
	global.DB.Model(&models.AlertContact{}).Where("tenant_id = ?", param.TenantId).Count(&count)
	if count >= constant.MAX_CONTACT_NUM {
		return errors.NewBusinessError("联系人限制创建" + strconv.Itoa(constant.MAX_CONTACT_NUM) + "个")
	}
	//每个联系人最多加入5个联系组
	if len(param.GroupIdList) >= constant.MAX_CONTACT_GROUP {
		return errors.NewBusinessError("每个联系人最多加入" + strconv.Itoa(constant.MAX_CONTACT_GROUP) + "个联系组")
	}

	currentTime := tools.GetNowStr()
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
	global.DB.Create(alertContact)

	//添加联系方式
	mpd.insertAlertContactInformation(param, param.Phone, 1, currentTime)
	mpd.insertAlertContactInformation(param, param.Email, 2, currentTime)

	//将联系人添加到组
	err := mpd.insertAlertContactGroupRel(param, contactId, currentTime)
	if err != nil {
		return err
	}
	//同步region
	mq.SendMsg(config.GetRocketmqConfig().AlertContactTopic, enums.InsertAlertContact, alertContact)
	return nil
}

func (mpd *AlertContactDao) UpdateAlertContact(param forms.AlertContactParam) error {
	if param.ContactName == "" {
		return errors.NewBusinessError("联系人名字不能为空")
	}
	//每个联系人最多加入5个联系组
	if len(param.GroupIdList) >= constant.MAX_CONTACT_GROUP {
		return errors.NewBusinessError("每个联系人最多加入" + strconv.Itoa(constant.MAX_CONTACT_GROUP) + "个联系组")
	}
	currentTime := tools.GetNowStr()
	var alertContact = &models.AlertContact{
		Id:          param.ContactId,
		TenantId:    param.TenantId,
		Name:        param.ContactName,
		Status:      1,
		Description: param.Description,
		UpdateTime:  currentTime,
	}
	global.DB.Model(alertContact).Updates(alertContact)

	//更新联系方式
	mpd.updateAlertContactInformation(param, currentTime)

	//更新联系人组关联
	err := mpd.updateAlertContactGroupRel(param, currentTime)
	if err != nil {
		return err
	}
	//同步region
	mq.SendMsg(config.GetRocketmqConfig().AlertContactTopic, enums.UpdateAlertContact, alertContact)
	return nil
}

func (mpd *AlertContactDao) DeleteAlertContact(param forms.AlertContactParam) {
	var model models.AlertContact
	global.DB.Delete(&model, param.ContactId)
	//删除联系方式
	mpd.deleteAlertContactInformation(param.ContactId)
	//删除联系人组关联
	mpd.deleteAlertContactGroupRel(param.ContactId)
	//同步region
	mq.SendMsg(config.GetRocketmqConfig().AlertContactTopic, enums.DeleteAlertContact, param.ContactId)
}

func (mpd *AlertContactDao) CertifyAlertContact(activeCode string) string {
	var model = &models.AlertContactInformation{}
	global.DB.Model(model).Where("active_code = ?", activeCode).Update("is_certify", 1)
	//同步region
	mq.SendMsg(config.GetRocketmqConfig().AlertContactTopic, enums.CertifyAlertContact, activeCode)
	return getTenantName(model.TenantId)
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
	global.DB.Create(alertContactInformation)
	// TODO 发送验证消息
	//sendMsg(param.TenantId, param.ContactId, no, noType, activeCode)
	// 同步region
	mq.SendMsg(config.GetRocketmqConfig().AlertContactTopic, enums.InsertAlertContactInformation, alertContactInformation)
}

//创建联系人组关联
func (mpd *AlertContactDao) insertAlertContactGroupRel(param forms.AlertContactParam, contactId string, currentTime string) error {
	var count int64
	for _, v := range param.GroupIdList {
		global.DB.Model(&models.AlertContactGroupRel{}).Where("tenant_id = ?", param.TenantId).Where("group_id = ?", v).Count(&count)
		if count >= constant.MAX_CONTACT_NUM {
			return errors.NewBusinessError("每组联系人限制创建" + strconv.Itoa(constant.MAX_CONTACT_NUM) + "个")
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
		global.DB.Create(alertContactGroupRel)
		// 同步region
		mq.SendMsg(config.GetRocketmqConfig().AlertContactTopic, enums.InsertAlertContactGroupRel, alertContactGroupRel)
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
	global.DB.Where("contact_id = ?", contactId).Delete(models.AlertContactInformation{})
	//同步region
	mq.SendMsg(config.GetRocketmqConfig().AlertContactTopic, enums.DeleteAlertContactInformation, contactId)
}

//删除联系人组关联
func (mpd *AlertContactDao) deleteAlertContactGroupRel(contactId string) {
	global.DB.Where("contact_id = ?", contactId).Delete(models.AlertContactGroupRel{})
	//同步region
	mq.SendMsg(config.GetRocketmqConfig().AlertContactTopic, enums.DeleteAlertContactGroupRelByContactId, contactId)
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
	var noticeMsgDTO = dtos.NoticeMsgDTO{}
	if noType == 1 {
		noticeMsgDTO.MsgEvent = dtos.MsgEvent{
			Type:   1, //TODO 枚举
			Source: dtos.VERIFY,
		}
	} else {
		noticeMsgDTO.MsgEvent = dtos.MsgEvent{
			Type:   2, //TODO 枚举
			Source: dtos.VERIFY,
		}
	}
	noticeMsgDTO.TenantId = tenantId
	noticeMsgDTO.SourceId = "activation-" + contactId
	var recvObjectBean = dtos.RecvObjectBean{
		RecvObject:     no,
		RecvObjectType: 2,
		NoticeContent:  tools.ToString(params),
	}
	noticeMsgDTO.RevObjectBean = recvObjectBean
	var noticeMsgDTOList []*dtos.NoticeMsgDTO
	noticeMsgDTOList = append(noticeMsgDTOList, &noticeMsgDTO)
	//TODO
	//service.NewMessageService(dao.NotificationRecord).SendMsg(noticeMsgDTOList, true)
}
