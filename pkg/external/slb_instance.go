package external

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/httputil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/service"
	"strings"
)

type SlbInstanceService struct {
	service.InstanceServiceImpl
}

type SlbQueryParam struct {
	LbUid     string   `json:"lbUid,omitempty"`
	Name      string   `json:"name,omitempty"`
	SlbId     string   `json:"slbId,omitempty"`
	SlbName   string   `json:"slbName,omitempty"`
	Address   string   `json:"address,omitempty"`
	StateList []string `json:"stateList,omitempty"`
	TenantId  string   `json:"tenantId,omitempty"`
}

type SlbQueryPageRequest struct {
	OrderName string      `json:"orderName,omitempty"`
	OrderRule string      `json:"orderRule,omitempty"`
	PageIndex int         `json:"pageIndex,omitempty"`
	PageSize  int         `json:"pageSize,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	//临时传递
	TenantId string
	IamInfo  service.IamInfo `json:"-"`
}

type SlbResponse struct {
	Code string             `json:"code"`
	Msg  string             `json:"msg"`
	Data SlbQueryPageResult `json:"data"`
}
type SlbQueryPageResult struct {
	Total int            `json:"total"`
	Rows  []*SlbInfoBean `json:"rows"`
}
type SlbInfoBean struct {
	LbUid       string `json:"lbUid"`
	Name        string `json:"name"`
	SlbName     string `json:"slbName"`
	SlbId       string `json:"slbId"`
	State       string `json:"state"`
	Address     string `json:"address"`
	NetworkName string `json:"networkName"`
	NetworkUid  string `json:"networkUid"`
	EipIp       string `json:"eipIp"`
	Spec        string `json:"spec"`
	Eip         struct {
		Ip         string      `json:"ip"`
		Name       interface{} `json:"name"`
		Bandwidth  int         `json:"bandwidth"`
		ExpireTime string      `json:"expireTime"`
		PayModel   interface{} `json:"payModel"`
		EipUid     string      `json:"eipUid"`
	} `json:"eip"`
	Listeners    []*Listener `json:"listeners"`
	ListenerList []*Listener `json:"listenerList"`
}

type Listener struct {
	Protocol     string `json:"protocol"`
	ProtocolPort int    `json:"protocolPort"`
	ListenerUid  string `json:"listenerUid"`
	ListenerName string `json:"listenerName"`
}

func (slb *SlbInstanceService) ConvertRealForm(form service.InstancePageForm) interface{} {
	queryParam := SlbQueryParam{
		Address:  form.ExtraAttr["privateIp"],
		SlbId:    form.InstanceId,
		SlbName:  form.InstanceName,
		TenantId: form.TenantId,
	}
	if strutil.IsNotBlank(form.StatusList) {
		queryParam.StateList = strings.Split(form.StatusList, ",")
	}
	return SlbQueryPageRequest{
		PageIndex: form.Current,
		PageSize:  form.PageSize,
		Data:      queryParam,
	}
}

func (slb *SlbInstanceService) DoRequest(url string, form interface{}) (interface{}, error) {
	var f = form.(SlbQueryPageRequest)
	respStr, err := httputil.HttpPostJson(url, f, nil)
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
				InstanceId:   d.SlbId,
				InstanceName: d.SlbName,
				Labels: []service.InstanceLabel{{
					Name:  "eipIp",
					Value: d.EipIp,
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
					Name:  "status",
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

func (slb *SlbInstanceService) ConvertRealAuthForm(form service.InstancePageForm) interface{} {
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
		IamInfo:   form.IamInfo,
	}
}

func (slb *SlbInstanceService) DoAuthRequest(url string, form interface{}) (interface{}, error) {
	var f = form.(SlbQueryPageRequest)
	respStr, err := httputil.HttpPostJson(url, form, slb.GetIamHeader(&f.IamInfo))
	if err != nil {
		return nil, err
	}
	var resp SlbResponse
	err = jsonutil.ToObjectWithError(respStr, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (slb *SlbInstanceService) ConvertAuthResp(realResp interface{}) (int, []service.InstanceCommonVO) {
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
					Name:  "status",
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
	if len(slb.Listeners) != 0 {
		for _, v := range slb.Listeners {
			listener = append(listener, v.ListenerName)
		}
	} else if len(slb.ListenerList) != 0 {
		for _, v := range slb.ListenerList {
			listener = append(listener, v.ListenerName)
		}
	}
	return strings.Join(listener, ",")
}
