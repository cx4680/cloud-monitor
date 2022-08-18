package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/snowflake"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/dto"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/enum/handler_type"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/model"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/service/external/message_center"
	util2 "code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/constant"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/k8s"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/util"
	"context"
	"errors"
	"strconv"
	"strings"
	"time"
)

type filter func(*form.AlarmRecordAlertsBean) (bool, error)

type AlarmRecordAddService struct {
	FilterChain     []filter
	AlarmRecordSvc  *AlarmRecordService
	AlarmHandlerSvc *service.AlarmHandlerService
	TenantSvc       *service.TenantService
	AlarmRecordDao  *dao.AlarmRecordDao
}

func NewAlarmRecordAddService(AlarmRecordSvc *AlarmRecordService, AlarmHandlerSvc *service.AlarmHandlerService, TenantSvc *service.TenantService) *AlarmRecordAddService {
	return &AlarmRecordAddService{
		FilterChain: []filter{func(alert *form.AlarmRecordAlertsBean) (bool, error) {
			ruleDesc := &k8s.AlarmDescription{}
			if err := jsonutil.ToObjectWithError(alert.Annotations.Description, ruleDesc); err != nil {
				return false, errors.New("requestId=" + alert.RequestId + ", 序列化告警数据失败")
			}
			// 判断该告警对应的规则是否有变化 资源组
			if strutil.IsNotBlank(ruleDesc.ResourceGroupId) {
				if num := dao.AlarmRecord.FindAlertRuleBindGroupNum(ruleDesc.Rule.BizId, ruleDesc.ResourceGroupId); num <= 0 {
					logger.Logger().Info("requestId=" + alert.RequestId + ",此告警规则已删除/禁用/解绑")
					return false, nil
				}
			} else if strutil.IsNotBlank(ruleDesc.ResourceId) {
				if num := dao.AlarmRecord.FindAlertRuleBindResourceNum(ruleDesc.Rule.BizId, ruleDesc.ResourceId); num <= 0 {
					logger.Logger().Info("requestId=" + alert.RequestId + ",此告警规则已删除/禁用/解绑")
					return false, nil
				}
			}
			items := dao.AlarmItem.GetItemListByRuleBizId(global.DB, ruleDesc.Rule.BizId)
			conditionId, err := util.MD5(items)
			if err != nil {
				logger.Logger().Info("requestId=", alert.RequestId, "告警触发条件MD5失败，此告警触发条件改变")
				return false, nil
			}
			if conditionId != strings.Split(alert.Labels.Alertname, "#")[2] {
				logger.Logger().Info("requestId=", alert.RequestId, "此告警触发条件改变, Alertname=", alert.Labels.Alertname, ", conditionId=", conditionId)
				return false, nil
			}
			return true, nil
		}},
		AlarmHandlerSvc: AlarmHandlerSvc,
		AlarmRecordSvc:  AlarmRecordSvc,
		TenantSvc:       TenantSvc,
		AlarmRecordDao:  dao.AlarmRecord,
	}
}

func (s *AlarmRecordAddService) Add(ctx *context.Context, f form.InnerAlarmRecordAddForm) error {
	if len(f.Alerts) == 0 {
		logger.Logger().Info("requestId=", util2.GetRequestId(ctx), ", alerts 信息为空")
		return errors.New("告警信息为空")
	}
	list, infoList, events := s.checkAndBuild(util2.GetRequestId(ctx), f.Alerts)
	logger.Logger().Info("requestId=", util2.GetRequestId(ctx), ", alarm data=", jsonutil.ToString(list), ", handler data=", jsonutil.ToString(events))
	//持久化
	if len(list) <= 0 || len(infoList) <= 0 {
		logger.Logger().Info("requestId=", util2.GetRequestId(ctx), " alarm info list is empty! ")
		return nil
	}

	if err := s.AlarmRecordSvc.InsertAndHandler(ctx, list, infoList, events); err != nil {
		return err
	}
	return nil
}

func (s *AlarmRecordAddService) checkAndBuild(requestId string, alerts []*form.AlarmRecordAlertsBean) ([]model.AlarmRecord, []model.AlarmInfo, []interface{}) {
	var list []model.AlarmRecord
	var infoList []model.AlarmInfo
	var alarmHandlerEventList []interface{}

	for _, alert := range alerts {
		a := alert
		a.RequestId = requestId
		if !s.predicate(a) {
			logger.Logger().Info("requestId=", a.RequestId, ", filter check false, alert=", jsonutil.ToString(a))
			continue
		}
		//解析告警规则信息
		ruleDesc := &k8s.AlarmDescription{}
		if err := jsonutil.ToObjectWithError(a.Annotations.Description, ruleDesc); err != nil {
			logger.Logger().Error("requestId=", a.RequestId, ", 序列化告警数据失败, ", jsonutil.ToString(a))
			continue
		}
		//获取告警联系人信息
		var contactGroups []*dto.ContactGroupInfo
		if ruleDesc.ContactGroupIds != nil && len(ruleDesc.ContactGroupIds) > 0 {
			contactGroups = s.AlarmRecordDao.FindContactInfoByGroupIds(ruleDesc.ContactGroupIds)
		}
		summaryData := s.parseSummary(a.Annotations.Summary)

		record := s.buildAlertRecord(a, ruleDesc, summaryData["currentValue"])
		list = append(list, *record)

		alarmInfo := s.buildAlarmInfo(a, record.BizId, ruleDesc, contactGroups)
		infoList = append(infoList, *alarmInfo)

		//告警处理
		handlerList := s.AlarmHandlerSvc.GetAlarmHandlerListByRuleId(ruleDesc.Rule.BizId)
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
				alarmHandlerEventList = append(alarmHandlerEventList, AlarmHandlerEvent{
					RequestId: requestId,
					Type:      handler.HandleType,
					Data:      data,
				})
			}
		}
	}
	return list, infoList, alarmHandlerEventList
}

func (s *AlarmRecordAddService) getDurationTime(now, startTime time.Time, period int) string {
	//持续时间，单位秒
	sub := now.Sub(startTime)
	return util.GetDateDiff(int(sub.Round(time.Second).Milliseconds()) + period*1000)
}

func (s *AlarmRecordAddService) buildAlarmInfo(alert *form.AlarmRecordAlertsBean, alarmBizId string, ruleDesc *k8s.AlarmDescription, contactGroups []*dto.ContactGroupInfo) *model.AlarmInfo {
	return &model.AlarmInfo{
		AlarmBizId:  alarmBizId,
		Summary:     alert.Annotations.Summary,
		Expression:  ruleDesc.ExprDetail,
		ContactInfo: jsonutil.ToString(contactGroups),
	}
}

func (s *AlarmRecordAddService) buildAlertRecord(alert *form.AlarmRecordAlertsBean, ruleDesc *k8s.AlarmDescription, cv string) *model.AlarmRecord {
	now := util2.GetNow()
	startTime := util2.TimeParseForZone(alert.StartsAt)
	ruleType := ruleDesc.Rule.Type

	period := ruleDesc.Rule.Period
	times := ruleDesc.Rule.Times
	level := ruleDesc.Rule.Level
	if constant.AlarmRuleTypeSingleMetric == ruleType {
		item := ruleDesc.RuleItems[0]
		period = item.TriggerCondition.Period
		times = item.TriggerCondition.Times
		level = item.Level
	}
	sourceId := ruleDesc.ResourceId
	if strutil.IsBlank(sourceId) {
		sourceId = ruleDesc.ResourceGroupId
	}
	return &model.AlarmRecord{
		BizId:          strconv.FormatInt(snowflake.GetWorker().NextId(), 10),
		RequestId:      alert.RequestId,
		Status:         alert.Status,
		TenantId:       ruleDesc.Rule.TenantID,
		RuleId:         ruleDesc.Rule.BizId,
		RuleName:       ruleDesc.Rule.Name,
		RuleSourceId:   ruleDesc.Rule.Source,
		RuleSourceType: ruleDesc.Rule.SourceType,
		MonitorType:    ruleDesc.Rule.MonitorType,
		SourceType:     ruleDesc.Rule.ProductName,
		SourceId:       sourceId,
		CurrentValue:   cv,
		StartTime:      startTime,
		EndTime:        util2.TimeParseForZone(alert.EndsAt),
		Duration:       s.getDurationTime(now, startTime, period*times),
		Level:          level,
		Region:         config.Cfg.Common.RegionName,
		CreateTime:     now,
		UpdateTime:     now,
	}
}

func (s *AlarmRecordAddService) buildAutoScalingData(alert *form.AlarmRecordAlertsBean, ruleDesc *k8s.AlarmDescription, param string) *dto.AutoScalingData {
	if "resolved" == alert.Status {
		//告警恢复不需要通知弹性伸缩
		return nil
	}
	//弹性伸缩一定有资源组数据
	if strutil.IsBlank(ruleDesc.ResourceGroupId) {
		//脏数据
		return nil
	}
	sm := s.parseSummary(alert.Annotations.Summary)
	it := alert.Labels.InstanceType

	return &dto.AutoScalingData{
		TenantId:        ruleDesc.Rule.TenantID,
		RuleId:          ruleDesc.Rule.BizId,
		ResourceGroupId: ruleDesc.ResourceGroupId,
		ResourceId:      sm["instance"],
		ResourceType:    it,
		Param:           param,
	}
}

//func (s *AlarmRecordAddService) genExprDetails(rule model.AlarmRule, items []model.AlarmItem) string {
//	ss := make([]string, len(items))
//	for i, item := range items {
//		ss[i] = dao.GetExpress2(*item.TriggerCondition)
//	}
//	if constant.AlarmRuleTypeSingleMetric == rule.Type {
//		return ss[0]
//	}
//
//	if constant.AlarmRuleCombinationAnd == rule.Combination {
//		return strings.Join(ss, " 并且 ")
//	}
//	return strings.Join(ss, " 或者 ")
//
//}

func (s *AlarmRecordAddService) buildNoticeData(alert *form.AlarmRecordAlertsBean, record *model.AlarmRecord, ruleDesc *k8s.AlarmDescription,
	contactGroups []*dto.ContactGroupInfo, resourceInfo map[string]string, ht int) *service.AlertMsgSendDTO {
	source := message_center.ALERT_OPEN
	if "resolved" == alert.Status {
		source = message_center.ALERT_CANCEL
	}

	instanceInfo := s.getInstanceInfo(ruleDesc.ResourceId, resourceInfo)
	logger.Logger().Info("instanceId=", ruleDesc.ResourceId, ", get instanceInfo=", jsonutil.ToString(instanceInfo))

	//period := ruleDesc.Rule.Period
	//times := ruleDesc.Rule.Times
	//if constant.AlarmRuleTypeSingleMetric == ruleDesc.Rule.Type {
	//item := ruleDesc.RuleItems[0]
	//period = item.TriggerCondition.Period
	//times = item.TriggerCondition.Times
	//}

	objMap := make(map[string]string)
	objMap["duration"] = record.Duration
	objMap["instanceInfo"] = instanceInfo
	objMap["product"] = ruleDesc.Rule.ProductName
	objMap["regionName"] = record.Region
	//objMap["metricName"] = ruleDesc.MonitorItem
	objMap["Name"] = ruleDesc.Rule.Name
	objMap["userName"] = s.TenantSvc.GetTenantInfo(ruleDesc.Rule.TenantID).Name
	objMap["alertTime"] = record.StartTime.Format(util2.FullTimeFmt)
	objMap["recoveryTime"] = record.EndTime.Format(util2.FullTimeFmt)
	objMap["exprDetails"] = ruleDesc.ExprDetail
	//f, _ := strconv.ParseFloat(record.CurrentValue, 64)
	//cv := fmt.Sprintf("%.2f", f)
	//objMap["currentValue"] = cv + ruleDesc.Unit

	//val := strconv.FormatFloat(ruleDesc.TargetValue, 'f', 0, 64)

	if handler_type.Email == ht {
		objMap["instanceAmount"] = "1"
		//objMap["period"] = util.GetDateDiff(period * 1000)
		//objMap["times"] = "持续" + strconv.Itoa(times) + "个周期"
		//objMap["time"] = util.GetDateDiff(period * 1000)
		//objMap["calType"] = dao.GetAlarmStatisticsText(ruleDesc.Statistic)
		//objMap["expression"] = dao.GetComparisonOperator(ruleDesc.ComparisonOperator)
		//if source == message_center.ALERT_CANCEL {
		//targetValue := ""
		//if strutil.IsNotBlank(ruleDesc.Unit) {
		//	targetValue = ruleDesc.Unit
		//}
		//targetValue = val + targetValue
		//objMap["targetValue"] = targetValue
		//} else if source == message_center.ALERT_OPEN {
		//objMap["targetValue"] = val
		//}

		//objMap["evaluationCount"] = "持续" + strconv.Itoa(times) + "个周期"
		// 暂不支持资源组
		objMap["InstanceID"] = ruleDesc.ResourceId
		//if strutil.IsBlank(ruleDesc.Unit) {
		//	objMap["unit"] = ""
		//} else {
		//	objMap["unit"] = ruleDesc.Unit
		//}
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
		AlertId:    record.BizId,
		SenderId:   ruleDesc.Rule.TenantID,
		SourceType: source,
		Msgs:       msgList,
	}
}

func (s *AlarmRecordAddService) buildNoticeMsg(contactGroups []*dto.ContactGroupInfo, rt message_center.ReceiveType, objMap map[string]string) *service.AlertMsgDTO {
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

func (s *AlarmRecordAddService) getInstanceInfo(resourceId string, resourceInfo map[string]string) string {
	if strutil.IsBlank(resourceId) {
		return ""
	}
	instance := s.AlarmRecordDao.FindFirstInstanceInfo(resourceId)
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

func (s *AlarmRecordAddService) buildMailTargets(contact dto.UserContactInfo) []string {
	if strutil.IsBlank(contact.Mail) {
		return nil
	}
	var list []string
	for _, m := range strings.Split(contact.Mail, ",") {
		list = append(list, m)
	}
	return list
}

func (s *AlarmRecordAddService) buildPhoneTargets(contact dto.UserContactInfo) []string {
	if strutil.IsBlank(contact.Phone) {
		return nil
	}
	var list []string
	for _, m := range strings.Split(contact.Phone, ",") {
		list = append(list, m)
	}
	return list
}

func (s *AlarmRecordAddService) parseSummary(summary string) map[string]string {
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

func (s *AlarmRecordAddService) predicate(alert *form.AlarmRecordAlertsBean) bool {
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
