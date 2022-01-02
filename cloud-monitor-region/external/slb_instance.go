package external

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/httputil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"strconv"
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
	Ip            string        `json:"ip"`
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
				}, {
					Name:  "id",
					Value: strconv.Itoa(d.Id),
				}},
			})
		}
	}
	return vo.Data.Total, list
}

func getListenerList(slb *SlbInfoBean) string {
	var listener []string
	for _, v := range slb.ListenerList {
		listener = append(listener, v.Name)
	}
	return strings.Join(listener, ",")
}
