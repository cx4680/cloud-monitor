package external

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/httputil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"strings"
)

type SlbInstanceService struct {
	service.InstanceServiceImpl
}

type SlbQueryParam struct {
	RegionCode  string   `json:"regionCode,omitempty"`
	Address     string   `json:"address,omitempty"`
	LbUid       string   `json:"lbUid,omitempty"`
	Name        string   `json:"name,omitempty"`
	Ip          string   `json:"ip,omitempty"`
	NetworkName string   `json:"networkName,omitempty"`
	NetworkUid  string   `json:"networkUid,omitempty"`
	StateList   []string `json:"stateList,omitempty"`
	TenantId    string   `json:"tenantId,omitempty"`
}

type SlbQueryPageRequest struct {
	OrderName string      `json:"orderName,omitempty"`
	OrderRule string      `json:"orderRule,omitempty"`
	PageIndex int         `json:"pageIndex,omitempty"`
	PageSize  int         `json:"pageSize,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	//临时传递
	TenantId string
}

type SlbResponse struct {
	Code string
	Msg  string
	Data SlbQueryPageResult
}
type SlbQueryPageResult struct {
	Total int
	Rows  []*SlbInfoBean
}

type SlbInfoBean struct {
	Name         string        `json:"name"`
	UserCode     string        `json:"userCode"`
	LbUid        string        `json:"lbUid"`
	State        string        `json:"state"`
	Address      string        `json:"address"`
	SubnetUid    string        `json:"subnetUid"`
	PortUid      interface{}   `json:"portUid"`
	RegionCode   string        `json:"regionCode"`
	ZoneCode     interface{}   `json:"zoneCode"`
	Remark       string        `json:"remark"`
	CreateTime   string        `json:"createTime"`
	UpdateTime   interface{}   `json:"updateTime"`
	NetworkName  string        `json:"networkName"`
	NetworkUid   string        `json:"networkUid"`
	SubnetName   string        `json:"subnetName"`
	PoolList     []interface{} `json:"poolList"`
	ListenerList []struct {
		Protocol     string `json:"protocol"`
		ProtocolPort int    `json:"protocolPort"`
		ListenerUid  string `json:"listenerUid"`
		ListenerName string `json:"listenerName"`
	} `json:"listenerList"`
	Eip struct {
		Ip         string      `json:"ip"`
		Name       interface{} `json:"name"`
		Bandwidth  int         `json:"bandwidth"`
		ExpireTime string      `json:"expireTime"`
		PayModel   interface{} `json:"payModel"`
		EipUid     int         `json:"eipUid"`
	} `json:"eip"`
	ExpireTime string `json:"expireTime"`
	PayModel   string `json:"payModel"`
	OrderId    string `json:"orderId"`
	Spec       string `json:"spec"`
}

func (slb *SlbInstanceService) ConvertRealForm(form service.InstancePageForm) interface{} {
	queryParam := SlbQueryParam{
		Address:  form.ExtraAttr["privateIp"],
		LbUid:    form.InstanceId,
		Name:     form.InstanceName,
		TenantId: form.TenantId,
	}
	if strutil.IsNotBlank(form.StatusList) {
		queryParam.StateList = strings.Split(form.StatusList, ",")
	}
	return SlbQueryPageRequest{
		PageIndex: form.Current,
		PageSize:  form.PageSize,
		Data:      queryParam,
		TenantId:  form.TenantId,
	}
}

func (slb *SlbInstanceService) DoRequest(url string, form interface{}) (interface{}, error) {
	var f = form.(SlbQueryPageRequest)
	respStr, err := httputil.HttpPostJson(url, form, map[string]string{"userCode": f.TenantId})
	if err != nil {
		return nil, err
	}
	var resp SlbResponse
	jsonutil.ToObject(respStr, &resp)
	return resp, nil
}

func (slb *SlbInstanceService) ConvertResp(realResp interface{}) (int, []service.InstanceCommonVO) {
	vo := realResp.(SlbResponse)
	var list []service.InstanceCommonVO
	if vo.Data.Total > 0 {
		for _, d := range vo.Data.Rows {
			list = append(list, service.InstanceCommonVO{
				InstanceId:   d.LbUid,
				InstanceName: d.Name,
				Labels: []service.InstanceLabel{{
					Name:  "eipIp",
					Value: d.Eip.Ip,
				}, {
					Name:  "privateIp",
					Value: d.Address,
				}, {
					Name:  "vpcName",
					Value: d.NetworkName,
				}, {
					Name:  "vpcId",
					Value: d.NetworkUid,
				}, {
					Name:  "state",
					Value: d.State,
				}, {
					Name:  "spec",
					Value: d.Spec,
				}, {
					Name:  "listener",
					Value: getListenerList(d),
				}},
			})
		}
	}
	return vo.Data.Total, list
}

func getListenerList(slb *SlbInfoBean) string {
	var listener []string
	for _, v := range slb.ListenerList {
		listener = append(listener, v.ListenerName)
	}
	return strings.Join(listener, ",")
}
