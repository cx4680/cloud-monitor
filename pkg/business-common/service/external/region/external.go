package region

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/httputil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"fmt"
)

type ExternService struct {
}

var RegionService = NewExternService()

func NewExternService() *ExternService {
	return &ExternService{}
}

func (externService *ExternService) GetRegionList(tenantId string) ([]ConfigItemVO, error) {
	path := fmt.Sprintf("%s?format=json&method=%s&appId=%s", config.Cfg.Common.Nk, "CESTC_UNHQ_queryPoolsByLoginId", "600006")
	json, err := httputil.HttpPostJson(path, map[string]string{"loginId": tenantId}, map[string]string{"userCode": tenantId})
	if err != nil {
		return nil, err
	}
	region := &InfoDTO{}
	jsonutil.ToObject(json, region)
	var configItemVOList []ConfigItemVO
	for _, v := range region.Result {
		configItemVO := ConfigItemVO{
			Code: v.PoolId,
			Name: v.PoolName,
			Data: v.PoolUrl,
		}
		configItemVOList = append(configItemVOList, configItemVO)
	}
	return configItemVOList, nil
}

type InfoDTO struct {
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

type ConfigItemVO struct {
	Id     string `json:"id"`
	BizId  string `json:"bizId"`
	PBizId string `json:"pBizId"` //配置名称
	Name   string `json:"name"`   //配置编码
	Code   string `json:"code"`   //配置编码
	Data   string `json:"data"`   //配置值
	Remark string `json:"remark"` //备注
}
