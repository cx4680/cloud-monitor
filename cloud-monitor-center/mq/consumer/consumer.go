package consumer

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

func InstanceHandler(msgs []*primitive.MessageExt) {
	alarmInstanceDao := dao.AlarmInstance
	for i := range msgs {
		var instances []*models.AlarmInstance
		fmt.Printf("subscribe callback: %v \n", msgs[i])
		err := json.Unmarshal(msgs[i].Body, &instances)
		if err != nil {
			continue
		}
		alarmInstanceDao.UpdateBatchInstanceName(instances)
	}
}
