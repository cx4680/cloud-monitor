package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"fmt"
	"net/http"
)

type ExternService struct {
}

func NewExternService() *ExternService {
	return &ExternService{}
}

func (externService *ExternService) GetRegionList(tenantId string) ([]*ResultBean, error) {
	path := fmt.Sprintf("%s?format=json&method=%s&appId=%s", config.GetCommonConfig().Nk, "QUERY_REGION_INFO", "CESTC_UNHQ_queryPoolsByLoginId")
	header := http.Header{}
	header.Add("userCode", tenantId)
	param := map[string]string{}
	param["loginId"] = tenantId
	json, err := tools.HttpHeaderPostJson(path, header, param)
	if err != nil {
		return nil, err
	}
	region := &RegionInfoDTO{}
	tools.ToObject(json, region)
	return region.result, nil
}

type RegionInfoDTO struct {
	respDesc string
	code     string
	message  string
	ok       bool
	respCode string
	result   []*ResultBean
}

type ResultBean struct {
	PoolUrl  string
	AllFlag  bool
	PoolId   string
	PoolName string
}
