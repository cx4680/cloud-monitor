package service

import (
	commonDao "code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	commonDtos "code.cestc.cn/ccos-ops/cloud-monitor/business-common/dtos"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysComponent/sysRocketMq"
	commonModels "code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	commonService "code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service/external/messageCenter"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/utils"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/utils/snowflake"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type filter func(*forms.AlertRecordAlertsBean) (bool, error)

type AlertRecordAddService struct {
	FilterChain     []filter
	AlertRecordSvc  *AlertRecordService
	AlarmHandlerSvc *commonService.AlarmHandlerService
	MessageSvc      *service.MessageService
	TenantSvc       *service.TenantService
	AlertRecordDao  *commonDao.AlertRecordDao
}

func NewAlertRecordAddService(AlertRecordSvc *AlertRecordService, AlarmHandlerSvc *commonService.AlarmHandlerService, MessageSvc *service.MessageService, TenantSvc *service.TenantService) *AlertRecordAddService {
	return &AlertRecordAddService{
		FilterChain: []filter{func(alert *forms.AlertRecordAlertsBean) (bool, error) {
			ruleDesc := &commonDtos.RuleDesc{}
			tools.ToObject(alert.Annotations.Description, ruleDesc)
			if ruleDesc == nil {
				return false, errors.New("序列化告警数据失败")
			}
			//	判断该条告警对应的规则是否删除、禁用、解绑
			if num := commonDao.AlertRecord.FindAlertRuleBindNum(ruleDesc.InstanceId, ruleDesc.RuleId); num <= 0 {
				logger.Logger().Info("此告警规则已删除/禁用/解绑")
				return false, nil
			}
			return true, nil
		}},
		AlarmHandlerSvc: AlarmHandlerSvc,
		AlertRecordSvc:  AlertRecordSvc,
		MessageSvc:      MessageSvc,
		TenantSvc:       TenantSvc,
		AlertRecordDao:  commonDao.AlertRecord,
	}
}

func (s *AlertRecordAddService) Add(f forms.InnerAlertRecordAddForm) error {
	if len(f.Alerts) == 0 {
		logger.Logger().Info("alerts 信息为空")
		return nil
	}
	list, handlerMap := s.checkAndBuild(f.Alerts)
	//持久化
	if list != nil && len(list) > 0 {
		if err := s.AlertRecordSvc.Persistence(s.AlertRecordSvc, sysRocketMq.RecordTopic, list); err != nil {
			return err
		}
	}
	//告警处置
	for t, p := range handlerMap {
		if p == nil || len(p) <= 0 {
			continue
		}
		if t == 1 || t == 2 {
			if err := s.MessageSvc.SendMsg(p, false); err != nil {
				return err
			}
		} else if t == 3 {
			//调用弹性伸缩
			for _, a := range p {
				data := a.(*commonDtos.AutoScalingData)
				respJson, err := tools.HttpPostJson(data.Param+"/inner/as/trigger", map[string]string{"ruleId": data.RuleId, "tenantId": data.TenantId}, nil)
				if err != nil {
					logger.Logger().Error("autoScaling request fail,", err)
				} else {
					logger.Logger().Info("autoScaling request success, resp=", respJson)
				}
			}
		}
	}

	return nil
}

func (s *AlertRecordAddService) checkAndBuild(alerts []*forms.AlertRecordAlertsBean) ([]commonModels.AlertRecord, map[int][]interface{}) {
	var list []commonModels.AlertRecord

	handlerMap := map[int][]interface{}{}

	for _, alert := range alerts {
		ret, err := s.predicate(alert)
		if err != nil {
			logger.Logger().Error("filter check error, alert=", tools.ToString(alert), err)
			continue
		}
		if !ret {
			logger.Logger().Info("filter check false, alert=", tools.ToString(alert))
			continue
		}
		//获取告警规则信息
		ruleDesc := &commonDtos.RuleDesc{}
		tools.ToObject(alert.Annotations.Description, ruleDesc)
		if ruleDesc == nil {
			logger.Logger().Error("序列化告警数据失败, ", tools.ToString(alert))
			continue
		}
		//获取告警联系人信息
		var contactGroups []commonDtos.ContactGroupInfo
		if ruleDesc.GroupList != nil && len(ruleDesc.GroupList) > 0 {
			contactGroups = s.AlertRecordDao.FindContactInfoByGroupIds(ruleDesc.GroupList)
		}

		record := s.buildAlertRecord(alert, ruleDesc, contactGroups)
		if record == nil {
			continue
		}
		list = append(list, *record)
		//告警处理
		//获取告警处理方式列表
		handlerList := s.AlarmHandlerSvc.GetAlarmHandlerListByRuleId(ruleDesc.RuleId)

		if handlerList != nil && len(handlerList) > 0 {
			for _, handler := range handlerList {
				if handlerMap[handler.HandleType] == nil {
					var pendingList []interface{}
					handlerMap[handler.HandleType] = pendingList
				}
				var data interface{}
				switch handler.HandleType {
				case 1:
				case 2:
					data = s.buildNoticeData(alert, record, ruleDesc, contactGroups, handler)
				case 3:
					data = s.buildAutoScalingData(alert, ruleDesc, handler.HandleParams)
				}
				if data != nil {
					handlerMap[handler.HandleType] = append(handlerMap[handler.HandleType], data)
				}
			}
		}
	}
	return list, handlerMap
}

func (s *AlertRecordAddService) buildAlertRecord(alert *forms.AlertRecordAlertsBean, ruleDesc *commonDtos.RuleDesc, contactGroups []commonDtos.ContactGroupInfo) *commonModels.AlertRecord {
	now := tools.GetNow()
	startTime := tools.TimeParseForZone(alert.StartsAt)
	//持续时间，单位秒
	duration := now.Second() - startTime.Second()

	alertAnnotationStr := tools.ToString(alert.Annotations)
	if tools.IsBlank(alertAnnotationStr) {
		logger.Logger().Info("告警数据为空, ", tools.ToString(alert))
		return nil
	}
	labelMap := s.getLabelMap(alertAnnotationStr)

	return &commonModels.AlertRecord{
		Id:           strconv.FormatInt(snowflake.GetWorker().NextId(), 10),
		Status:       alert.Status,
		TenantId:     ruleDesc.TenantId,
		RuleId:       ruleDesc.RuleId,
		RuleName:     ruleDesc.RuleName,
		MonitorType:  ruleDesc.MonitorType,
		SourceType:   ruleDesc.Product,
		SourceId:     ruleDesc.InstanceId,
		Summary:      alertAnnotationStr,
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
}

func (s *AlertRecordAddService) buildAutoScalingData(alert *forms.AlertRecordAlertsBean, ruleDesc *commonDtos.RuleDesc, param string) *commonDtos.AutoScalingData {
	if "resolved" == alert.Status {
		//告警恢复不需要通知弹性伸缩
		return nil
	}
	if tools.IsBlank(ruleDesc.ResourceGroupId) {
		//脏数据
		return nil
	}

	return &commonDtos.AutoScalingData{
		TenantId:        ruleDesc.TenantId,
		RuleId:          ruleDesc.RuleId,
		ResourceGroupId: ruleDesc.ResourceGroupId,
		Param:           param,
	}
}

func (s *AlertRecordAddService) buildNoticeData(alert *forms.AlertRecordAlertsBean, record *commonModels.AlertRecord, ruleDesc *commonDtos.RuleDesc, contactGroups []commonDtos.ContactGroupInfo, handler commonModels.AlarmHandler) *service.AlertMsgSendDTO {
	source := messageCenter.ALERT_OPEN
	if "resolved" == alert.Status {
		source = messageCenter.ALERT_CANCEL
	}

	handleType := handler.HandleType
	targetList := s.buildTargets(handleType, contactGroups)
	instanceInfo := s.getInstanceInfo(alert.Annotations.Summary)

	objMap := make(map[string]string)
	objMap["duration"] = record.Duration
	objMap["instanceInfo"] = instanceInfo
	objMap["product"] = ruleDesc.Product
	objMap["regionName"] = record.Region
	objMap["metricName"] = ruleDesc.MonitorItem
	objMap["Name"] = ruleDesc.RuleName
	objMap["userName"] = s.TenantSvc.GetTenantInfo(ruleDesc.TenantId).Name

	f, _ := strconv.ParseFloat(record.CurrentValue, 64)
	cv := fmt.Sprintf("%.2f", f)
	objMap["currentValue"] = cv + ruleDesc.Unit

	if 1 == handleType {
		objMap["instanceAmount"] = "1"
		objMap["period"] = utils.GetDateDiff(ruleDesc.Time)
		objMap["times"] = "持续" + strconv.Itoa(ruleDesc.Time) + "个周期"
		objMap["time"] = utils.GetDateDiff(ruleDesc.Period * 1000)
		objMap["calType"] = ruleDesc.Statistic
		objMap["expression"] = ruleDesc.ComparisonOperator
		if source == messageCenter.ALERT_CANCEL {
			targetValue := ""
			if tools.IsNotBlank(ruleDesc.Unit) {
				targetValue = ruleDesc.Unit
			}
			targetValue = strconv.Itoa(ruleDesc.TargetValue) + targetValue
			objMap["targetValue"] = targetValue
		} else if source == messageCenter.ALERT_OPEN {
			objMap["targetValue"] = strconv.Itoa(ruleDesc.TargetValue)
		}

		objMap["evaluationCount"] = "持续" + strconv.Itoa(ruleDesc.Time) + "个周期"
		objMap["InstanceID"] = ruleDesc.InstanceId
		if tools.IsBlank(ruleDesc.Unit) {
			objMap["unit"] = ""
		} else {
			objMap["unit"] = ruleDesc.Unit
		}
	}

	return &service.AlertMsgSendDTO{
		AlertId: record.Id,
		Msg: messageCenter.MessageSendDTO{
			SenderId:   ruleDesc.TenantId,
			Target:     targetList,
			SourceType: source,
			Content:    tools.ToString(objMap),
		},
	}
}

func (s *AlertRecordAddService) getInstanceInfo(summary string) string {
	var builder strings.Builder

	labelMap := s.getLabelMap(summary)
	instance := s.AlertRecordDao.FindFirstInstanceInfo(labelMap["instance"])

	if tools.IsNotBlank(instance.InstanceName) {
		builder.WriteString(instance.InstanceName)
		builder.WriteString("/")
	}
	if tools.IsNotBlank(instance.Ip) {
		builder.WriteString(instance.Ip)
	}

	delete(labelMap, "instance")
	delete(labelMap, "name")
	delete(labelMap, "currentValue")

	for _, value := range labelMap {
		builder.WriteString("/")
		builder.WriteString(value)
	}

	return builder.String()
}

func (s *AlertRecordAddService) buildTargets(channel int, contactGroups []commonDtos.ContactGroupInfo) []messageCenter.MessageTargetDTO {
	var targetList []messageCenter.MessageTargetDTO

	for _, group := range contactGroups {
		for _, contact := range group.Contacts {
			if 1 == channel {
				if mailTargetList := s.buildMailTargets(contact); mailTargetList != nil {
					targetList = append(targetList, mailTargetList...)
				}
			} else if 2 == channel {
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

func (s *AlertRecordAddService) getLabelMap(summary string) map[string]string {
	var labelMap = make(map[string]string)
	for _, s := range strings.Split(summary, ",") {
		labels := strings.Split(s, "=")
		if len(labels) == 2 {
			labelMap[labels[0]] = labels[1]
		}
	}
	return labelMap
}

func (s *AlertRecordAddService) predicate(alert *forms.AlertRecordAlertsBean) (bool, error) {
	for _, f := range s.FilterChain {
		ret, err := f(alert)
		if err != nil {
			return false, err
		}
		if !ret {
			return ret, nil
		}

	}
	return true, nil
}
