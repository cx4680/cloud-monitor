package service

import (
	commonDao "code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	commonDtos "code.cestc.cn/ccos-ops/cloud-monitor/business-common/dtos"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysComponent/sysRocketMq"
	commonModels "code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service/external/messageCenter"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/utils"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/utils/snowflake"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Filter interface {
	DoFilter(alert *forms.AlertRecordAlertsBean) (bool, error)
}

type AlertRecordAddService struct {
	FilterChain    []Filter
	AlertRecordSvc *AlertRecordService
	MessageSvc     *service.MessageService
	TenantSvc      *service.TenantService
	AlertRecordDao *commonDao.AlertRecordDao
}

func NewAlertRecordAddService(AlertRecordSvc *AlertRecordService, MessageSvc *service.MessageService, TenantSvc *service.TenantService) *AlertRecordAddService {
	return &AlertRecordAddService{
		FilterChain: []Filter{&RuleChangeFilter{
			AlertRecordDao: commonDao.AlertRecord,
		}},
		AlertRecordSvc: AlertRecordSvc,
		MessageSvc:     MessageSvc,
		TenantSvc:      TenantSvc,
		AlertRecordDao: commonDao.AlertRecord,
	}
}

func (s *AlertRecordAddService) Add(f forms.InnerAlertRecordAddForm) error {
	if len(f.Alerts) == 0 {
		log.Println("alerts 信息为空")
		return nil
	}
	list, alertMsgList := s.checkAndBuild(f.Alerts)
	//持久化
	if list != nil && len(list) > 0 {
		if err := s.persistence(list); err != nil {
			return err
		}
	}

	//发送通知
	if alertMsgList != nil && len(alertMsgList) > 0 {
		if err := s.sendNotice(alertMsgList); err != nil {
			return err
		}
	}
	return nil
}

func (s *AlertRecordAddService) checkAndBuild(alerts []*forms.AlertRecordAlertsBean) ([]commonModels.AlertRecord, []service.AlertMsgSendDTO) {
	var list []commonModels.AlertRecord
	var alertMsgList []service.AlertMsgSendDTO
	for _, alert := range alerts {
		ret, err := s.predicate(alert)
		if err != nil {
			log.Printf("filter check error, alert=%s,  %v\n", tools.ToString(alert), err)
			continue
		}
		if !ret {
			log.Printf("filter check false, alert=%s\n", tools.ToString(alert))
			continue
		}
		alertAnnotations := tools.ToString(alert.Annotations)
		if tools.IsBlank(alertAnnotations) {
			log.Printf("告警数据为空, %s\n", tools.ToString(alert))
			continue
		}
		labelMap := getLabelMap(alertAnnotations)
		now := tools.GetNow()
		startTime := tools.TimeParseForZone(alert.StartsAt)
		//持续时间，单位秒
		duration := now.Second() - startTime.Second()

		ruleDesc := &commonDtos.RuleDesc{}
		tools.ToObject(alert.Annotations.Description, ruleDesc)
		if ruleDesc == nil {
			log.Printf("序列化告警数据失败, %s\n", tools.ToString(alert))
			continue
		}
		var contactGroups []commonDtos.ContactGroupInfo
		if ruleDesc.GroupList != nil && len(ruleDesc.GroupList) > 0 {
			contactGroups = s.AlertRecordDao.FindContactInfoByGroupIds(ruleDesc.GroupList)
		}
		record := commonModels.AlertRecord{
			Id:           strconv.FormatInt(snowflake.GetWorker().NextId(), 10),
			Status:       alert.Status,
			TenantId:     ruleDesc.TenantId,
			RuleId:       ruleDesc.RuleId,
			RuleName:     ruleDesc.RuleName,
			MonitorType:  ruleDesc.MonitorType,
			SourceType:   ruleDesc.Product,
			SourceId:     ruleDesc.InstanceId,
			Summary:      alertAnnotations,
			CurrentValue: labelMap["currentValue"],
			StartTime:    tools.TimeToStr(startTime, tools.FullTimeFmt),
			EndTime:      tools.TimeToStr(tools.TimeParseForZone(alert.EndsAt), tools.FullTimeFmt),
			TargetValue:  strconv.Itoa(ruleDesc.TargetValue),
			Expression:   ruleDesc.Express,
			Duration:     strconv.Itoa(duration),
			Level:        ruleDesc.Level,
			AlarmKey:     ruleDesc.MetricName,
			Region:       config.GetCommonConfig().RegionName,
			NoticeStatus: "success",
			ContactInfo:  tools.ToString(contactGroups),
			CreateTime:   tools.TimeToStr(now, tools.FullTimeFmt),
			UpdateTime:   tools.TimeToStr(now, tools.FullTimeFmt),
		}
		list = append(list, record)

		channel := ruleDesc.NotifyChannel
		targetList := s.buildTargets(channel, contactGroups)

		st := messageCenter.ALERT_OPEN
		if "resolved" == alert.Status {
			st = messageCenter.ALERT_CANCEL
		}

		instance := s.AlertRecordDao.FindFirstInstanceInfo(ruleDesc.InstanceId)
		if instance == nil {
			log.Printf("未查询到实例信息, %s\n", tools.ToString(alert))
			continue
		}
		instanceInfo := s.getInstanceInfo(instance, alert.Annotations.Summary)

		objMap := make(map[string]string)

		objMap["duration"] = record.Duration
		objMap["instanceInfo"] = instanceInfo
		objMap["product"] = ruleDesc.Product
		objMap["regionName"] = record.Region
		objMap["metricName"] = ruleDesc.MonitorItem
		objMap["Name"] = ruleDesc.RuleName
		objMap["userName"] = s.TenantSvc.GetTenantInfo(ruleDesc.TenantId).Name

		cv := fmt.Sprintf("%.2f", record.CurrentValue)
		unit := ruleDesc.Unit
		objMap["currentValue"] = cv + unit

		if "mail" == channel {
			objMap["instanceAmount"] = "1"
			objMap["period"] = utils.GetDateDiff(ruleDesc.Time)
			objMap["times"] = getTime(ruleDesc.Time)
			objMap["time"] = utils.GetDateDiff(ruleDesc.Period * 1000)
			objMap["calType"] = ruleDesc.Statistic
			objMap["expression"] = ruleDesc.ComparisonOperator
			if st == messageCenter.ALERT_CANCEL {
				targetValue := ""
				if tools.IsNotBlank(ruleDesc.Unit) {
					targetValue = ruleDesc.Unit
				}
				targetValue = string(rune(ruleDesc.TargetValue)) + targetValue
				objMap["targetValue"] = targetValue
			} else if st == messageCenter.ALERT_OPEN {
				objMap["targetValue"] = string(rune(ruleDesc.TargetValue))
			}

			objMap["evaluationCount"] = getTime(ruleDesc.Time)
			objMap["InstanceID"] = ruleDesc.InstanceId
			if tools.IsBlank(ruleDesc.Unit) {
				objMap["unit"] = ""
			} else {
				objMap["unit"] = ruleDesc.Unit
			}
		}

		alertMsgList = append(alertMsgList, service.AlertMsgSendDTO{
			AlertId: record.Id,
			Msg: messageCenter.MessageSendDTO{
				SenderId:   ruleDesc.TenantId,
				Target:     targetList,
				SourceType: st,
				Content:    tools.ToString(objMap),
			},
		})
	}
	return list, alertMsgList
}

func (s *AlertRecordAddService) getInstanceInfo(instance *commonModels.AlarmInstance, summary string) string {
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

func (s *AlertRecordAddService) buildTargets(channel string, contactGroups []commonDtos.ContactGroupInfo) []messageCenter.MessageTargetDTO {
	var targetList []messageCenter.MessageTargetDTO

	for _, group := range contactGroups {
		for _, contact := range group.Contacts {
			if "all" == channel {
				if mailTargetList := s.buildMailTargets(contact); mailTargetList != nil {
					targetList = append(targetList, mailTargetList...)
				}
				if phoneTargetList := s.buildPhoneTargets(contact); phoneTargetList != nil {
					targetList = append(targetList, phoneTargetList...)
				}
			} else if "email" == channel {
				if mailTargetList := s.buildMailTargets(contact); mailTargetList != nil {
					targetList = append(targetList, mailTargetList...)
				}
			} else if "phone" == channel {
				if phoneTargetList := s.buildPhoneTargets(contact); phoneTargetList != nil {
					targetList = append(targetList, phoneTargetList...)
				}
			}

		}
	}
	return targetList
}

func (s *AlertRecordAddService) buildMailTargets(contact commonDtos.UserContactInfo) []messageCenter.MessageTargetDTO {
	if tools.IsBlank(contact.Mail) {
		return nil
	}
	var list []messageCenter.MessageTargetDTO
	for _, m := range strings.Split(contact.Mail, ",") {
		list = append(list, messageCenter.MessageTargetDTO{
			Addr: m,
			Type: messageCenter.Email,
		})
	}
	return list
}

func (s *AlertRecordAddService) buildPhoneTargets(contact commonDtos.UserContactInfo) []messageCenter.MessageTargetDTO {
	if tools.IsBlank(contact.Phone) {
		return nil
	}
	var list []messageCenter.MessageTargetDTO
	for _, m := range strings.Split(contact.Phone, ",") {
		list = append(list, messageCenter.MessageTargetDTO{
			Addr: m,
			Type: messageCenter.Phone,
		})
	}
	return list
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

func getTime(time int) string {
	return "持续" + string(rune(time)) + "个周期"
}

func (s *AlertRecordAddService) predicate(alert *forms.AlertRecordAlertsBean) (bool, error) {
	for _, filter := range s.FilterChain {
		ret, err := filter.DoFilter(alert)
		if err != nil {
			return false, err
		}
		if !ret {
			return ret, nil
		}

	}
	return true, nil
}

func (s *AlertRecordAddService) persistence(list []commonModels.AlertRecord) error {
	return s.AlertRecordSvc.Persistence(s.AlertRecordSvc, sysRocketMq.RecordTopic, list)
}

func (s *AlertRecordAddService) sendNotice(alertMsgList []service.AlertMsgSendDTO) error {
	return s.MessageSvc.SendMsg(alertMsgList, false)
}

type RuleChangeFilter struct {
	AlertRecordDao *commonDao.AlertRecordDao
}

func (f *RuleChangeFilter) DoFilter(alert *forms.AlertRecordAlertsBean) (bool, error) {
	ruleDesc := &commonDtos.RuleDesc{}
	tools.ToObject(alert.Annotations.Description, ruleDesc)
	if ruleDesc == nil {
		return false, errors.New("序列化告警数据失败")
	}
	//	判断该条告警对应的规则是否删除、禁用、解绑
	if num := f.AlertRecordDao.FindAlertRuleBindNum(ruleDesc.InstanceId, ruleDesc.RuleId); num <= 0 {
		log.Println("此告警规则已删除/禁用/解绑")
		return false, nil
	}
	return true, nil
}
