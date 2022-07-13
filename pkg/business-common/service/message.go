package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/snowflake"
	commonDao "code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/dto"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/errors"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	commonModels "code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/model"
	message_center2 "code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/service/external/message_center"
	util2 "code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/constant"
	"strconv"
)

type MessageService struct {
	NotificationRecordDao *commonDao.NotificationRecordDao
	MCS                   *message_center2.Service
}

type AlertMsgSendDTO struct {
	AlertId    string
	SenderId   string
	SourceType message_center2.MsgSource
	Msgs       []AlertMsgDTO
}

type AlertMsgDTO struct {
	Type    message_center2.ReceiveType
	Targets []string
	Content string
}

func NewMessageService(messageCenterService *message_center2.Service) *MessageService {
	return &MessageService{
		NotificationRecordDao: commonDao.NotificationRecord,
		MCS:                   messageCenterService,
	}
}

func (s *MessageService) TargetFilter(targetList []string, rt message_center2.ReceiveType, senderId string) []string {
	if targetList == nil || len(targetList) <= 0 {
		return nil
	}
	var addrs []string
	for _, t := range targetList {
		if message_center2.Email == rt {
			addrs = append(addrs, t)
			continue
		}
		//短信需要验证是否超限
		num := s.NotificationRecordDao.GetTenantPhoneCurrentMonthRecordNum(senderId)
		if s.checkSentNum(senderId, num) {
			addrs = append(addrs, t)
		} else {
			logger.Logger().Infof("too many records have been sent, send refused, sender=%s \n", senderId)
		}
	}
	//去重
	addrs = util2.RemoveDuplicateElement(addrs)
	return addrs
}

func (s *MessageService) SendAlarmNotice(msgList []interface{}) error {
	channels := s.MCS.GetRemoteChannels()
	if len(channels) == 0 {
		logger.Logger().Info("There is no message center for this env")
		return nil
	}

	if msgList == nil || len(msgList) <= 0 {
		return nil
	}

	var recordList []commonModels.NotificationRecord
	var sendMsgList []message_center2.MessageSendDTO
	//统计短信数量
	var smsSender []string
	for _, alertMsg := range msgList {
		am := alertMsg.(*AlertMsgSendDTO)
		for _, msg := range am.Msgs {
			newTargets := s.TargetFilter(msg.Targets, msg.Type, am.SenderId)
			sendMsgList = append(sendMsgList, message_center2.MessageSendDTO{
				SenderId:   am.SenderId,
				Type:       msg.Type,
				SourceType: am.SourceType,
				Targets:    newTargets,
				Content:    msg.Content,
			})
			if msg.Type == message_center2.Phone {
				smsSender = append(smsSender, msg.Targets...)
			}
			for _, addr := range newTargets {
				recordList = append(recordList, commonModels.NotificationRecord{
					BizId:            strconv.FormatInt(snowflake.GetWorker().NextId(), 10),
					SenderId:         am.SenderId,
					SourceId:         am.AlertId,
					SourceType:       uint8(am.SourceType),
					TargetAddress:    addr,
					NotificationType: uint8(msg.Type),
					Result:           1,
					CreateTime:       util2.GetNow(),
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

	//	sync record
	//擦除自增Id，解决同步到中心化Id冲突问题
	for i, _ := range recordList {
		recordList[i].Id = 0
	}
	// 发送短信余量不足提醒
	if len(smsSender) > 0 {
		smsSender = util2.RemoveDuplicateElement(smsSender)
	}
	return nil
}

// SendActivateMsg 发送激活信息
func (s *MessageService) SendActivateMsg(msg message_center2.MessageSendDTO, contactId string) {
	//send msg
	var ret uint8 = 1
	if err := s.MCS.Send(msg); err != nil {
		logger.Logger().Errorf("message send error, %v\n\n", err)
		ret = 0
	}
	record := &commonModels.NotificationRecord{
		BizId:            strconv.FormatInt(snowflake.GetWorker().NextId(), 10),
		SenderId:         msg.SenderId,
		SourceId:         contactId,
		SourceType:       uint8(message_center2.VERIFY),
		TargetAddress:    msg.Targets[0],
		NotificationType: uint8(msg.Type),
		Result:           ret,
		CreateTime:       util2.GetNow(),
	}
	s.NotificationRecordDao.Insert(global.DB, record)
	record.Id = 0
}

// SmsMarginReminder 短信余量提醒
func (s *MessageService) SmsMarginReminder(sender string) {
	count := s.NotificationRecordDao.GetTenantSMSLackRecordNum(sender)
	if count > 0 {
		//已发送过短信提醒
		return
	}
	alreadySendNum := s.NotificationRecordDao.GetTenantPhoneCurrentMonthRecordNum(sender)

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

	remainderMsg := message_center2.MessageSendDTO{
		SenderId:   sender,
		Type:       message_center2.Phone,
		Targets:    []string{serialNumber},
		SourceType: message_center2.SMS_LACK,
		Content:    jsonutil.ToString(params),
	}
	if err := s.MCS.Send(remainderMsg); err != nil {
		return
	}
	//保存发送记录
	record := &commonModels.NotificationRecord{
		SenderId:         sender,
		SourceId:         "sms-lack-" + sender,
		SourceType:       uint8(message_center2.SMS_LACK),
		TargetAddress:    serialNumber,
		NotificationType: uint8(message_center2.Phone),
		Result:           1,
		CreateTime:       util2.GetNow(),
	}
	s.NotificationRecordDao.Insert(global.DB, record)

	record.Id = 0
}

func (s *MessageService) checkSentNum(tenantId string, num int) bool {
	//check local
	if num > constant.MaxSmsNum {
		logger.Logger().Info("user ", tenantId, " already used more", constant.MaxSmsNum, ", send sms refused.")
		return false
	}
	return true
}

func (s *MessageService) GetTenantCurrentMonthSmsUsedNum(tenantId string) (int, error) {
	if tenantId == "" {
		return 0, errors.NewBusinessError("租户不能为空")
	}
	num := s.NotificationRecordDao.GetTenantSMSLackRecordNum(tenantId)
	return num, nil
}
