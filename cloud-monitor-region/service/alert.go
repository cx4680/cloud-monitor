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
			ruleDesc := &commonDtos.RuleDesc{}
			jsonutil.ToObject(alert.Annotations.Description, ruleDesc)
			if ruleDesc == nil {
				return false, errors.New("序列化告警数据失败")
			}
			// 判断该告警对应的规则是否有变化 资源组
			if strutil.IsNotBlank(ruleDesc.ResourceGroupId) {
				if num := commonDao.AlertRecord.FindAlertRuleBindGroupNum(ruleDesc.RuleId, ruleDesc.ResourceGroupId); num <= 0 {
					logger.Logger().Info("此告警规则已删除/禁用/解绑")
					return false, nil
				}
			} else if strutil.IsNotBlank(ruleDesc.ResourceId) {
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

func (s *AlertRecordAddService) Add(f form.InnerAlertRecordAddForm) error {
	if len(f.Alerts) == 0 {
		logger.Logger().Info("alerts 信息为空")
		return nil
	}
	list, handlerMap := s.checkAndBuild(f.Alerts)
	logger.Logger().Info("alarm data:", jsonutil.ToString(list), " || ", jsonutil.ToString(handlerMap))
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
				logger.Logger().Error("send alarm message fail,", err)
			}
		} else if t == handler_type.Http {
			//调用弹性伸缩
			for _, a := range p {
				data := a.(*commonDtos.AutoScalingData)
				respJson, err := httputil.HttpPostJson(data.Param+"/inner/as/trigger", map[string]string{"ruleId": data.RuleId, "tenantId": data.TenantId}, nil)
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

func (s *AlertRecordAddService) checkAndBuild(alerts []*form.AlertRecordAlertsBean) ([]commonModels.AlertRecord, map[int][]interface{}) {
	var list []commonModels.AlertRecord
	handlerMap := map[int][]interface{}{}

	for _, alert := range alerts {
		ret, err := s.predicate(alert)
		if err != nil {
			logger.Logger().Error("filter check error, alert=", jsonutil.ToString(alert), err)
			continue
		}
		if !ret {
			logger.Logger().Info("filter check false, alert=", jsonutil.ToString(alert))
			continue
		}
		//获取告警规则信息
		ruleDesc := &commonDtos.RuleDesc{}
		jsonutil.ToObject(alert.Annotations.Description, ruleDesc)
		if ruleDesc == nil {
			logger.Logger().Error("序列化告警数据失败, ", jsonutil.ToString(alert))
			continue
		}
		//获取告警联系人信息
		var contactGroups []*commonDtos.ContactGroupInfo
		if ruleDesc.GroupList != nil && len(ruleDesc.GroupList) > 0 {
			contactGroups = s.AlertRecordDao.FindContactInfoByGroupIds(ruleDesc.GroupList)
		}

		alertAnnotationStr := jsonutil.ToString(alert.Annotations)
		if strutil.IsBlank(alertAnnotationStr) {
			logger.Logger().Info("告警数据为空, ", jsonutil.ToString(alert))
			continue
		}
		labelMap := s.getLabelMap(alert.Annotations.Summary)

		record := s.buildAlertRecord(alert, ruleDesc, contactGroups, labelMap["currentValue"])
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
			case handler_type.Email, handler_type.Sms:
				data = s.buildNoticeData(alert, record, ruleDesc, contactGroups, handler.HandleType)
			case handler_type.Http:
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

func (s *AlertRecordAddService) getDurationTime(now, startTime time.Time, period int) string {
	//持续时间，单位秒
	sub := now.Sub(startTime)
	return util.GetDateDiff(int(sub.Round(time.Second).Milliseconds()) + period*1000)
}

//截取联系人字符串，防止联系人太多导致数据保存失败
func sub(list []*commonDtos.ContactGroupInfo, num int) string {
	str := jsonutil.ToString(list)
	if num > 5 {
		logger.Logger().Info("contact str too long, set it empty")
		return ""
	}
	if len(str) > 1000 {
		return sub(list[0:len(list)/2], num+1)
	}
	return str
}

func (s *AlertRecordAddService) buildAlertRecord(alert *form.AlertRecordAlertsBean, ruleDesc *commonDtos.RuleDesc, contactGroups []*commonDtos.ContactGroupInfo, cv string) *commonModels.AlertRecord {
	now := commonUtil.GetNow()
	startTime := commonUtil.TimeParseForZone(alert.StartsAt)

	val := strconv.FormatFloat(ruleDesc.TargetValue, 'f', 2, 64)

	sourceId := ruleDesc.ResourceId
	if strutil.IsBlank(sourceId) {
		sourceId = ruleDesc.ResourceGroupId
	}
	contactStr := sub(contactGroups, 1)

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
		CurrentValue: cv,
		StartTime:    commonUtil.TimeToStr(startTime, commonUtil.FullTimeFmt),
		EndTime:      commonUtil.TimeToStr(commonUtil.TimeParseForZone(alert.EndsAt), commonUtil.FullTimeFmt),
		TargetValue:  val,
		Expression:   ruleDesc.Express,
		Duration:     s.getDurationTime(now, startTime, ruleDesc.Period),
		Level:        ruleDesc.Level,
		AlarmKey:     ruleDesc.MetricName,
		Region:       config.Cfg.Common.RegionName,
		NoticeStatus: "success",
		ContactInfo:  contactStr,
		CreateTime:   global.JsonTime{Time: now},
		UpdateTime:   global.JsonTime{Time: now},
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
	contactGroups []*commonDtos.ContactGroupInfo, ht int) *service.AlertMsgSendDTO {
	source := message_center.ALERT_OPEN
	if "resolved" == alert.Status {
		source = message_center.ALERT_CANCEL
	}

	instanceInfo := s.getInstanceInfo(ruleDesc.ResourceId)
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

	val := strconv.FormatFloat(ruleDesc.TargetValue, 'f', 2, 64)

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

	//TODO
	if ht == handler_type.Sms {
		for _, group := range contactGroups {
			for _, contact := range group.Contacts {
				msgList = append(msgList, service.AlertMsgDTO{
					Type:    message_center.Phone,
					Targets: s.buildPhoneTargets(contact),
					Content: jsonutil.ToString(objMap),
				})
			}
		}
	} else if ht == handler_type.Email {
		for _, group := range contactGroups {
			for _, contact := range group.Contacts {
				msgList = append(msgList, service.AlertMsgDTO{
					Type:    message_center.Email,
					Targets: s.buildMailTargets(contact),
					Content: jsonutil.ToString(objMap),
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

func (s *AlertRecordAddService) getInstanceInfo(resourceId string) string {
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

func (s *AlertRecordAddService) predicate(alert *form.AlertRecordAlertsBean) (bool, error) {
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
