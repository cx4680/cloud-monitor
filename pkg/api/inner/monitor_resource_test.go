package inner

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/httputil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/model"
	"fmt"
	"log"
	"testing"
)

func Test_GetResourceList(t *testing.T) {
	respStr, err := httputil.HttpGet("http://cloud-monitor.vmcbc06.intranet.cecloudcs.com/inner/monitorResource/list?abbreviation=slb&tenantId=220506072800500")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(respStr)
}

func Test_test(t *testing.T) {
	str := `{
		"errorCode": "0",
		"success": true,
		"module":[
			{
			"Id":1,
			"instance_name":"test"
			},
			{
			"Id":2,
			"instance_name":"jim"
			}
		]
	}`

	var r global.Resp
	err := jsonutil.ToObjectWithError(str, &r)
	if err != nil {
		log.Fatal(err)
		return
	}
	if list, ok := r.Module.([]interface{}); ok {
		for _, temp := range list {
			var ai model.AlarmInstance
			err = jsonutil.ToObjectWithError(jsonutil.ToString(temp), &ai)
			if err != nil {
				log.Fatal(err)
				return
			}
			fmt.Println(ai.Id)
		}
	}

}
