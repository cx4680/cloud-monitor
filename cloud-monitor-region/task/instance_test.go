package task

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/vo"
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
	"unsafe"
)

func TestDelete(t *testing.T) {
	count := 10
	dbInstanceList := make([]*models.AlarmInstance, count)
	for count > 0 {
		count--
		dbInstanceList[count] = &models.AlarmInstance{
			InstanceID:   "id" + strconv.Itoa(count),
			InstanceName: "name" + strconv.Itoa(count),
		}
	}

	count2 := 5
	instanceInfoList := make([]*models.AlarmInstance, count2)
	for count2 > 0 {
		count2--
		instanceInfoList[count2] = &models.AlarmInstance{
			InstanceID:   "id" + strconv.Itoa(count2),
			InstanceName: "name" + strconv.Itoa(count2),
		}
	}
	DeleteNotExistsInstances("xx", dbInstanceList, instanceInfoList)
}

func TestDelete1(t *testing.T) {
	s := "xxx"
	d := vo.PageVO{
		Records: unsafe.Pointer(&s),
	}
	marshal, err := json.Marshal(d)
	if err != nil {

	}
	fmt.Printf("%+v", string(marshal))
}

func TestSafe1(t *testing.T) {
	d := vo.PageVO{
		Records: "gg",
	}
	marshal, err := json.Marshal(d)
	if err != nil {

	}
	fmt.Printf("%+v", string(marshal))
}
