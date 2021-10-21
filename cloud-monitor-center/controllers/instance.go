package controllers

import (
	"code.cestc.cn/ccos-ops/cloud-monitor-center/dao"
	"github.com/gin-gonic/gin"
)

type InstanceCtl struct {
	dao *dao.InstanceDao
}

func NewInstanceCtl(dao *dao.InstanceDao) *InstanceCtl {
	return &InstanceCtl{dao: dao}
}

func (ctl *InstanceCtl) Page(c *gin.Context) {

}

func (ctl *InstanceCtl) Unbind(c *gin.Context) {

}

func (ctl *InstanceCtl) Bind(c *gin.Context) {

}

func (ctl *InstanceCtl) GetRuleList(c *gin.Context) {

}
