package external

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
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
}

type SlbQueryPageRequest struct {
	OrderName string      `json:"orderName,omitempty"`
	OrderRule string      `json:"OrderRule,omitempty"`
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

func (slb *SlbInstanceService) convertRealForm(form service.InstancePageForm) interface{} {
	queryParam := SlbQueryParam{
		Address:   form.ExtraAttr["privateIp"],
		LbUid:     form.InstanceId,
		Name:      form.InstanceName,
		StateList: form.StatusList,
	}
	return SlbQueryPageRequest{
		PageIndex: form.Current,
		PageSize:  form.PageSize,
		Data:      queryParam,
		TenantId:  form.TenantId,
	}
}

func (slb *SlbInstanceService) doRequest(url string, form interface{}) (interface{}, error) {
	var f = form.(SlbQueryPageRequest)
	respStr, err := tools.HttpPostJson(url, form, map[string]string{"userCode": f.TenantId})
	if err != nil {
		return nil, err
	}
	var resp SlbResponse
	tools.ToObject(respStr, &resp)
	return resp, nil
}

func (slb *SlbInstanceService) convertResp(realResp interface{}) (int, []service.InstanceCommonVO) {
	vo := realResp.(SlbResponse)
	var list []service.InstanceCommonVO
	if vo.Data.Total > 0 {
		for _, d := range vo.Data.Rows {
			list = append(list, service.InstanceCommonVO{
				Id:   string(rune(d.Id)),
				Name: d.Name,
				Labels: []service.InstanceLabel{{
					Name:  "subnetName",
					Value: d.SubnetName,
				}, {
					Name:  "status",
					Value: d.State,
				}},
			})
		}
	}
	return vo.Data.Total, list
}
