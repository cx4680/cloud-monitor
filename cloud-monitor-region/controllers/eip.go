package controllers

import (
	"github.com/gin-gonic/gin"
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
// @Tags InstanceCtl
// @Accept json
// @Produce json
// @Param id query  string true "id"
// @Success 200 {object} vo.InstanceVO
// @Router /hawkeye/eip/page [get]
func (ic *EipCtl) Page(ctx *gin.Context) {

}
