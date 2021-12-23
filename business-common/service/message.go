package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/constants"
	commonDao "code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dtos"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/errors"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysComponent/sysRocketMq"
	commonModels "code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service/external/messageCenter"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/utils/snowflake"
	"log"
	"strconv"
	"time"
)

type MessageService struct {
	NotificationRecordDao *commonDao.NotificationRecordDao
	MCS                   *messageCenter.Service
}

type AlertMsgSendDTO struct {
	AlertId    string
	SenderId   string
	SourceType messageCenter.MsgSource
	Msgs       []AlertMsgDTO
}

type AlertMsgDTO struct {
	Type    messageCenter.ReceiveType
	Targets []string
	Content string
}

func NewMessageService(messageCenterService *messageCenter.Service) *MessageService {
	return &MessageService{
		NotificationRecordDao: commonDao.NotificationRecord,
		MCS:                   messageCenterService,
	}
}

func (s *MessageService) TargetFilter(targetList []string, rt messageCenter.ReceiveType, senderId string, isCenter bool) []string {
	if targetList == nil || len(targetList) <= 0 {
		return nil
	}
	var addrs []string
	for _, t := range targetList {
		if messageCenter.Email == rt {
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
	return addrs
}

func (s *MessageService) SendAlarmNotice(msgList []interface{}) error {
	if !config.GetCommonConfig().HasNoticeModel {
		logger.Logger().Info("There is no message center for this env")
		return nil
	}
	if msgList == nil || len(msgList) <= 0 {
		return nil
	}

	var recordList []commonModels.NotificationRecord
	var sendMsgList []messageCenter.MessageSendDTO
	//统计短信数量
	var smsSender []string
	for _, alertMsg := range msgList {
		am := alertMsg.(*AlertMsgSendDTO)
		for _, msg := range am.Msgs {
			newTargets := s.TargetFilter(msg.Targets, msg.Type, am.SenderId, false)
			sendMsgList = append(sendMsgList, messageCenter.MessageSendDTO{
				SenderId:   am.SenderId,
				Type:       msg.Type,
				SourceType: am.SourceType,
				Targets:    newTargets,
				Content:    msg.Content,
			})
			if msg.Type == messageCenter.Phone {
				smsSender = append(smsSender, msg.Targets...)
			}
			for _, addr := range newTargets {
				recordList = append(recordList, commonModels.NotificationRecord{
					SenderId:         am.SenderId,
					SourceId:         am.AlertId,
					SourceType:       int(am.SourceType),
					TargetAddress:    addr,
					NotificationType: int(msg.Type),
					Result:           1,
				})
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
		_ = sysRocketMq.SendRocketMqMsg(sysRocketMq.RocketMqMsg{
			Topic:   sysRocketMq.NotificationSyncTopic,
			Content: tools.ToString(recordList),
		})
		// 发送短信余量不足提醒
		if len(smsSender) > 0 {
			smsSender = tools.RemoveDuplicateElement(smsSender)
			//通过MQ异步解耦
			_ = sysRocketMq.SendRocketMqMsg(sysRocketMq.RocketMqMsg{
				Topic:   sysRocketMq.SmsMarginReminderTopic,
				Content: tools.ToString(smsSender),
			})
		}
	}

	return nil
}

// SmsMarginReminder 短信余量提醒
func (s *MessageService) SmsMarginReminder(sender string) {
	count := s.NotificationRecordDao.GetTenantSMSLackRecordNum(sender)
	if count > 0 {
		//已发送过短信提醒
		return
	}
	alreadySendNum := s.getUserCurrentMonthSmsUsedNum(sender)
	if alreadySendNum < constants.ThresholdSmsNum {
		//未达到提醒阈值
		return
	}
	tenantDTO := dtos.TenantDTO{}
	loginName := tenantDTO.Name
	serialNumber := tenantDTO.Phone

	params := make(map[string]string)
	params["userName"] = loginName
	params["msgUsed"] = strconv.Itoa(alreadySendNum)
	params["msgLeft"] = strconv.Itoa(constants.MaxSmsNum - alreadySendNum)
	params["msgInitial"] = strconv.Itoa(constants.MaxSmsNum)

	remainderMsg := messageCenter.MessageSendDTO{
		SenderId:   sender,
		Type:       messageCenter.Phone,
		Targets:    []string{serialNumber},
		SourceType: messageCenter.SMS_LACK,
		Content:    tools.ToString(params),
	}
	if err := s.MCS.Send(remainderMsg); err != nil {
		return
	}
	//保存发送记录
	s.NotificationRecordDao.Insert(global.DB, commonModels.NotificationRecord{
		SenderId:         sender,
		SourceId:         "sms-lack-" + sender,
		SourceType:       int(messageCenter.SMS_LACK),
		TargetAddress:    serialNumber,
		NotificationType: int(messageCenter.Phone),
		Result:           1,
	})
}

func (s *MessageService) checkSentNum(tenantId string, num int, isCenter bool) bool {
	//check local
	if num > constants.MaxSmsNum {
		logger.Logger().Info("user ", tenantId, " already used more", constants.MaxSmsNum, ", send sms refused.")
		return false
	}
	//	check remote
	if !isCenter {
		num = s.getUserCurrentMonthSmsUsedNum(tenantId)
		if num > constants.MaxSmsNum {
			logger.Logger().Info("user ", tenantId, " already used more", constants.MaxSmsNum, ", send sms refused.")
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
	if num < constants.ThresholdSmsNum {
		return nil
	}
	tenantDTO := dtos.TenantDTO{}
	logingName := tenantDTO.Name
	serialNumber := tenantDTO.Phone

	params := make(map[string]string)
	params["userName"] = logingName
	params["msgUsed"] = strconv.Itoa(num)
	params["msgLeft"] = strconv.Itoa(constants.MaxSmsNum - num)
	params["msgInitial"] = strconv.Itoa(constants.MaxSmsNum)

	var noticeMsgDTOList []*dtos.NoticeMsgDTO
	noticeMsgDTO := dtos.NoticeMsgDTO{
		SourceId: "sms-lack-" + sender,
		TenantId: sender,
		MsgEvent: dtos.MsgEvent{
			Type:   1, //TODO 枚举
			Source: dtos.SMS_LACK,
		},
		RevObjectBean: dtos.RecvObjectBean{
			RecvObjectType: 1, //TODO 枚举
			RecvObject:     serialNumber,
			NoticeContent:  tools.ToString(params),
		},
	}
	noticeMsgDTOList = append(noticeMsgDTOList, &noticeMsgDTO)
	//TODO send to message center
	//s.sendToMsgCenter(noticeMsgDTOList)

	return s.saveNotificationRecords(noticeMsgDTOList)

}

func (s *MessageService) saveNotificationRecords(noticeMsgDTOS []*dtos.NoticeMsgDTO) []commonModels.NotificationRecord {
	var recordList []commonModels.NotificationRecord
	now := time.Now().Format("2006-01-02 15:04:05")
	for _, noticeMsgDTO := range noticeMsgDTOS {
		recordList = append(recordList, commonModels.NotificationRecord{
			Id:               strconv.FormatInt(snowflake.GetWorker().NextId(), 10),
			SenderId:         noticeMsgDTO.TenantId,
			SourceId:         noticeMsgDTO.SourceId,
			SourceType:       int(noticeMsgDTO.MsgEvent.Source),
			TargetAddress:    noticeMsgDTO.RevObjectBean.RecvObject,
			NotificationType: int(noticeMsgDTO.MsgEvent.Type),
			Result:           1,
			CreateTime:       now,
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
	resp, err := tools.HttpGet(config.GetCommonConfig().HawkeyeCenterPath + "/inner/getUsage?tenantId=" + tenantId)
	if err != nil {
		log.Fatal("获取用户短信月使用量出错, tenantId=" + tenantId)
		return 0
	}
	var respObj ResultDTO
	tools.ToObject(resp, &respObj)
	return respObj.Module
}

func (s *MessageService) GetTenantCurrentMonthSmsUsedNum(tenantId string) (int, error) {
	if tenantId == "" {
		return 0, errors.NewBusinessError("租户不能为空")
	}
	num := s.NotificationRecordDao.GetTenantSMSLackRecordNum(tenantId)
	return num, nil
}
