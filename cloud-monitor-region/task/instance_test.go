package task

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/vo"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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
	get, err := http.Get("http://10.255.191.82:30090/api/v1/query?query=" + url.QueryEscape("ebs_stock_all{instance=\"10.150.132.138:8080\"}"))
	if err != nil {
		fmt.Printf("%v", err)
	}
	fmt.Printf("%v", get)
}

func TestSafe2(t *testing.T) {
	apiUrl := "http://127.0.0.1:9090/get"
	// URL param
	data := url.Values{}
	data.Set("name", "小王子")
	data.Set("age", "18")
	u, err := url.ParseRequestURI(apiUrl)
	if err != nil {
		fmt.Printf("parse url requestUrl failed, err:%v\n", err)
	}
	u.RawQuery = data.Encode()
	fmt.Println(u.String())
	resp, _ := http.Get(u.String())
	fmt.Println(resp)
}
