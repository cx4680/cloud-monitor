package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/constants"
	commonDao "code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dtos"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/sysComponent/sysRocketMq"
	commonModels "code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service/external/messageCenter"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/utils/snowflake"
	"github.com/google/uuid"
	"log"
	"strconv"
	"strings"
	"time"
)

type MessageService struct {
	NotificationRecordDao *commonDao.NotificationRecordDao
	MCS                   *messageCenter.Service
}

func NewMessageService(notificationRecordDao *commonDao.NotificationRecordDao, messageCenterService *messageCenter.Service) *MessageService {
	return &MessageService{
		NotificationRecordDao: notificationRecordDao,
		MCS:                   messageCenterService,
	}
}

func (s *MessageService) TargetFilter(targetList []messageCenter.MessageTargetDTO, senderId string, isCenter bool) {
	if targetList == nil || len(targetList) <= 0 {
		return
	}
	var nt []messageCenter.MessageTargetDTO
	for _, t := range targetList {
		if messageCenter.Email == t.Type {
			nt = append(nt, t)
		} else {
			num := s.NotificationRecordDao.GetTenantPhoneCurrentMonthRecordNum(senderId)
			if s.checkSentNum(senderId, num, isCenter) {
				nt = append(nt, t)
			} else {
				log.Printf("too many records have been sent, send refused, sender=%s \n", senderId)
			}
		}
	}
	targetList = nt
}

type AlertMsgSendDTO struct {
	AlertId string
	Msg     messageCenter.MessageSendDTO
}

func (s *MessageService) SendMsgNew(msgList []AlertMsgSendDTO, isCenter bool) error {
	if !config.GetCommonConfig().HasNoticeModel {
		log.Println("There is no message center for this env")
		return nil
	}
	if msgList == nil || len(msgList) <= 0 {
		return nil
	}
	var recordList []commonModels.NotificationRecord
	for _, m := range msgList {
		msg := m.Msg
		//TODO 确认切片引用
		s.TargetFilter(msg.Target, msg.SenderId, isCenter)

		//send msg
		if err := s.MCS.Send(msg); err != nil {
			log.Printf("message send error, %v\n\n", err)
			return err
		}
		for _, t := range msg.Target {
			recordList = append(recordList, commonModels.NotificationRecord{
				SenderId:         msg.SenderId,
				SourceId:         m.AlertId,
				SourceType:       int(msg.SourceType),
				TargetAddress:    t.Addr,
				NotificationType: int(t.Type),
				Result:           1,
			})
		}
	}
	//save record local
	s.NotificationRecordDao.InsertBatch(global.DB, recordList)

	//	sync record to center
	if !isCenter {
		//TODO topic
		sysRocketMq.SendMsg("notification_sync", tools.ToString(recordList))
	}
	return nil

}

func (s *MessageService) checkSentNum(tenantId string, num int, isCenter bool) bool {
	//check local
	if num > constants.MaxSmsNum {
		log.Println("user ", tenantId, " already used more", constants.MaxSmsNum, ", send sms refused.")
		return false
	}
	//	check remote
	if !isCenter {
		num = s.getUserCurrentMonthSmsUsedNum(tenantId)
		if num > constants.MaxSmsNum {
			log.Println("user ", tenantId, " already used more", constants.MaxSmsNum, ", send sms refused.")
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
	params["msgUsed"] = string(rune(num))
	params["msgLeft"] = string(rune(constants.MaxSmsNum - num))
	params["msgInitial"] = string(rune(constants.MaxSmsNum))

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
	s.sendToMsgCenter(noticeMsgDTOList)

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

func (s *MessageService) sendToMsgCenter(noticeMsgDTOS []*dtos.NoticeMsgDTO) {
	smsMessageDTO := buildSmsMessageReqDTO(noticeMsgDTOS)
	s.doSendToMsgCenter(smsMessageDTO)
}

func buildSmsMessageReqDTO(noticeMsgDTOS []*dtos.NoticeMsgDTO) *dtos.SmsMessageReqDTO {
	var msgList []dtos.MessagesBean
	var collect map[dtos.MsgEvent][]*dtos.NoticeMsgDTO
	for _, noticeMsgDTO := range noticeMsgDTOS {
		if collect[noticeMsgDTO.MsgEvent] == nil {
			var msgDTOS []*dtos.NoticeMsgDTO
			collect[noticeMsgDTO.MsgEvent] = msgDTOS
		}
		collect[noticeMsgDTO.MsgEvent] = append(collect[noticeMsgDTO.MsgEvent], noticeMsgDTO)
	}
	for event, noticeMsxgDTOS1 := range collect {
		var list []dtos.RecvObjectBean
		for _, dto := range noticeMsxgDTOS1 {
			list = append(list, dto.RevObjectBean)
		}
		msgList = append(msgList, dtos.MessagesBean{
			MsgEventCode:   string(rune(event.Type)) + string(rune(event.Source)),
			RecvObjectList: list,
		})
	}
	return &dtos.SmsMessageReqDTO{
		BusinessId: strings.ReplaceAll(uuid.New().String(), "-", ""),
		InModeCode: constants.AppCode,
		Messages:   msgList,
		ReferTime:  "",
	}
}

func (s *MessageService) doSendToMsgCenter(smsMessageReqDTO *dtos.SmsMessageReqDTO) {

	resp, err := tools.HttpPostJson(config.GetCommonConfig().SmsCenterPath, *smsMessageReqDTO)
	if err != nil {
		log.Fatal("send message to msgCenter fail", err)
	}
	log.Println("send message to msgCenter resp=" + resp)
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
