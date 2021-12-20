package controllers

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/tools"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/validator/translate"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/vo"
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
	var f forms.AlertRecordPageQueryForm
	if err := c.ShouldBindJSON(&f); err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	tenantId := c.GetString(global.TenantId)
	page := dao.AlertRecord.GetPageList(global.DB, tenantId, f)
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", page))
}

func (a *AlertRecordController) GetDetail(c *gin.Context) {
	recordId := c.Query("recordId")
	if tools.IsBlank(recordId) {
		c.JSON(http.StatusBadRequest, global.NewError("参数异常"))
		return
	}
	c.JSON(http.StatusOK, dao.AlertRecord.GetById(global.DB, recordId))
}

func (a *AlertRecordController) GetAlertRecordTotal(c *gin.Context) {
	var f forms.AlertRecordPageQueryForm
	if err := c.ShouldBindQuery(&f); err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	tenantId := c.GetString(global.TenantId)
	d, _ := time.ParseDuration("24h")
	d7, _ := time.ParseDuration("-168h")
	var start, end string
	//没有传日期则计算7天内的数据
	if f.StartTime == "" || f.EndTime == "" {
		now := tools.GetNow()
		end = tools.TimeToStr(now.Add(d), tools.DayTimeFmt)
		start = tools.TimeToStr(now.Add(d7), tools.DayTimeFmt)
	} else {
		start = tools.TimeToStr(tools.StrToTime(tools.FullTimeFmt, f.StartTime), tools.DayTimeFmt)
		end = tools.TimeToStr(tools.StrToTime(tools.FullTimeFmt, f.EndTime).Add(d), tools.DayTimeFmt)
	}
	total := dao.AlertRecord.GetAlertRecordTotal(global.DB, tenantId, f.Region, start, end)
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", total))
}

func (a *AlertRecordController) GetRecordNumHistory(c *gin.Context) {
	var f forms.AlertRecordPageQueryForm
	if err := c.ShouldBindQuery(&f); err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	tenantId := c.GetString(global.TenantId)
	d, _ := time.ParseDuration("24h")
	startDate := tools.StrToTime(tools.FullTimeFmt, f.StartTime)
	endDate := tools.StrToTime(tools.FullTimeFmt, f.EndTime).Add(d)

	start := tools.TimeToStr(startDate, tools.DayTimeFmt)
	end := tools.TimeToStr(endDate, tools.DayTimeFmt)
	numList := dao.AlertRecord.GetRecordNumHistory(global.DB, tenantId, f.Region, start, end)
	//补充无数据的日期，该日期的历史数据为0
	resultMap := make(map[string]int)
	for _, v := range numList {
		resultMap[v.DayTime] = v.Number
	}
	var data []vo.RecordNumHistory
	for endDate.Sub(startDate) > 0 {
		recordNumHistory := vo.RecordNumHistory{
			DayTime: tools.TimeToStr(startDate, tools.DayTimeFmt),
			Number:  resultMap[tools.TimeToStr(startDate, tools.DayTimeFmt)],
		}
		data = append(data, recordNumHistory)
		startDate = startDate.Add(d)
	}
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", data))
}
