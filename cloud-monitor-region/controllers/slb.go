package controllers

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/external/slb"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/utils"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/validator/translate"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/vo"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SlbCtl struct {
}

func NewSlbCtl() *SlbCtl {
	return &SlbCtl{}
}

// Page
// @Summary Page
// @Schemes
// @Description GetById
// @Tags InstanceCtl
// @Accept json
// @Produce json
// @Param id query  string true "id"
// @Success 200 {object} vo.InstanceVO
// @Router /hawkeye/eip/page [get]
func (ic *SlbCtl) Page(c *gin.Context) {
	var params = &SlbPageParam{}
	if err := c.ShouldBindQuery(params); err != nil {
		c.JSON(http.StatusBadRequest, translate.GetErrorMsg(err))
		return
	}
	slbParams := slb.QueryParam{Address: params.PrivateIp, Name: params.InstanceName, StateList: params.StatusList, LbUid: params.InstanceId}
	tenantId, exists := c.Get("tenantId")
	if !exists {
		c.JSON(http.StatusBadRequest, "tenantId not exists")
		return
	}
	ret, err := slb.GetSlbInstancePage(&slbParams, params.Current, params.PageSize, tenantId.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, translate.GetErrorMsg(err))
		return
	}
	pageVO := vo.PageVO{Total: ret.Total, Size: params.PageSize, Current: params.Current}
	if ret != nil && ret.Total > 0 {
		list := make([]interface{}, ret.Total)
		for index, row := range ret.Rows {
			instanceVO := SlbInfoVo{
				InstanceId:   row.LbUid,
				InstanceName: row.Name,
				Spec:         row.Spec,
				Status:       row.State,
				PrivateIp:    row.Address,
				Id:           row.Id,
				Eip: Eip{
					Id: row.FloatingIpUid,
					Ip: utils.If(&row.Eip != nil, row.Eip.Ip, nil).(string),
				},
				Vpc: Vpc{Id: row.NetworkUid, Name: row.NetworkName},
			}
			listenerList := row.ListenerList
			listeners := make([]Listener, len(listenerList))
			instanceVO.ListenerList = listeners
			for index, listener := range listenerList {
				newListener := Listener{
					Name:         listener.Name,
					Id:           listener.ListenerUid,
					Protocol:     listener.Protocol,
					ProtocolPort: listener.ProtocolPort,
				}
				listeners[index] = newListener
			}
			list[index] = instanceVO
		}
		pageVO.Records = list
	}
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", pageVO))
}

type SlbPageParam struct {
	InstanceId   string   `json:"instanceId"`
	InstanceName string   `json:"instanceName"`
	PrivateIp    string   `json:"privateIp"`
	StatusList   []string `json:"statusList"`
	Current      int      `form:"current,default=1"`
	PageSize     int      `form:"pageSize,default=10"`
}
type SlbInfoVo struct {
	Id           int        `json:"id"`
	InstanceId   string     `json:"instanceId"`
	InstanceName string     `json:"instanceName"`
	Status       string     `json:"status"`
	Spec         string     `json:"spec"`
	PrivateIp    string     `json:"privateIp"`
	Eip          Eip        `json:"eip"`
	Vpc          Vpc        `json:"vpc"`
	ListenerList []Listener `json:"listenerList"`
}
type Eip struct {
	Ip string `json:"ip"`
	Id string `json:"id"`
}
type Vpc struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
type Listener struct {
	Name         string `json:"name"`
	Id           string `json:"id"`
	Protocol     string `json:"protocol"`
	ProtocolPort int    `json:"protocolPort"`
}
