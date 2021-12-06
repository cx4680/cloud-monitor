package controllers

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/external/cbr"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/validator/translate"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/vo"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type CbrInstanceCtl struct {
}

func NewCbrCtl() *CbrInstanceCtl {
	return &CbrInstanceCtl{}
}

func (cic *CbrInstanceCtl) Page(c *gin.Context) {
	var params = CbrPageParam{}
	if err := c.ShouldBindQuery(params); err != nil {
		c.JSON(http.StatusBadRequest, translate.GetErrorMsg(err))
		return
	}
	var cbrParam = &cbr.QueryParam{
		TenantId:   params.TenantId,
		VaultId:    params.InstanceId,
		VaultName:  params.InstanceName,
		Status:     strings.Join(params.StatusList, ","),
		PageNumber: strconv.Itoa(params.Current),
		PageSize:   strconv.Itoa(params.PageSize),
	}
	cbrPageVO, err := cbr.PageList(cbrParam)
	if err != nil {
		c.JSON(http.StatusOK, global.NewSuccess("查询失败", err))
	}
	list := make([]interface{}, cbrPageVO.Total_count)
	if cbrPageVO != nil && cbrPageVO.Total_count > 0 {
		for index, data := range cbrPageVO.Data {
			instanceVO := CbrInfoVo{
				InstanceId:   data.VaultId,
				InstanceName: data.Name,
				Region:       data.Region,
				Status:       data.Status,
				Type:         data.Type,
				Capacity:     data.Capacity,
				UsedCapacity: data.UsedCapacity,
			}
			list[index] = instanceVO
		}
	}
	pageVO := vo.PageVO{
		Current: params.Current,
		Size:    params.PageSize,
		Total:   cbrPageVO.Total_count,
		Records: list,
	}
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", pageVO))
}

type CbrPageParam struct {
	TenantId     string   `form:"tenantId,omitempty"`
	InstanceId   string   `form:"instanceId,omitempty"`
	InstanceName string   `form:"instanceName,omitempty"`
	StatusList   []string `form:"statusList,omitempty"`
	Current      int      `form:"current,default=1,omitempty"`
	PageSize     int      `form:"pageSize,default=10,omitempty"`
}

type CbrInfoVo struct {
	InstanceId   string `json:"instanceId,omitempty"`
	InstanceName string `json:"instanceName,omitempty"`
	Region       string `json:"region,omitempty"`
	Status       string `json:"status,omitempty"`
	Type         string `json:"type,omitempty"`
	Capacity     string `json:"capacity,omitempty"`
	UsedCapacity string `json:"usedCapacity,omitempty"`
}
