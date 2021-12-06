package controllers

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/vo"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/external/nat"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/validator/translate"
	"github.com/gin-gonic/gin"
	"net/http"
)

type NatCtl struct {
}

func NewNatCtl() *NatCtl {
	return &NatCtl{}
}

func (ic *NatCtl) Page(c *gin.Context) {
	var params = &NatPageParam{}
	if err := c.ShouldBindQuery(params); err != nil {
		c.JSON(http.StatusBadRequest, translate.GetErrorMsg(err))
		return
	}
	natParams := nat.QueryParam{Name: params.InstanceName,Uid: params.InstanceId}
	tenantId, exists := c.Get("tenantId")
	if !exists {
		c.JSON(http.StatusBadRequest, "tenantId not exists")
		return
	}
	ret, err := nat.GetNatInstancePage(&natParams, params.Current, params.PageSize, tenantId.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, translate.GetErrorMsg(err))
		return
	}
	pageVO := vo.PageVO{Total: ret.Total, Size: params.PageSize, Current: params.Current}
	pageVO.Records = ret.Rows
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", pageVO))
}

type NatPageParam struct {
	InstanceId   string   `form:"instanceId,omitempty"`
	InstanceName string   `form:"instanceName,omitempty"`
	PrivateIp    string   `form:"privateIp,omitempty"`
	StatusList   []string `form:"statusList,omitempty"`
	Current      int      `form:"current,default=1"`
	PageSize     int      `form:"pageSize,default=10"`
}