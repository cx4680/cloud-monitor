package eip

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/external"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/pkg/config"
	"encoding/json"
	"github.com/pkg/errors"
	"strings"
)

type QueryParam struct {
	IpAddress    string `json:"ipAddress,omitempty"`
	Status       int    `json:"status,omitempty"`
	InstanceType int    `json:"instanceType,omitempty"`
	Uid          string `json:"uid,omitempty"`
	RegionCode   string `json:"regionCode,omitempty"`
	StatusList   []int  `json:"statusList,omitempty"`
}

func GetEipInstancePage(form *QueryParam, pageIndex int, pageSize int, userCode string) (*QueryPageResult, error) {
	request := &external.QueryPageRequest{Data: form, PageIndex: pageIndex, PageSize: pageSize}
	resp, err := external.PageList(userCode, request, config.GetConfig().Nk+"?appId=600006&format=json&method=CESTC_UNHQ_getEipInfoList")
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
	Rows  []InfoBean
}

type InfoBean struct {
	BandWidth     int         `json:"bandWidth"`
	ReleaseTime   string      `json:"releaseTime"`
	Instance      interface{} `json:"instance"`
	Data          string      `json:"data"`
	Isp           string      `json:"isp"`
	Description   interface{} `json:"description"`
	PayModel      int         `json:"payModel"`
	FreezeTime    string      `json:"freezeTime"`
	UserCode      string      `json:"userCode"`
	InstanceUid   interface{} `json:"instanceUid"`
	Uid           string      `json:"uid"`
	RegionCode    string      `json:"regionCode"`
	ResourceCode  string      `json:"resourceCode"`
	ExpireDay     int         `json:"expireDay"`
	Id            int         `json:"id"`
	BandWidthUid  string      `json:"bandWidthUid"`
	FloatingIp    string      `json:"floatingIp"`
	InstanceType  interface{} `json:"instanceType"`
	IpAddress     string      `json:"ipAddress"`
	FloatingIpUid interface{} `json:"floatingIpUid"`
	PortUid       interface{} `json:"portUid"`
	ExpireTime    string      `json:"expireTime"`
	CreateTime    string      `json:"createTime"`
	Name          string      `json:"name"`
	Status        int         `json:"status"`
}
