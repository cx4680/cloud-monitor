package consumer

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/enum"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

func MonitorProductHandler(msgs []*primitive.MessageExt) {
	for i := range msgs {
		var MqMsg form.MqMsg
		fmt.Printf("subscribe callback: %v \n", msgs[i])
		err := json.Unmarshal(msgs[i].Body, &MqMsg)
		if err != nil {
			logger.Logger().Error(err.Error())
		}
		switch MqMsg.EventEum {
		case enum.ChangeMonitorProductStatus:
			data := jsonutil.ToString(MqMsg.Data)
			var param form.MonitorProductParam
			jsonutil.ToObject(data, &param)
			dao.MonitorProduct.ChangeStatus(global.DB, param.BizIdList, param.Status)
		}
	}
}
