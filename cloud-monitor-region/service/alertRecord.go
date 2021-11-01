package service

import (
	commonDao "code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	commonDtos "code.cestc.cn/ccos-ops/cloud-monitor/business-common/dtos"
	commonModels "code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/dtos"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/forms"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// TODO 环境变量注入
var regionName string

type AlertRecordService struct {
	dao           *commonDao.AlertRecordDao
	tenantService *TenantService
}

func NewAlertRecordService(alertRecordDao *commonDao.AlertRecordDao, service *TenantService) *AlertRecordService {
	return &AlertRecordService{dao: alertRecordDao, tenantService: service}
}

func (s *AlertRecordService) AddAlertRecord(f *forms.InnerAlertRecordAddForm) error {

	if f.Alerts == nil || len(f.Alerts) <= 0 {
		log.Println("alerts 信息为空")
		return nil
	}

	var records []*commonModels.AlertRecord

	for _, alert := range f.Alerts {
		ruleDesc := &commonDtos.RuleDesc{}
		if err := json.Unmarshal([]byte(alert.Annotations.Description), ruleDesc); err != nil {
			return errors.New("序列化告警数据失败")
		}
		//	判断该条告警对应的规则是否删除、禁用、解绑
		if num := s.dao.FindAlertRuleBindNum(ruleDesc.InstanceId, ruleDesc.RuleId); num <= 0 {
			log.Println("此告警规则已删除/禁用/解绑")
			continue
		}
		// 构建告警记录列表、通知列表
		alertRecord := buildAlertRecord(alert, ruleDesc)
		records = append(records, alertRecord)

	}

	//todo

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

func (s *AlertRecordService) buildMsg(alert *forms.AlertRecordAlertsBean, ruleDesc *commonDtos.RuleDesc, alertRecord *commonModels.AlertRecord) *[]*dtos.NoticeMsg {
	if ruleDesc.GroupList == nil || len(ruleDesc.GroupList) <= 0 {
		log.Println("告警组为空")
		return nil
	}
	contactGroups := s.dao.FindContactInfoByGroupIds(ruleDesc.GroupList)
	var newContactGroups []commonDtos.ContactGroupInfo
	for _, group := range *contactGroups {
		if group.Contacts != nil && len(group.Contacts) > 0 {
			newContactGroups = append(newContactGroups, group)
		}
	}
	//返回联系人json信息注入到告警记录中
	alertRecord.ContactInfo = tools.ToString(contactGroups)

	instance := s.dao.FindFirstInstanceInfo(ruleDesc.InstanceId)
	if instance == nil {
		log.Println("未查询到实例信息")
		return nil
	}
	instanceInfo := s.buildInstanceInfo(instance, alert.Annotations.Summary)

	newRegionName := regionName
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
	params["userName"] = s.tenantService.GetTenantInfo(ruleDesc.TenantId).Name

	cv := fmt.Sprintf("%.2f", alertRecord.CurrentValue)
	unit := ruleDesc.Unit
	params["currentValue"] = cv + unit

	//channel := ruleDesc.NotifyChannel
	//TODO
	//String notifyChannel = descDTO.getNotifyChannel();
	//            List<NoticeMsgDTO> noticeMsgDTOS = buildSendMsg(alertRecord, contactInfoList, params, notifyChannel,descDTO);
	//            return noticeMsgDTOS;
	return nil
}

//func buildSendMsg(alertRecord models.AlertRecord, contactInfoList *[]*RecordContactInfo) *[]*dtos.NoticeMsg {
//
//}

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
		Region:       "",
		NoticeStatus: "",
		ContactInfo:  "",
		CreateTime:   now.Format("2006-01-02 15:04:05"),
		UpdateTime:   now.Format("2006-01-02 15:04:05"),
	}
}
