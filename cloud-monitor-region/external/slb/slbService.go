package slb

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/external"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/pkg/config"
	"encoding/json"
	"github.com/pkg/errors"
	"strings"
)

type QueryParam struct {
	RegionCode  string   `json:"regionCode,omitempty"`
	Address     string   `json:"address,omitempty"`
	LbUid       string   `json:"lbUid,omitempty"`
	Name        string   `json:"name,omitempty"`
	Ip          string   `json:"ip,omitempty"`
	NetworkName string   `json:"networkName,omitempty"`
	NetworkUid  string   `json:"networkUid,omitempty"`
	StateList   []string `json:"stateList,omitempty"`
}

func GetSlbInstancePage(form *QueryParam, pageIndex int, pageSize int, userCode string) (*QueryPageResult, error) {
	request := &external.QueryPageRequest{Data: form, PageIndex: pageIndex, PageSize: pageSize}
	resp, err := external.PageList(userCode, request, config.GetConfig().Nk+"?appId=600006&format=json&method=CESTC_UNHQ_queryLVSMachineList")
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
	SubnetId    int         `json:"subnetId"`
	UnbindEip   interface{} `json:"unbindEip"`
	OldState    interface{} `json:"oldState"`
	NetworkName string      `json:"networkName"`
	CrmAttrList []struct {
		AttrCode  string `json:"attrCode"`
		AttrValue string `json:"attrValue"`
		AttrName  string `json:"attrName"`
	} `json:"crmAttrList"`
	Remark   string `json:"remark"`
	PayModel string `json:"payModel"`
	Eip      struct {
		ExpireTime    interface{} `json:"expireTime"`
		Bandwidth     int         `json:"bandwidth"`
		Ip            string      `json:"ip"`
		Name          interface{} `json:"name"`
		FloatingIpUid string      `json:"floatingIpUid"`
		CrmAttrList   []struct {
			AttrCode  string  `json:"attrCode"`
			AttrValue *string `json:"attrValue"`
			AttrName  string  `json:"attrName"`
		} `json:"crmAttrList"`
		PayModel interface{} `json:"payModel"`
		Id       int         `json:"id"`
	} `json:"eip"`
	PoolList      []interface{} `json:"poolList"`
	UserCode      string        `json:"userCode"`
	Spec          string        `json:"spec"`
	SubnetName    string        `json:"subnetName"`
	RegionCode    string        `json:"regionCode"`
	ResourceCode  string        `json:"resourceCode"`
	NetworkId     int           `json:"networkId"`
	Id            int           `json:"id"`
	State         string        `json:"state"`
	BandWidthUid  string        `json:"bandWidthUid"`
	Address       string        `json:"address"`
	LbUid         string        `json:"lbUid"`
	Ip            interface{}   `json:"ip"`
	StateList     interface{}   `json:"stateList"`
	EipInstaceUid interface{}   `json:"eipInstaceUid"`
	FloatingIpUid string        `json:"floatingIpUid"`
	UpdateTime    interface{}   `json:"updateTime"`
	PortUid       interface{}   `json:"portUid"`
	ExpireTime    string        `json:"expireTime"`
	CreateTime    string        `json:"createTime"`
	ListenerList  []struct {
		ListenerUid  string `json:"listenerUid"`
		Protocol     string `json:"protocol"`
		ProtocolPort int    `json:"protocolPort"`
		Name         string `json:"name"`
	} `json:"listenerList"`
	SubnetUid  string      `json:"subnetUid"`
	Name       string      `json:"name"`
	NetworkUid string      `json:"networkUid"`
	ZoneCode   interface{} `json:"zoneCode"`
}
