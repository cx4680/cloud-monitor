package controller

import (
	commonDtos "code.cestc.cn/ccos-ops/cloud-monitor/business-common/dto"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/validator/translate"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/vo"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type AlarmRecordController struct {
}

func NewAlarmRecordController() *AlarmRecordController {
	return &AlarmRecordController{}
}

func (a *AlarmRecordController) GetPageList(c *gin.Context) {
	var f form.AlarmRecordPageQueryForm
	if err := c.ShouldBindJSON(&f); err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	tenantId, _ := util.GetTenantId(c)
	page := dao.AlarmRecord.GetPageList(global.DB, tenantId, f)
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", page))
}

func (a *AlarmRecordController) GetAlarmContactInfo(c *gin.Context) {
	bizId := c.Query("bizId")
	if strutil.IsBlank(bizId) {
		c.JSON(http.StatusBadRequest, global.NewError("参数错误"))
		return
	}
	tenantId, err := util.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, global.NewError("身份获取错误"))
		return
	}
	r := dao.AlarmRecord.GetByBizIdAndTenantId(global.DB, bizId, tenantId)
	if strutil.IsBlank(r.ContactInfo) {
		c.JSON(http.StatusOK, global.NewSuccess("查询成功", nil))
		return
	}
	var contactGroups []*commonDtos.ContactGroupInfo
	jsonutil.ToObject(r.ContactInfo, &contactGroups)
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", contactGroups))
}

func (a *AlarmRecordController) GetAlarmRecordTotal(c *gin.Context) {
	var f form.AlarmRecordPageQueryForm
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
	total := dao.AlarmRecord.GetAlarmRecordTotal(global.DB, tenantId, f.Region, start, end)
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", total))
}

func (a *AlarmRecordController) GetRecordNumHistory(c *gin.Context) {
	var f form.AlarmRecordPageQueryForm
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
	numList := dao.AlarmRecord.GetRecordNumHistory(global.DB, tenantId, f.Region, start, end)
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
