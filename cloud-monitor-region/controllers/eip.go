package controllers

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/external/eip"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/validator/translate"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/vo"
	"github.com/gin-gonic/gin"
	"net/http"
)

type EipCtl struct {
}

func NewEipCtl() *EipCtl {
	return &EipCtl{}
}

// Page
// @Summary Page
// @Schemes
// @Description GetById
// @Tags EipCtl
// @Accept json
// @Produce json
// @Param id query  string true "id"
// @Success 200 {object} vo.InstanceVO
// @Router /hawkeye/eip/page [get]
func (ic *EipCtl) Page(c *gin.Context) {
	var params = &EipPageParam{}
	if err := c.ShouldBindQuery(params); err != nil {
		c.JSON(http.StatusBadRequest, translate.GetErrorMsg(err))
		return
	}
	eipParams := eip.QueryParam{IpAddress: params.Ip, Uid: params.InstanceId, StatusList: params.StatusList}
	tenantId, exists := c.Get("tenantId")
	if !exists {
		c.JSON(http.StatusBadRequest, "tenantId not exists")
		return
	}
	ret, err := eip.GetEipInstancePage(&eipParams, params.Current, params.PageSize, tenantId.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, translate.GetErrorMsg(err))
		return
	}
	pageVO := vo.PageVO{Total: ret.Total, Size: params.PageSize, Current: params.Current}
	if ret != nil && ret.Total > 0 {
		list := make([]interface{}, ret.Total)
		for index, row := range ret.Rows {
			instanceVO := EipInfoVo{
				InstanceId:     row.Uid,
				InstanceName:   row.Name,
				EipAddress:     row.IpAddress,
				Status:         row.Status,
				BandWidth:      row.BandWidth,
				BindInstanceId: row.BandWidthUid,
			}
			list[index] = instanceVO
		}
		pageVO.Records = list
	}
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", pageVO))
}

type EipPageParam struct {
	InstanceId string `form:"instanceId"`
	Ip         string `form:"ip"`
	Status     string `form:"status"`
	StatusList []int  `form:"statusList"`
	Current    int    `form:"current,default=1"`
	PageSize   int    `form:"pageSize,default=10"`
}
type EipInfoVo struct {
	EipAddress     string `json:"eipAddress,omitempty"`
	InstanceName   string `json:"instanceName,omitempty"`
	InstanceId     string `json:"instanceId,omitempty"`
	Status         int    `json:"status,omitempty"`
	BandWidth      int    `json:"bandWidth,omitempty"`
	BindInstanceId string `json:"bindInstanceId,omitempty"`
}
