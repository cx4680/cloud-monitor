package ecs

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/external"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"encoding/json"
	"errors"
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

func GetUserInstancePage(form *VmParams, pageIndex int, pageSize int, userCode string) (*QueryPageResult, error) {
	request := &external.QueryPageRequest{Data: form, PageIndex: pageIndex, PageSize: pageSize}
	resp, err := external.PageList(userCode, request, config.GetCommonConfig().Nk+"?appId=600006&format=json&method=CESTC_UNHQ_pageQueryServerLIst")
	if err != nil {
		return nil, err
	}
	result := &Response{}
	if json.Unmarshal(resp, result); err != nil {
		logger.Logger().Errorf("check result parse  failed, err:%v\n", err)
		return nil, err
	}
	if strings.EqualFold(result.Code, "0") {
		return &result.Data, nil
	}
	return nil, errors.New(result.Msg)
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
