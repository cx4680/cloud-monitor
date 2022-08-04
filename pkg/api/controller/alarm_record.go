package controller

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	commonDtos "code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/dto"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	util2 "code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/service"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/validator/translate"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/vo"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type AlarmRecordController struct {
	service *service.AlarmRecordService
}

func NewAlarmRecordController() *AlarmRecordController {
	return &AlarmRecordController{service.NewAlarmRecordService()}
}

func (a *AlarmRecordController) GetPageList(c *gin.Context) {
	var f form.AlarmRecordPageQueryForm
	if err := c.ShouldBindJSON(&f); err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	tenantId, _ := util2.GetTenantId(c)
	page := dao.AlarmRecord.GetPageList(global.DB, tenantId, f)
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", page))
}

func (a *AlarmRecordController) GetAlarmContactInfo(c *gin.Context) {
	bizId := c.Query("bizId")
	if strutil.IsBlank(bizId) {
		c.JSON(http.StatusBadRequest, global.NewError("参数错误"))
		return
	}
	tenantId, err := util2.GetTenantId(c)
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
	tenantId, iamUserId, err := util2.GetTenantIdAndUserId(c)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
		return
	}
	isIamLogin := service.CheckIamLogin(tenantId, iamUserId)
	if isIamLogin {
		result, err := a.service.GetAlarmRecordTotalByIam(f)
		if err != nil {
			c.JSON(http.StatusOK, global.NewError(err.Error()))
			return
		}
		c.JSON(http.StatusOK, global.NewSuccess("查询成功", result))
		return
	}
	d, _ := time.ParseDuration("24h")
	d7, _ := time.ParseDuration("-168h")
	var start, end string
	//没有传日期则计算7天内的数据
	if f.StartTime == "" || f.EndTime == "" {
		now := util2.GetNow()
		end = util2.TimeToStr(now.Add(d), util2.DayTimeFmt)
		start = util2.TimeToStr(now.Add(d7), util2.DayTimeFmt)
	} else {
		start = util2.TimeToStr(util2.StrToTime(util2.FullTimeFmt, f.StartTime), util2.DayTimeFmt)
		end = util2.TimeToStr(util2.StrToTime(util2.FullTimeFmt, f.EndTime).Add(d), util2.DayTimeFmt)
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
	tenantId, iamUserId, err := util2.GetTenantIdAndUserId(c)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
		return
	}
	isIamLogin := service.CheckIamLogin(tenantId, iamUserId)
	if isIamLogin {
		f.TenantId = tenantId
		f.IamUserId = iamUserId
		result, err := a.service.GetRecordNumHistoryByIam(f)
		if err != nil {
			c.JSON(http.StatusOK, global.NewError(err.Error()))
			return
		}
		c.JSON(http.StatusOK, global.NewSuccess("查询成功", result))
		return
	}
	d, _ := time.ParseDuration("24h")
	startDate := util2.StrToTime(util2.FullTimeFmt, f.StartTime)
	endDate := util2.StrToTime(util2.FullTimeFmt, f.EndTime).Add(d)

	start := util2.TimeToStr(startDate, util2.DayTimeFmt)
	end := util2.TimeToStr(endDate, util2.DayTimeFmt)
	numList := dao.AlarmRecord.GetRecordNumHistory(global.DB, tenantId, f.Region, start, end)
	//补充无数据的日期，该日期的历史数据为0
	resultMap := make(map[string]int)
	for _, v := range numList {
		resultMap[v.DayTime] = v.Number
	}
	var data []vo.RecordNumHistory
	for endDate.Sub(startDate) > 0 {
		recordNumHistory := vo.RecordNumHistory{
			DayTime: util2.TimeToStr(startDate, util2.DayTimeFmt),
			Number:  resultMap[util2.TimeToStr(startDate, util2.DayTimeFmt)],
		}
		data = append(data, recordNumHistory)
		startDate = startDate.Add(d)
	}
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", data))
}

func (a *AlarmRecordController) GetLevelTotal(c *gin.Context) {
	tenantId, iamUserId, err := util2.GetTenantIdAndUserId(c)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
		return
	}
	var f form.AlarmRecordPageQueryForm
	if err = c.ShouldBindQuery(&f); err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	f.TenantId = tenantId
	f.IamUserId = iamUserId
	result, err := a.service.GetLevelTotal(f)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", result))
}

func (a *AlarmRecordController) GetTotalByProduct(c *gin.Context) {
	tenantId, iamUserId, err := util2.GetTenantIdAndUserId(c)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
		return
	}
	var f form.AlarmRecordPageQueryForm
	if err = c.ShouldBindQuery(&f); err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	f.TenantId = tenantId
	f.IamUserId = iamUserId
	result, err := a.service.GetTotalByProduct(f)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", result))
}

func (a *AlarmRecordController) GetPageListByProduct(c *gin.Context) {
	tenantId, iamUserId, err := util2.GetTenantIdAndUserId(c)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
		return
	}
	var f = form.AlarmRecordPageQueryForm{PageNum: 1, PageSize: 10}
	if err = c.ShouldBindJSON(&f); err != nil {
		c.JSON(http.StatusBadRequest, global.NewError(translate.GetErrorMsg(err)))
		return
	}
	f.TenantId = tenantId
	f.IamUserId = iamUserId
	page, err := a.service.GetPageListByProduct(f)
	if err != nil {
		c.JSON(http.StatusOK, global.NewError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, global.NewSuccess("查询成功", page))
}
