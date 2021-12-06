package ecs

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/vo"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type UserInstanceQueryForm struct {
	tenantId string

	instanceId string

	instanceName string

	status int

	statusList []int

	current int

	pageSize int
}

type VmParams struct {
	HostId     string `json:"hostId,omitempty"`
	HostName   string `json:"hostName,omitempty"`
	Status     int    `json:"status,omitempty"`
	StatusList []int  `json:"statusList,omitempty"`
}

func PageList(form *forms.EcsQueryPageForm) (*vo.EcsPageVO, error) {
	var ecsInnerGateway = config.GetEcsConfig().InnerGateway
	path := ecsInnerGateway + "/noauth/ecs/PageList"
	jsonStr, _ := json.Marshal(form)
	resp, err := http.Post(path, "application/json", strings.NewReader(string(jsonStr)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	result, _ := ioutil.ReadAll(resp.Body)
	var ecsQueryPageVO vo.EcsQueryPageVO
	json.Unmarshal(result, &ecsQueryPageVO)
	return &ecsQueryPageVO.Data, nil
}

type Response struct {
	Code string
	Msg  string
	Data QueryPageResult
}
type QueryPageResult struct {
	Total int
	Rows  []VmInfoBean
}

type VmInfoBean struct {
	Id                   int
	HostId               string
	HostName             string
	Ip                   string
	Status               int
	TemplateSpecInfoBean *TemplateSpecInfoBean
}

type TemplateSpecInfoBean struct {
	RegionCode string
}
