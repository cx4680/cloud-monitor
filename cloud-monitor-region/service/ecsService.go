package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/vo"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

var ecsInnerGateway = config.GetEcsConfig().InnerGateway

func EcsPageList(form forms.EcsQueryPageForm) vo.EcsPageVO {
	path := ecsInnerGateway + "/noauth/ecs/PageList"
	jsonStr, _ := json.Marshal(form)
	resp, err := http.Post(path, "application/json", strings.NewReader(string(jsonStr)))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	result, _ := ioutil.ReadAll(resp.Body)
	var ecsQueryPageVO vo.EcsQueryPageVO
	json.Unmarshal(result, &ecsQueryPageVO)
	return ecsQueryPageVO.Data
}
