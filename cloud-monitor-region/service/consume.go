package service

import (
	commonDtos "code.cestc.cn/ccos-ops/cloud-monitor/business-common/dto"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/enum/handler_type"
	commonService "code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service/external/message_center"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/httputil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"context"
	"time"

	"github.com/toolkits/pkg/concurrent/semaphore"
)

func StartHandleAlarmEvent(ctx context.Context) {
	messageSvc := commonService.NewMessageService(message_center.NewService())
	sema := semaphore.NewSemaphore(10)
	duration := time.Duration(100) * time.Millisecond
	for {
		events := AlarmHandlerQueue.PopBackBy(100)
		if len(events) == 0 {
			time.Sleep(duration)
			continue
		}
		consume(events, messageSvc, sema)
	}
}

func consume(events []interface{}, messageSvc *commonService.MessageService, sema *semaphore.Semaphore) {
	for _, event := range events {
		e := event.(AlarmHandlerEvent)

		sema.Acquire()
		go func(e AlarmHandlerEvent) {
			defer func() {
				if err := recover(); err != nil {
					logger.Logger().Error("requestId=", e.RequestId, ", run time error,", err)
				}
			}()
			defer sema.Release()

			consumeOneEvent(messageSvc, e)

		}(e)
	}

}

func consumeOneEvent(messageSvc *commonService.MessageService, e AlarmHandlerEvent) {
	t := e.Type
	if t == handler_type.Sms || t == handler_type.Email {
		if err := messageSvc.SendAlarmNotice([]interface{}{e.Data}); err != nil {
			logger.Logger().Error("requestId=", e.RequestId, "send alarm message fail, data=", jsonutil.ToString(e), err)
		}
	} else if t == handler_type.Http {
		//调用弹性伸缩
		data := e.Data.(*commonDtos.AutoScalingData)
		respJson, err := httputil.HttpPostJson(data.Param, map[string]string{"ruleId": data.RuleId, "tenantId": data.TenantId}, nil)
		if err != nil {
			logger.Logger().Error("requestId=", e.RequestId, ", autoScaling request fail, data=", jsonutil.ToString(e), err)
		} else {
			logger.Logger().Info("requestId=", e.RequestId, ", autoScaling request success, resp=", respJson)
		}
	} else {
		logger.Logger().Error("data type error, ", jsonutil.ToString(e))
	}
}
