package nat

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/external"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"encoding/json"
	"errors"
	"strings"
)

type QueryParam struct {
	Name string
	Uid  string
}

func GetNatInstancePage(form *QueryParam, pageIndex int, pageSize int, userCode string) (*QueryPageResult, error) {
	request := &external.QueryPageRequest{Data: form, PageIndex: pageIndex, PageSize: pageSize}
	resp, err := external.PageList(userCode, request, config.GetCommonConfig().Nk+"?appId=600006&format=json&method=CESTC_UNHQ_natPage")
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
	Rows  []*InfoBean
}

type InfoBean struct {
	VpcName       string
	ReleaseTime   string
	OrderId       string
	Specification string
	Description   string
	PayModel      int
	UpdateTime    string
	FreezeTime    string
	UserCode      string
	SubnetName    string
	EipNum        int
	SubnetCidr    string
	Uid           string
	VpcUid        string
	ExpireTime    string
	ResourceCode  string
	StatusList    string
	ExpireDay     int
	CreateTime    string
	SubnetUid     string
	Name          string
	Id            int
	AdminStateUp  bool
	Status        string
}
