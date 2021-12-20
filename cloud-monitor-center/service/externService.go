package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/vo"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"fmt"
)

type ExternService struct {
}

func NewExternService() *ExternService {
	return &ExternService{}
}

func (externService *ExternService) GetRegionList(tenantId string) ([]vo.ConfigItemVO, error) {
	path := fmt.Sprintf("%s?format=json&method=%s&appId=%s", config.GetCommonConfig().Nk, "CESTC_UNHQ_queryPoolsByLoginId", "600006")
	json, err := tools.HttpPostJson(path, map[string]string{"loginId": tenantId}, map[string]string{"userCode": tenantId})
	if err != nil {
		return nil, err
	}
	region := &RegionInfoDTO{}
	tools.ToObject(json, region)
	var configItemVOList []vo.ConfigItemVO
	for _, v := range region.Result {
		configItemVO := vo.ConfigItemVO{
			Code: v.PoolId,
			Name: v.PoolName,
			Data: v.PoolUrl,
		}
		configItemVOList = append(configItemVOList, configItemVO)
	}
	return configItemVOList, nil
}

type RegionInfoDTO struct {
	RespDesc string        `json:"respDesc"`
	Code     string        `json:"code"`
	Message  string        `json:"message"`
	Ok       bool          `json:"ok"`
	RespCode string        `json:"respCode"`
	Result   []*ResultBean `json:"result"`
}

type ResultBean struct {
	PoolUrl  string `json:"poolUrl"`
	AllFlag  bool   `json:"allFlag"`
	PoolId   string `json:"poolId"`
	PoolName string `json:"poolName"`
}
