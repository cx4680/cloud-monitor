package controllers

import "github.com/gin-gonic/gin"

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
func (ic *SlbCtl) Page(ctx *gin.Context) {

}
