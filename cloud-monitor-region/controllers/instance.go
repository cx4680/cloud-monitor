package controllers

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/external/ecs"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/utils"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/validator/translate"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-region/vo"
	"github.com/gin-gonic/gin"
	"net/http"
)

type InstanceCtl struct {
	dao *dao.InstanceDao
}

func NewInstanceCtl(dao *dao.InstanceDao) *InstanceCtl {
	return &InstanceCtl{dao}
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
// @Router /hawkeye/instance/page [get]
func (ic *InstanceCtl) Page(c *gin.Context) {
	var params = &UserInstancePageQueryForm{}
	if err := c.ShouldBindQuery(params); err != nil {
		c.JSON(http.StatusBadRequest, translate.GetErrorMsg(err))
		return
	}
	vmParams := ecs.VmParams{HostId: params.InstanceId, HostName: params.InstanceName, Status: params.Status, StatusList: params.StatusList}
	tenantId, exists := c.Get("tenantId")
	if !exists {
		c.JSON(http.StatusBadRequest, "tenantId not exists")
		return
	}
	ret, err := ecs.GetUserInstancePage(&vmParams, params.Current, params.PageSize, tenantId.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, translate.GetErrorMsg(err))
		return
	}
	pageVO := vo.PageVO{Total: ret.Total, Size: params.PageSize, Current: params.Current}
	if ret != nil && ret.Total > 0 {
		list := make([]interface{}, ret.Total)
		for index, row := range ret.Rows {
			instanceVO := InstanceVO{
				Id:           row.Id,
				InstanceId:   row.HostId,
				InstanceName: row.HostName,
				Ip:           row.Ip,
				Status:       row.Status,
			}
			if row.TemplateSpecInfoBean != nil {
				instanceVO.Region = row.TemplateSpecInfoBean.RegionCode
			}
			list[index] = instanceVO
		}
		pageVO.Records = list
	}
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", pageVO))

}

// GetInstanceNumByRegion
// @Summary GetInstanceNumByRegion
// @Schemes
// @Description GetInstanceNumByRegion
// @Tags InstanceCtl
// @Accept json
// @Produce json
// @Param id query  string true "id"
// @Success 200 {object} vo.AlarmInstanceRegionVO
// @Router /hawkeye/instance/page [get]
func (ic *InstanceCtl) GetInstanceNumByRegion(ctx *gin.Context) {
	tenantId, exists := ctx.Get("tenantId")
	if !exists {
		ctx.JSON(http.StatusBadRequest, "tenantId not exists")
		return
	}
	vmParams := ecs.VmParams{}
	ret, err := ecs.GetUserInstancePage(&vmParams, 1, 10, tenantId.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, translate.GetErrorMsg(err))
		return
	}
	bindNum := ic.dao.GetInstanceNum(tenantId.(string))
	total := utils.If(bindNum > ret.Total, bindNum, ret.Total)
	vo := &AlarmInstanceRegionVO{Total: total.(int), BindNum: bindNum}
	ctx.JSON(http.StatusOK, global.NewSuccess("查询成功", vo))
}

type UserInstancePageQueryForm struct {
	InstanceId   string `form:"instanceId,omitempty"`
	InstanceName string `form:"instanceName,omitempty"`
	Status       int    `form:"status,omitempty"`
	StatusList   []int  `form:"statusList,omitempty"`
	Current      int    `form:"current,omitempty,default=1"`
	PageSize     int    `form:"pageSize,omitempty,default=10"`
}

type InstanceVO struct {
	Id           int
	InstanceId   string
	InstanceName string
	Region       string
	Ip           string

	/**
	 * 0: 待处理 1: 启动中 2: 运行中 3: 关机中 4: 停止 5: 重启中 6: 异常 7: 删除 12:未识别, 其他: 未知
	 */
	Status int
}

type AlarmInstanceRegionVO struct {
	Total int

	/**
	 * 已绑定实例数
	 */
	BindNum int
}
