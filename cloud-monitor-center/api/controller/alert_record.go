package controller

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/validator/translate"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/vo"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type AlertRecordController struct {
}

func NewAlertRecordController() *AlertRecordController {
	return &AlertRecordController{}
}

func (a *AlertRecordController) GetPageList(c *gin.Context) {
	var f form.AlertRecordPageQueryForm
	if err := c.ShouldBindJSON(&f); err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	tenantId, _ := util.GetTenantId(c)
	page := dao.AlertRecord.GetPageList(global.DB, tenantId, f)
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", page))
}

func (a *AlertRecordController) GetDetail(c *gin.Context) {
	recordId := c.Query("recordId")
	if strutil.IsBlank(recordId) {
		c.JSON(http.StatusBadRequest, global.NewError("参数异常"))
		return
	}
	tenantId, err := util.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, global.NewError("获取用户信息失败"))
		return
	}
	c.JSON(http.StatusOK, dao.AlertRecord.GetByIdAndTenantId(global.DB, recordId, tenantId))
}

func (a *AlertRecordController) GetAlertRecordTotal(c *gin.Context) {
	var f form.AlertRecordPageQueryForm
	if err := c.ShouldBindQuery(&f); err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	tenantId, _ := util.GetTenantId(c)
	d, _ := time.ParseDuration("24h")
	d7, _ := time.ParseDuration("-168h")
	var start, end string
	//没有传日期则计算7天内的数据
	if f.StartTime == "" || f.EndTime == "" {
		now := util.GetNow()
		end = util.TimeToStr(now.Add(d), util.DayTimeFmt)
		start = util.TimeToStr(now.Add(d7), util.DayTimeFmt)
	} else {
		start = util.TimeToStr(util.StrToTime(util.FullTimeFmt, f.StartTime), util.DayTimeFmt)
		end = util.TimeToStr(util.StrToTime(util.FullTimeFmt, f.EndTime).Add(d), util.DayTimeFmt)
	}
	total := dao.AlertRecord.GetAlertRecordTotal(global.DB, tenantId, f.Region, start, end)
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", total))
}

func (a *AlertRecordController) GetRecordNumHistory(c *gin.Context) {
	var f form.AlertRecordPageQueryForm
	if err := c.ShouldBindQuery(&f); err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	tenantId, _ := util.GetTenantId(c)
	d, _ := time.ParseDuration("24h")
	startDate := util.StrToTime(util.FullTimeFmt, f.StartTime)
	endDate := util.StrToTime(util.FullTimeFmt, f.EndTime).Add(d)

	start := util.TimeToStr(startDate, util.DayTimeFmt)
	end := util.TimeToStr(endDate, util.DayTimeFmt)
	numList := dao.AlertRecord.GetRecordNumHistory(global.DB, tenantId, f.Region, start, end)
	//补充无数据的日期，该日期的历史数据为0
	resultMap := make(map[string]int)
	for _, v := range numList {
		resultMap[v.DayTime] = v.Number
	}
	var data []vo.RecordNumHistory
	for endDate.Sub(startDate) > 0 {
		recordNumHistory := vo.RecordNumHistory{
			DayTime: util.TimeToStr(startDate, util.DayTimeFmt),
			Number:  resultMap[util.TimeToStr(startDate, util.DayTimeFmt)],
		}
		data = append(data, recordNumHistory)
		startDate = startDate.Add(d)
	}
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", data))
}
