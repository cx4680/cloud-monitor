package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/constant"
	commonDao "code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dto"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/errors"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sys_component/sys_rocketmq"
	commonModels "code.cestc.cn/ccos-ops/cloud-monitor/business-common/model"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service/external/message_center"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/httputil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/snowflake"
	"log"
	"strconv"
)

type MessageService struct {
	NotificationRecordDao *commonDao.NotificationRecordDao
	MCS                   *message_center.Service
}

type AlertMsgSendDTO struct {
	AlertId    string
	SenderId   string
	SourceType message_center.MsgSource
	Msgs       []AlertMsgDTO
}

type AlertMsgDTO struct {
	Type    message_center.ReceiveType
	Targets []string
	Content string
}

func NewMessageService(messageCenterService *message_center.Service) *MessageService {
	return &MessageService{
		NotificationRecordDao: commonDao.NotificationRecord,
		MCS:                   messageCenterService,
	}
}

func (s *MessageService) TargetFilter(targetList []string, rt message_center.ReceiveType, senderId string, isCenter bool) []string {
	if targetList == nil || len(targetList) <= 0 {
		return nil
	}
	var addrs []string
	for _, t := range targetList {
		if message_center.Email == rt {
			addrs = append(addrs, t)
			continue
		}
		//短信需要验证是否超限
		num := s.NotificationRecordDao.GetTenantPhoneCurrentMonthRecordNum(senderId)
		if s.checkSentNum(senderId, num, isCenter) {
			addrs = append(addrs, t)
		} else {
			logger.Logger().Infof("too many records have been sent, send refused, sender=%s \n", senderId)
		}
	}
	//去重
	addrs = util.RemoveDuplicateElement(addrs)
	return addrs
}

func (s *MessageService) SendAlarmNotice(msgList []interface{}) error {
	if config.MsgOpen != config.Cfg.Common.MsgIsOpen {
		logger.Logger().Info("There is no message center for this env")
		return nil
	}
	if msgList == nil || len(msgList) <= 0 {
		return nil
	}

	var recordList []commonModels.NotificationRecord
	var sendMsgList []message_center.MessageSendDTO
	//统计短信数量
	var smsSender []string
	for _, alertMsg := range msgList {
		am := alertMsg.(*AlertMsgSendDTO)
		for _, msg := range am.Msgs {
			newTargets := s.TargetFilter(msg.Targets, msg.Type, am.SenderId, false)
			sendMsgList = append(sendMsgList, message_center.MessageSendDTO{
				SenderId:   am.SenderId,
				Type:       msg.Type,
				SourceType: am.SourceType,
				Targets:    newTargets,
				Content:    msg.Content,
			})
			if msg.Type == message_center.Phone {
				smsSender = append(smsSender, msg.Targets...)
			}
			for _, addr := range newTargets {
				recordList = append(recordList, commonModels.NotificationRecord{
					Id:               strconv.FormatInt(snowflake.GetWorker().NextId(), 10),
					SenderId:         am.SenderId,
					SourceId:         am.AlertId,
					SourceType:       uint8(am.SourceType),
					TargetAddress:    addr,
					NotificationType: uint8(msg.Type),
					Result:           1,
					CreateTime:       util.GetNow(),
				})
			}
		}
	}
	//send msg
	if err := s.MCS.SendBatch(sendMsgList); err != nil {
		logger.Logger().Errorf("message send error, %v\n\n", err)
		return err
	}
	//save record local
	s.NotificationRecordDao.InsertBatch(global.DB, recordList)

	//	sync record to center
	_ = sys_rocketmq.SendRocketMqMsg(sys_rocketmq.RocketMqMsg{
		Topic:   sys_rocketmq.NotificationSyncTopic,
		Content: jsonutil.ToString(recordList),
	})
	// 发送短信余量不足提醒
	if len(smsSender) > 0 {
		smsSender = util.RemoveDuplicateElement(smsSender)
		//通过MQ异步解耦
		_ = sys_rocketmq.SendRocketMqMsg(sys_rocketmq.RocketMqMsg{
			Topic:   sys_rocketmq.SmsMarginReminderTopic,
			Content: jsonutil.ToString(smsSender),
		})
	}
	return nil
}

// SendCertifyMsg 发送激活信息
func (s *MessageService) SendCertifyMsg(msg message_center.MessageSendDTO, contactId string) {
	//send msg
	var ret uint8 = 1
	if err := s.MCS.Send(msg); err != nil {
		logger.Logger().Errorf("message send error, %v\n\n", err)
		ret = 0
	}
	record := commonModels.NotificationRecord{
		Id:               strconv.FormatInt(snowflake.GetWorker().NextId(), 10),
		SenderId:         msg.SenderId,
		SourceId:         contactId,
		SourceType:       uint8(message_center.VERIFY),
		TargetAddress:    msg.Targets[0],
		NotificationType: uint8(msg.Type),
		Result:           ret,
		CreateTime:       util.GetNow(),
	}
	s.NotificationRecordDao.Insert(global.DB, record)
}

// SmsMarginReminder 短信余量提醒
func (s *MessageService) SmsMarginReminder(sender string) {
	count := s.NotificationRecordDao.GetTenantSMSLackRecordNum(sender)
	if count > 0 {
		//已发送过短信提醒
		return
	}
	alreadySendNum := s.getUserCurrentMonthSmsUsedNum(sender)
	if alreadySendNum < constant.ThresholdSmsNum {
		//未达到提醒阈值
		return
	}
	tenantDTO := dto.TenantDTO{}
	loginName := tenantDTO.Name
	serialNumber := tenantDTO.Phone

	params := make(map[string]string)
	params["userName"] = loginName
	params["msgUsed"] = strconv.Itoa(alreadySendNum)
	params["msgLeft"] = strconv.Itoa(constant.MaxSmsNum - alreadySendNum)
	params["msgInitial"] = strconv.Itoa(constant.MaxSmsNum)

	remainderMsg := message_center.MessageSendDTO{
		SenderId:   sender,
		Type:       message_center.Phone,
		Targets:    []string{serialNumber},
		SourceType: message_center.SMS_LACK,
		Content:    jsonutil.ToString(params),
	}
	if err := s.MCS.Send(remainderMsg); err != nil {
		return
	}
	//保存发送记录
	s.NotificationRecordDao.Insert(global.DB, commonModels.NotificationRecord{
		SenderId:         sender,
		SourceId:         "sms-lack-" + sender,
		SourceType:       uint8(message_center.SMS_LACK),
		TargetAddress:    serialNumber,
		NotificationType: uint8(message_center.Phone),
		Result:           1,
		CreateTime:       util.GetNow(),
	})
}

func (s *MessageService) checkSentNum(tenantId string, num int, isCenter bool) bool {
	//check local
	if num > constant.MaxSmsNum {
		logger.Logger().Info("user ", tenantId, " already used more", constant.MaxSmsNum, ", send sms refused.")
		return false
	}
	//	check remote
	if !isCenter {
		num = s.getUserCurrentMonthSmsUsedNum(tenantId)
		if num > constant.MaxSmsNum {
			logger.Logger().Info("user ", tenantId, " already used more", constant.MaxSmsNum, ", send sms refused.")
			return false
		}
	}
	return true
}

func (s *MessageService) sendNotification(sender string, num int) []commonModels.NotificationRecord {
	count := s.NotificationRecordDao.GetTenantSMSLackRecordNum(sender)
	if count > 0 {
		return nil
	}
	if num < constant.ThresholdSmsNum {
		return nil
	}
	tenantDTO := dto.TenantDTO{}
	logingName := tenantDTO.Name
	serialNumber := tenantDTO.Phone

	params := make(map[string]string)
	params["userName"] = logingName
	params["msgUsed"] = strconv.Itoa(num)
	params["msgLeft"] = strconv.Itoa(constant.MaxSmsNum - num)
	params["msgInitial"] = strconv.Itoa(constant.MaxSmsNum)

	var noticeMsgDTOList []*dto.NoticeMsgDTO
	noticeMsgDTO := dto.NoticeMsgDTO{
		SourceId: "sms-lack-" + sender,
		TenantId: sender,
		MsgEvent: dto.MsgEvent{
			Type:   1, //TODO 枚举
			Source: dto.SMS_LACK,
		},
		RevObjectBean: dto.RecvObjectBean{
			RecvObjectType: 1, //TODO 枚举
			RecvObject:     serialNumber,
			NoticeContent:  jsonutil.ToString(params),
		},
	}
	noticeMsgDTOList = append(noticeMsgDTOList, &noticeMsgDTO)
	//TODO send to message center
	//s.sendToMsgCenter(noticeMsgDTOList)

	return s.saveNotificationRecords(noticeMsgDTOList)

}

func (s *MessageService) saveNotificationRecords(noticeMsgDTOS []*dto.NoticeMsgDTO) []commonModels.NotificationRecord {
	var recordList []commonModels.NotificationRecord
	for _, noticeMsgDTO := range noticeMsgDTOS {
		recordList = append(recordList, commonModels.NotificationRecord{
			Id:               strconv.FormatInt(snowflake.GetWorker().NextId(), 10),
			SenderId:         noticeMsgDTO.TenantId,
			SourceId:         noticeMsgDTO.SourceId,
			SourceType:       uint8(noticeMsgDTO.MsgEvent.Source),
			TargetAddress:    noticeMsgDTO.RevObjectBean.RecvObject,
			NotificationType: uint8(noticeMsgDTO.MsgEvent.Type),
			Result:           1,
			CreateTime:       util.GetNow(),
		})
	}
	s.NotificationRecordDao.InsertBatch(global.DB, recordList)
	return recordList
}

type ResultDTO struct {
	ErrorMsg   string
	ErrorCode  string
	Success    bool
	Module     int
	AllowRetry bool
	ErrorArgs  []interface{}
}

func (s *MessageService) getUserCurrentMonthSmsUsedNum(tenantId string) int {
	resp, err := httputil.HttpGet(config.Cfg.Common.HawkeyeCenterPath + "/inner/getUsage?tenantId=" + tenantId)
	if err != nil {
		log.Fatal("获取用户短信月使用量出错, tenantId=" + tenantId)
		return 0
	}
	var respObj ResultDTO
	jsonutil.ToObject(resp, &respObj)
	return respObj.Module
}

func (s *MessageService) GetTenantCurrentMonthSmsUsedNum(tenantId string) (int, error) {
	if tenantId == "" {
		return 0, errors.NewBusinessError("租户不能为空")
	}
	num := s.NotificationRecordDao.GetTenantSMSLackRecordNum(tenantId)
	return num, nil
}
