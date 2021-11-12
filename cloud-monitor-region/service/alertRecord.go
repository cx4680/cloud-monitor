package service

import (
	commonDao "code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	commonDtos "code.cestc.cn/ccos-ops/cloud-monitor/business-common/dtos"
	commonModels "code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/dtos"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/mq"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/utils"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type AlertRecordService struct {
	AlertRecordDao *commonDao.AlertRecordDao
	TenantService  *TenantService
	MessageService *MessageService
}

func NewAlertRecordService(alertRecordDao *commonDao.AlertRecordDao, tenantService *TenantService, messageService *MessageService) *AlertRecordService {
	return &AlertRecordService{AlertRecordDao: alertRecordDao, TenantService: tenantService, MessageService: messageService}
}

func (s *AlertRecordService) RecordAlertAndSendMessage(f *forms.InnerAlertRecordAddForm) error {
	ruleList := f.Alerts

	if ruleList == nil || len(ruleList) <= 0 {
		log.Println("alerts 信息为空")
		return nil
	}
	var list []*commonModels.AlertRecord
	var msgDTOList []*dtos.NoticeMsgDTO

	for _, alert := range ruleList {
		ruleDesc := &commonDtos.RuleDesc{}

		tools.ToObject(alert.Annotations.Description, ruleDesc)
		if ruleDesc == nil {
			return errors.New("序列化告警数据失败")
		}
		//	判断该条告警对应的规则是否删除、禁用、解绑
		if num := s.AlertRecordDao.FindAlertRuleBindNum(ruleDesc.InstanceId, ruleDesc.RuleId); num <= 0 {
			log.Println("此告警规则已删除/禁用/解绑")
			continue
		}
		// 构建告警记录列表、通知列表
		alertRecord := buildAlertRecord(alert, ruleDesc)

		noticeMsgDTOS := s.buildMsg(alert, ruleDesc, alertRecord)
		if noticeMsgDTOS != nil && len(noticeMsgDTOS) > 0 {
			msgDTOList = append(msgDTOList, noticeMsgDTOS...)
		}
		if noticeMsgDTOS != nil && len(noticeMsgDTOS) > 0 {
			alertRecord.NoticeStatus = "success"
		} else {
			alertRecord.NoticeStatus = "fail"
		}
		list = append(list, alertRecord)
	}
	s.AlertRecordDao.InsertBatch(list)
	if config.GetCommonConfig().HasNoticeModel {
		notificationRecords := s.MessageService.SendMsg(msgDTOList, false)
		mq.SendNotificationRecordMsg(notificationRecords)
	}
	mq.SendAlertRecordMsg(list)
	return nil
}

func getLabelMap(summary string) map[string]string {
	var labelMap = make(map[string]string)
	for _, s := range strings.Split(summary, ",") {
		labels := strings.Split(s, "=")
		if len(labels) == 2 {
			labelMap[labels[0]] = labels[1]
		}
	}
	return labelMap
}

func parseDateStr(source string) string {
	parse, _ := time.Parse("2006-01-02T15:04:05Z", source)
	return time.Unix(parse.Unix(), 0).Format("2006-01-02 15:04:05")
}

func (s *AlertRecordService) buildInstanceInfo(instance *commonModels.AlarmInstance, summary string) string {
	var builder strings.Builder
	if tools.IsNotBlank(instance.InstanceName) {
		builder.WriteString(instance.InstanceName)
		builder.WriteString("/")
	}
	if tools.IsNotBlank(instance.Ip) {
		builder.WriteString(instance.Ip)
	}
	labelMap := getLabelMap(summary)
	delete(labelMap, "instance")
	delete(labelMap, "name")
	delete(labelMap, "currentValue")

	for _, value := range labelMap {
		builder.WriteString("/")
		builder.WriteString(value)
	}

	return builder.String()
}

func (s *AlertRecordService) buildMsg(alert *forms.AlertRecordAlertsBean, ruleDesc *commonDtos.RuleDesc, alertRecord *commonModels.AlertRecord) []*dtos.NoticeMsgDTO {
	if ruleDesc.GroupList == nil || len(ruleDesc.GroupList) <= 0 {
		log.Println("告警组为空")
		return nil
	}
	contactGroups := s.AlertRecordDao.FindContactInfoByGroupIds(ruleDesc.GroupList)
	var newContactGroups []*commonDtos.ContactGroupInfo
	for _, group := range contactGroups {
		if group.Contacts != nil && len(group.Contacts) > 0 {
			newContactGroups = append(newContactGroups, &group)
		}
	}
	//返回联系人json信息注入到告警记录中
	alertRecord.ContactInfo = tools.ToString(contactGroups)

	instance := s.AlertRecordDao.FindFirstInstanceInfo(ruleDesc.InstanceId)
	if instance == nil {
		log.Println("未查询到实例信息")
		return nil
	}
	instanceInfo := s.buildInstanceInfo(instance, alert.Annotations.Summary)

	newRegionName := config.GetCommonConfig().RegionName
	if tools.IsNotBlank(instance.RegionName) {
		newRegionName = instance.RegionName
	}

	params := make(map[string]string)
	params["duration"] = alertRecord.Duration
	params["instanceInfo"] = instanceInfo
	params["product"] = ruleDesc.Product
	params["regionName"] = newRegionName
	params["metricName"] = ruleDesc.MonitorItem
	params["Name"] = ruleDesc.RuleName
	params["userName"] = s.TenantService.GetTenantInfo(ruleDesc.TenantId).Name

	cv := fmt.Sprintf("%.2f", alertRecord.CurrentValue)
	unit := ruleDesc.Unit
	params["currentValue"] = cv + unit

	channel := ruleDesc.NotifyChannel
	return buildSendMsg(alertRecord, newContactGroups, params, channel, ruleDesc)
}

func buildSendMsg(alertRecord *commonModels.AlertRecord, contactInfoList []*commonDtos.ContactGroupInfo, params map[string]string, notifyChannel string, ruleDesc *commonDtos.RuleDesc) []*dtos.NoticeMsgDTO {
	msgSource := dtos.MsgSourceDTO{
		Type:     dtos.ALERT_OPEN,
		SourceId: alertRecord.Id,
	}
	if "resolved" == alertRecord.Status {
		msgSource.Type = dtos.ALERT_CANCEL
		params["recoveryTime"] = alertRecord.EndTime
	} else {
		params["alertTime"] = alertRecord.StartTime
	}
	var list []*dtos.NoticeMsgDTO
	for _, recordContactInfo := range contactInfoList {
		for _, userContactInfo := range recordContactInfo.Contacts {

			if (notifyChannel == "email" || notifyChannel == "all") && tools.IsNotBlank(userContactInfo.Mail) {
				for _, email := range strings.Split(userContactInfo.Mail, ",") {
					list = append(list, buildMsgDTO(&msgSource, alertRecord, params, dtos.Email, email, ruleDesc))
				}
			}

			if (notifyChannel == "phone" || notifyChannel == "all") && tools.IsNotBlank(userContactInfo.Phone) {
				for _, phone := range strings.Split(userContactInfo.Mail, ",") {
					list = append(list, buildMsgDTO(&msgSource, alertRecord, params, dtos.Phone, phone, ruleDesc))
				}

			}

		}
	}
	return list

}

func getTime(time int) string {
	return "持续" + string(rune(time)) + "个周期"
}

func buildMsgDTO(msgSource *dtos.MsgSourceDTO, alertRecord *commonModels.AlertRecord, params map[string]string, receiverType dtos.ReceiveType, no string, descDTO *commonDtos.RuleDesc) *dtos.NoticeMsgDTO {

	objMap := make(map[string]string)
	if dtos.Email == receiverType {
		objMap["instanceAmount"] = "1"
		objMap["period"] = utils.GetDateDiff(descDTO.Time)
		objMap["times"] = getTime(descDTO.Time)
		objMap["time"] = utils.GetDateDiff(descDTO.Period * 1000)
		objMap["calType"] = descDTO.Statistic
		objMap["expression"] = string(descDTO.ComparisonOperator)
		if receiverType == dtos.Email && msgSource.Type == dtos.ALERT_CANCEL {
			targetValue := ""
			if tools.IsNotBlank(descDTO.Unit) {
				targetValue = descDTO.Unit
			}
			targetValue = string(rune(descDTO.TargetValue)) + targetValue
			objMap["targetValue"] = targetValue
		} else if receiverType == dtos.Email && msgSource.Type == dtos.ALERT_OPEN {
			objMap["targetValue"] = string(rune(descDTO.TargetValue))
		}

		objMap["evaluationCount"] = getTime(descDTO.Time)
		objMap["InstanceID"] = descDTO.InstanceId
		if tools.IsBlank(descDTO.Unit) {
			objMap["unit"] = ""
		} else {
			objMap["unit"] = descDTO.Unit
		}
	}

	return &dtos.NoticeMsgDTO{
		SourceId: msgSource.SourceId,
		TenantId: alertRecord.TenantId,
		MsgEvent: dtos.MsgEvent{
			Type:   receiverType,
			Source: msgSource.Type,
		},
		RevObjectBean: dtos.RecvObjectBean{
			RecvObjectType: receiverType,
			RecvObject:     no,
			NoticeContent:  tools.ToString(objMap),
		},
	}
}

func buildAlertRecord(alert *forms.AlertRecordAlertsBean, ruleDesc *commonDtos.RuleDesc) *commonModels.AlertRecord {
	summary, err := json.Marshal(alert.Annotations)
	if err != nil {
		log.Println(err)
		summary = []byte{}
	}
	labelMap := getLabelMap(string(summary))

	now := time.Now()
	startTimeStr := parseDateStr(alert.StartsAt)
	startTime, _ := time.Parse("2016-01-02 15:04:05", startTimeStr)
	duration := now.Second() - startTime.Second()

	return &commonModels.AlertRecord{
		Status:       alert.Status,
		TenantId:     ruleDesc.TenantId,
		RuleId:       ruleDesc.RuleId,
		RuleName:     ruleDesc.RuleName,
		MonitorType:  ruleDesc.MonitorType,
		SourceType:   ruleDesc.Product,
		SourceId:     ruleDesc.InstanceId,
		Summary:      string(summary),
		CurrentValue: labelMap["currentValue"],
		StartTime:    startTimeStr,
		EndTime:      parseDateStr(alert.EndsAt),
		TargetValue:  strconv.Itoa(ruleDesc.TargetValue),
		Expression:   ruleDesc.Express,
		Duration:     strconv.Itoa(duration),
		Level:        ruleDesc.Level,
		AlarmKey:     ruleDesc.MetricName,
		Region:       config.GetCommonConfig().RegionName,
		NoticeStatus: "",
		ContactInfo:  "",
		CreateTime:   now.Format("2006-01-02 15:04:05"),
		UpdateTime:   now.Format("2006-01-02 15:04:05"),
	}
}
