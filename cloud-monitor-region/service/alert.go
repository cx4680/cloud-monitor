package service

import (
	commonDao "code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	commonDtos "code.cestc.cn/ccos-ops/cloud-monitor/business-common/dtos"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/enums/handlerType"
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
	"time"
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
			// 判断该告警对应的规则是否有变化 资源组
			if tools.IsNotBlank(ruleDesc.ResourceGroupId) {
				if num := commonDao.AlertRecord.FindAlertRuleBindGroupNum(ruleDesc.RuleId, ruleDesc.ResourceGroupId); num <= 0 {
					logger.Logger().Info("此告警规则已删除/禁用/解绑")
					return false, nil
				}
			} else if tools.IsNotBlank(ruleDesc.ResourceId) {
				if num := commonDao.AlertRecord.FindAlertRuleBindResourceNum(ruleDesc.RuleId, ruleDesc.ResourceId); num <= 0 {
					logger.Logger().Info("此告警规则已删除/禁用/解绑")
					return false, nil
				}
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
	logger.Logger().Info("alarm data:", tools.ToString(list), " || ", tools.ToString(handlerMap))
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
		if t == handlerType.Sms || t == handlerType.Email {
			if err := s.MessageSvc.SendAlarmNotice(p); err != nil {
				logger.Logger().Error("send alarm message fail,", err)
			}
		} else if t == handlerType.Http {
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
		var contactGroups []*commonDtos.ContactGroupInfo
		if ruleDesc.GroupList != nil && len(ruleDesc.GroupList) > 0 {
			contactGroups = s.AlertRecordDao.FindContactInfoByGroupIds(ruleDesc.GroupList)
		}

		alertAnnotationStr := tools.ToString(alert.Annotations)
		if tools.IsBlank(alertAnnotationStr) {
			logger.Logger().Info("告警数据为空, ", tools.ToString(alert))
			continue
		}
		labelMap := s.getLabelMap(alert.Annotations.Summary)

		record := s.buildAlertRecord(alert, ruleDesc, contactGroups, labelMap)
		if record == nil {
			continue
		}
		list = append(list, *record)
		//告警处理
		//获取告警处理方式列表
		handlerList := s.AlarmHandlerSvc.GetAlarmHandlerListByRuleId(ruleDesc.RuleId)

		if handlerList == nil || len(handlerList) <= 0 {
			continue
		}

		for _, handler := range handlerList {
			if handlerMap[handler.HandleType] == nil {
				var pendingList []interface{}
				handlerMap[handler.HandleType] = pendingList
			}
			var data interface{}
			switch handler.HandleType {
			case handlerType.Email, handlerType.Sms:
				data = s.buildNoticeData(alert, record, ruleDesc, contactGroups, handler.HandleType)
			case handlerType.Http:
				data = s.buildAutoScalingData(alert, ruleDesc, handler.HandleParams)
			default:

			}
			if data != nil {
				handlerMap[handler.HandleType] = append(handlerMap[handler.HandleType], data)
			}
		}

	}
	return list, handlerMap
}

func (s *AlertRecordAddService) getDurationTime(now, startTime time.Time) string {
	//持续时间，单位秒
	sub := now.Sub(startTime)
	return utils.GetDateDiff(int(sub.Round(time.Second).Seconds()))
}

func (s *AlertRecordAddService) buildAlertRecord(alert *forms.AlertRecordAlertsBean, ruleDesc *commonDtos.RuleDesc, contactGroups []*commonDtos.ContactGroupInfo, labelMap map[string]string) *commonModels.AlertRecord {
	now := tools.GetNow()
	startTime := tools.TimeParseForZone(alert.StartsAt)

	val := strconv.FormatFloat(ruleDesc.TargetValue, 'f', 2, 64)

	sourceId := ruleDesc.ResourceId
	if tools.IsBlank(sourceId) {
		sourceId = ruleDesc.ResourceGroupId
	}

	return &commonModels.AlertRecord{
		Id:           strconv.FormatInt(snowflake.GetWorker().NextId(), 10),
		Status:       alert.Status,
		TenantId:     ruleDesc.TenantId,
		RuleId:       ruleDesc.RuleId,
		RuleName:     ruleDesc.RuleName,
		MonitorType:  ruleDesc.MonitorType,
		SourceType:   ruleDesc.Product,
		SourceId:     sourceId,
		Summary:      alert.Annotations.Summary,
		CurrentValue: labelMap["currentValue"],
		StartTime:    tools.TimeToStr(startTime, tools.FullTimeFmt),
		EndTime:      tools.TimeToStr(tools.TimeParseForZone(alert.EndsAt), tools.FullTimeFmt),
		TargetValue:  val,
		Expression:   ruleDesc.Express,
		Duration:     s.getDurationTime(now, startTime),
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
	//弹性伸缩一定有资源组数据
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

func (s *AlertRecordAddService) buildNoticeData(alert *forms.AlertRecordAlertsBean, record *commonModels.AlertRecord, ruleDesc *commonDtos.RuleDesc,
	contactGroups []*commonDtos.ContactGroupInfo, ht int) *service.AlertMsgSendDTO {
	source := messageCenter.ALERT_OPEN
	if "resolved" == alert.Status {
		source = messageCenter.ALERT_CANCEL
	}

	instanceInfo := s.getInstanceInfo(ruleDesc.ResourceId, alert.Annotations.Summary)

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

	val := strconv.FormatFloat(ruleDesc.TargetValue, 'f', 2, 64)

	if handlerType.Email == ht {
		objMap["instanceAmount"] = "1"
		objMap["period"] = utils.GetDateDiff(ruleDesc.Period * 1000)
		objMap["times"] = "持续" + strconv.Itoa(ruleDesc.Time) + "个周期"
		objMap["time"] = utils.GetDateDiff(ruleDesc.Period * 1000)
		objMap["calType"] = commonDao.GetAlarmStatisticsText(ruleDesc.Statistic)
		objMap["expression"] = commonDao.GetComparisonOperator(ruleDesc.ComparisonOperator)
		if source == messageCenter.ALERT_CANCEL {
			targetValue := ""
			if tools.IsNotBlank(ruleDesc.Unit) {
				targetValue = ruleDesc.Unit
			}
			targetValue = val + targetValue
			objMap["targetValue"] = targetValue
		} else if source == messageCenter.ALERT_OPEN {
			objMap["targetValue"] = val
		}

		objMap["evaluationCount"] = "持续" + strconv.Itoa(ruleDesc.Time) + "个周期"
		// 暂不支持资源组
		objMap["InstanceID"] = ruleDesc.ResourceId
		if tools.IsBlank(ruleDesc.Unit) {
			objMap["unit"] = ""
		} else {
			objMap["unit"] = ruleDesc.Unit
		}
	}

	var msgList []service.AlertMsgDTO

	//TODO
	if ht == handlerType.Sms {
		for _, group := range contactGroups {
			for _, contact := range group.Contacts {
				msgList = append(msgList, service.AlertMsgDTO{
					Type:    messageCenter.Phone,
					Targets: s.buildPhoneTargets(contact),
					Content: tools.ToString(objMap),
				})
			}
		}
	} else if ht == handlerType.Email {
		for _, group := range contactGroups {
			for _, contact := range group.Contacts {
				msgList = append(msgList, service.AlertMsgDTO{
					Type:    messageCenter.Email,
					Targets: s.buildMailTargets(contact),
					Content: tools.ToString(objMap),
				})
			}
		}
	}

	return &service.AlertMsgSendDTO{
		AlertId:    record.Id,
		SenderId:   ruleDesc.TenantId,
		SourceType: source,
		Msgs:       msgList,
	}
}

func (s *AlertRecordAddService) getInstanceInfo(resourceId, summary string) string {
	var builder strings.Builder

	labelMap := s.getLabelMap(summary)
	instance := &commonModels.AlarmInstance{}
	if tools.IsNotBlank(resourceId) {
		instance = s.AlertRecordDao.FindFirstInstanceInfo(resourceId)
	}

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

func (s *AlertRecordAddService) buildMailTargets(contact commonDtos.UserContactInfo) []string {
	if tools.IsBlank(contact.Mail) {
		return nil
	}
	var list []string
	for _, m := range strings.Split(contact.Mail, ",") {
		list = append(list, m)
	}
	return list
}

func (s *AlertRecordAddService) buildPhoneTargets(contact commonDtos.UserContactInfo) []string {
	if tools.IsBlank(contact.Phone) {
		return nil
	}
	var list []string
	for _, m := range strings.Split(contact.Phone, ",") {
		list = append(list, m)
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
