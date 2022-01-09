package service

import (
	commonDao "code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	commonDtos "code.cestc.cn/ccos-ops/cloud-monitor/business-common/dto"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/enum/handler_type"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sys_component/sys_rocketmq"
	commonModels "code.cestc.cn/ccos-ops/cloud-monitor/business-common/model"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	commonService "code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service/external/message_center"
	commonUtil "code.cestc.cn/ccos-ops/cloud-monitor/business-common/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/httputil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/snowflake"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type filter func(*form.AlertRecordAlertsBean) (bool, error)

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
		FilterChain: []filter{func(alert *form.AlertRecordAlertsBean) (bool, error) {
			rule := &commonDtos.RuleDesc{}
			if err := jsonutil.ToObjectWithError(alert.Annotations.Description, rule); err != nil {
				return false, errors.New("requestId=" + alert.RequestId + ", 序列化告警数据失败")
			}
			// 判断该告警对应的规则是否有变化 资源组
			if strutil.IsNotBlank(rule.ResourceGroupId) {
				if num := commonDao.AlertRecord.FindAlertRuleBindGroupNum(rule.RuleId, rule.ResourceGroupId); num <= 0 {
					logger.Logger().Info("requestId=" + alert.RequestId + ",此告警规则已删除/禁用/解绑")
					return false, nil
				}
			} else if strutil.IsNotBlank(rule.ResourceId) {
				if num := commonDao.AlertRecord.FindAlertRuleBindResourceNum(rule.RuleId, rule.ResourceId); num <= 0 {
					logger.Logger().Info("requestId=" + alert.RequestId + ",此告警规则已删除/禁用/解绑")
					return false, nil
				}
			}
			dbRule := commonDao.AlarmRule.GetById(global.DB, rule.RuleId)
			cond := dbRule.RuleCondition
			if !(cond.MetricName == rule.MetricName &&
				cond.Period == rule.Period &&
				cond.Times == rule.Time &&
				cond.Statistics == rule.Statistic &&
				cond.ComparisonOperator == rule.ComparisonOperator &&
				cond.Threshold == rule.TargetValue) {
				logger.Logger().Info("requestId=", alert.RequestId, "此告警触发条件改变")
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

func (s *AlertRecordAddService) Add(requestId string, f form.InnerAlertRecordAddForm) error {
	if len(f.Alerts) == 0 {
		logger.Logger().Info("requestId=", requestId, ", alerts 信息为空")
		return nil
	}
	list, handlerMap := s.checkAndBuild(requestId, f.Alerts)
	logger.Logger().Info("requestId=", requestId, ", alarm data=", jsonutil.ToString(list))
	logger.Logger().Info("requestId=", requestId, ", handler dat=", jsonutil.ToString(handlerMap))
	//持久化
	if list != nil && len(list) > 0 {
		if err := s.AlertRecordSvc.Persistence(s.AlertRecordSvc, sys_rocketmq.RecordTopic, list); err != nil {
			return err
		}
	}
	//告警处置
	for t, p := range handlerMap {
		if p == nil || len(p) <= 0 {
			continue
		}
		if t == handler_type.Sms || t == handler_type.Email {
			if err := s.MessageSvc.SendAlarmNotice(p); err != nil {
				logger.Logger().Error("requestId=", requestId, ", send alarm message fail,", err)
			}
		} else if t == handler_type.Http {
			//调用弹性伸缩
			for _, a := range p {
				data := a.(*commonDtos.AutoScalingData)
				respJson, err := httputil.HttpPostJson(data.Param, map[string]string{"ruleId": data.RuleId, "tenantId": data.TenantId}, nil)
				if err != nil {
					logger.Logger().Error("requestId=", requestId, ", autoScaling request fail,", err)
				} else {
					logger.Logger().Info("requestId=", requestId, ", autoScaling request success, resp=", respJson)
				}
			}
		}
	}

	return nil
}

func (s *AlertRecordAddService) checkAndBuild(requestId string, alerts []*form.AlertRecordAlertsBean) ([]commonModels.AlertRecord, map[int][]interface{}) {
	var list []commonModels.AlertRecord
	handlerMap := s.initHandlerMap()

	for _, alert := range alerts {
		a := alert
		a.RequestId = requestId
		if !s.predicate(a) {
			logger.Logger().Info("requestId=", a.RequestId, ", filter check false, alert=", jsonutil.ToString(a))
			continue
		}
		//解析告警规则信息
		ruleDesc := &commonDtos.RuleDesc{}
		if err := jsonutil.ToObjectWithError(a.Annotations.Description, ruleDesc); err != nil {
			logger.Logger().Error("requestId=", a.RequestId, ", 序列化告警数据失败, ", jsonutil.ToString(a))
			continue
		}
		//获取告警联系人信息
		var contactGroups []*commonDtos.ContactGroupInfo
		if ruleDesc.GroupList != nil && len(ruleDesc.GroupList) > 0 {
			contactGroups = s.AlertRecordDao.FindContactInfoByGroupIds(ruleDesc.GroupList)
		}
		summaryData := s.parseSummary(a.Annotations.Summary)

		record := s.buildAlertRecord(a, ruleDesc, contactGroups, summaryData["currentValue"])
		list = append(list, *record)

		//告警处理
		handlerList := s.AlarmHandlerSvc.GetAlarmHandlerListByRuleId(ruleDesc.RuleId)
		if handlerList == nil || len(handlerList) <= 0 {
			continue
		}
		for _, handler := range handlerList {
			var data interface{}
			switch handler.HandleType {
			case handler_type.Email, handler_type.Sms:
				data = s.buildNoticeData(a, record, ruleDesc, contactGroups, summaryData, handler.HandleType)
			case handler_type.Http:
				data = s.buildAutoScalingData(a, ruleDesc, handler.HandleParams)
			default:

			}
			if data != nil {
				handlerMap[handler.HandleType] = append(handlerMap[handler.HandleType], data)
			}
		}
	}
	return list, handlerMap
}

func (s *AlertRecordAddService) initHandlerMap() map[int][]interface{} {
	m := map[int][]interface{}{}
	m[handler_type.Email] = []interface{}{}
	m[handler_type.Sms] = []interface{}{}
	m[handler_type.Http] = []interface{}{}
	return m
}

func (s *AlertRecordAddService) getDurationTime(now, startTime time.Time, period int) string {
	//持续时间，单位秒
	sub := now.Sub(startTime)
	return util.GetDateDiff(int(sub.Round(time.Second).Milliseconds()) + period*1000)
}

func (s *AlertRecordAddService) buildAlertRecord(alert *form.AlertRecordAlertsBean, ruleDesc *commonDtos.RuleDesc, contactGroups []*commonDtos.ContactGroupInfo, cv string) *commonModels.AlertRecord {
	now := commonUtil.GetNow()
	startTime := commonUtil.TimeParseForZone(alert.StartsAt)
	val := strconv.FormatFloat(ruleDesc.TargetValue, 'f', 2, 64)
	sourceId := ruleDesc.ResourceId
	if strutil.IsBlank(sourceId) {
		sourceId = ruleDesc.ResourceGroupId
	}
	return &commonModels.AlertRecord{
		Id:             strconv.FormatInt(snowflake.GetWorker().NextId(), 10),
		RequestId:      alert.RequestId,
		Status:         alert.Status,
		TenantId:       ruleDesc.TenantId,
		RuleId:         ruleDesc.RuleId,
		RuleName:       ruleDesc.RuleName,
		RuleSourceId:   ruleDesc.Source,
		RuleSourceType: ruleDesc.SourceType,
		MonitorType:    ruleDesc.MonitorType,
		SourceType:     ruleDesc.Product,
		SourceId:       sourceId,
		Summary:        alert.Annotations.Summary,
		CurrentValue:   cv,
		StartTime:      commonUtil.TimeToStr(startTime, commonUtil.FullTimeFmt),
		EndTime:        commonUtil.TimeToStr(commonUtil.TimeParseForZone(alert.EndsAt), commonUtil.FullTimeFmt),
		TargetValue:    val,
		Expression:     ruleDesc.Express,
		Duration:       s.getDurationTime(now, startTime, ruleDesc.Period*ruleDesc.Time),
		Level:          ruleDesc.Level,
		AlarmKey:       ruleDesc.MetricName,
		ContactInfo:    jsonutil.ToString(contactGroups),
		Region:         config.Cfg.Common.RegionName,
		NoticeStatus:   "success",
		CreateTime:     now,
		UpdateTime:     now,
	}
}

func (s *AlertRecordAddService) buildAutoScalingData(alert *form.AlertRecordAlertsBean, ruleDesc *commonDtos.RuleDesc, param string) *commonDtos.AutoScalingData {
	if "resolved" == alert.Status {
		//告警恢复不需要通知弹性伸缩
		return nil
	}
	//弹性伸缩一定有资源组数据
	if strutil.IsBlank(ruleDesc.ResourceGroupId) {
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

func (s *AlertRecordAddService) buildNoticeData(alert *form.AlertRecordAlertsBean, record *commonModels.AlertRecord, ruleDesc *commonDtos.RuleDesc,
	contactGroups []*commonDtos.ContactGroupInfo, resourceInfo map[string]string, ht int) *service.AlertMsgSendDTO {
	source := message_center.ALERT_OPEN
	if "resolved" == alert.Status {
		source = message_center.ALERT_CANCEL
	}

	instanceInfo := s.getInstanceInfo(ruleDesc.ResourceId, resourceInfo)
	logger.Logger().Info("instanceId=", ruleDesc.ResourceId, ", get instanceInfo=", jsonutil.ToString(instanceInfo))

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

	val := strconv.FormatFloat(ruleDesc.TargetValue, 'f', 0, 64)

	if handler_type.Email == ht {
		objMap["instanceAmount"] = "1"
		objMap["period"] = util.GetDateDiff(ruleDesc.Period * 1000)
		objMap["times"] = "持续" + strconv.Itoa(ruleDesc.Time) + "个周期"
		objMap["time"] = util.GetDateDiff(ruleDesc.Period * 1000)
		objMap["calType"] = commonDao.GetAlarmStatisticsText(ruleDesc.Statistic)
		objMap["expression"] = commonDao.GetComparisonOperator(ruleDesc.ComparisonOperator)
		if source == message_center.ALERT_CANCEL {
			targetValue := ""
			if strutil.IsNotBlank(ruleDesc.Unit) {
				targetValue = ruleDesc.Unit
			}
			targetValue = val + targetValue
			objMap["targetValue"] = targetValue
		} else if source == message_center.ALERT_OPEN {
			objMap["targetValue"] = val
		}

		objMap["evaluationCount"] = "持续" + strconv.Itoa(ruleDesc.Time) + "个周期"
		// 暂不支持资源组
		objMap["InstanceID"] = ruleDesc.ResourceId
		if strutil.IsBlank(ruleDesc.Unit) {
			objMap["unit"] = ""
		} else {
			objMap["unit"] = ruleDesc.Unit
		}
	}

	var msgList []service.AlertMsgDTO
	if ht == handler_type.Sms {
		msg := s.buildNoticeMsg(contactGroups, message_center.Phone, objMap)
		if msg != nil {
			msgList = append(msgList, *msg)
		}
	} else if ht == handler_type.Email {
		msg := s.buildNoticeMsg(contactGroups, message_center.Email, objMap)
		if msg != nil {
			msgList = append(msgList, *msg)
		}
	}
	return &service.AlertMsgSendDTO{
		AlertId:    record.Id,
		SenderId:   ruleDesc.TenantId,
		SourceType: source,
		Msgs:       msgList,
	}
}

func (s *AlertRecordAddService) buildNoticeMsg(contactGroups []*commonDtos.ContactGroupInfo, rt message_center.ReceiveType, objMap map[string]string) *service.AlertMsgDTO {
	var ts []string
	for _, group := range contactGroups {
		for _, contact := range group.Contacts {
			var targets []string
			if message_center.Phone == rt {
				targets = s.buildPhoneTargets(contact)
			} else {
				targets = s.buildMailTargets(contact)
			}
			if targets != nil && len(targets) > 0 {
				ts = append(ts, targets...)
			}
		}
	}
	if len(ts) <= 0 {
		return nil
	}
	return &service.AlertMsgDTO{
		Type:    rt,
		Targets: ts,
		Content: jsonutil.ToString(objMap),
	}
}

func (s *AlertRecordAddService) getInstanceInfo(resourceId string, resourceInfo map[string]string) string {
	if strutil.IsBlank(resourceId) {
		return ""
	}
	instance := s.AlertRecordDao.FindFirstInstanceInfo(resourceId)
	if instance == nil {
		return ""
	}
	var ss []string
	if strutil.IsNotBlank(instance.InstanceName) {
		ss = append(ss, instance.InstanceName)
	}
	if strutil.IsNotBlank(instance.Ip) {
		ss = append(ss, instance.Ip)
	}
	for k, v := range resourceInfo {
		if k == "instance" {
			continue
		}
		if k == "name" {
			continue
		}
		if k == "currentValue" {
			continue
		}
		ss = append(ss, v)
	}
	return strings.Join(ss, "/")
}

func (s *AlertRecordAddService) buildMailTargets(contact commonDtos.UserContactInfo) []string {
	if strutil.IsBlank(contact.Mail) {
		return nil
	}
	var list []string
	for _, m := range strings.Split(contact.Mail, ",") {
		list = append(list, m)
	}
	return list
}

func (s *AlertRecordAddService) buildPhoneTargets(contact commonDtos.UserContactInfo) []string {
	if strutil.IsBlank(contact.Phone) {
		return nil
	}
	var list []string
	for _, m := range strings.Split(contact.Phone, ",") {
		list = append(list, m)
	}
	return list
}

func (s *AlertRecordAddService) parseSummary(summary string) map[string]string {
	var obj = make(map[string]string)
	if strutil.IsBlank(summary) {
		return obj
	}
	for _, s := range strings.Split(summary, ",") {
		labels := strings.Split(s, "=")
		if len(labels) == 2 {
			obj[labels[0]] = labels[1]
		}
	}
	return obj
}

func (s *AlertRecordAddService) predicate(alert *form.AlertRecordAlertsBean) bool {
	for _, f := range s.FilterChain {
		ret, err := f(alert)
		if err != nil {
			logger.Logger().Error("requestId=", alert.RequestId, ", alarm predicate fail, error=", err)
			return false
		}
		if !ret {
			return ret
		}
	}
	return true
}
