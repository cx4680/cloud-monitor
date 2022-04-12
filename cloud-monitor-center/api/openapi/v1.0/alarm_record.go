package v1_0

import (
	commonDtos "code.cestc.cn/ccos-ops/cloud-monitor/business-common/dto"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/openapi"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global/openapi/channel_type"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/cloud-monitor-center/vo"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
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
	reqParam := AlarmRecordPageQueryForm{
		PageNumber: 1,
		PageSize:   10,
	}
	if err := c.ShouldBindQuery(&reqParam); err != nil {
		logger.Logger().Info(err)
		c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.GetErrorCode(err), c))
		return
	}
	tenantId, err := util.GetTenantId(c)
	if err != nil {
		logger.Logger().Info(err)
		c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.InvalidParameter, c))
		return
	}
	voParam := form.AlarmRecordPageQueryForm{
		PageNum:    reqParam.PageNumber,
		PageSize:   reqParam.PageSize,
		Level:      reqParam.AlarmLevel,
		ResourceId: reqParam.ResourceId,
		StartTime:  FormatTime(reqParam.StartTime),
		EndTime:    FormatTime(reqParam.EndTime),
		RuleName:   reqParam.RuleName,
		Status:     reqParam.Status,
	}
	pageVo := dao.AlarmRecord.GetPageList(global.DB, tenantId, voParam)
	result := AlarmHistoryPage{
		ResCommonPage: *openapi.NewResCommonPage(c, pageVo),
		Histories:     nil,
	}
	records := pageVo.Records.([]vo.AlarmRecordPageVO)
	for _, item := range records {
		history := AlarmHistory{
			BizId:        item.BizId,
			Status:       item.Status,
			RuleId:       item.RuleId,
			RuleName:     item.RuleName,
			MonitorType:  item.MonitorType,
			SourceType:   item.SourceType,
			SourceId:     item.SourceId,
			CurrentValue: item.CurrentValue,
			StartTime:    item.StartTime,
			EndTime:      item.EndTime,
			TargetValue:  item.TargetValue,
			Expression:   item.Expression,
			Duration:     item.Duration,
			Level:        item.Level,
			MetricCode:   item.AlarmKey,
			CreateTime:   item.CreateTime.Format(util.FullTimeFmt),
		}
		result.Histories = append(result.Histories, history)
	}
	c.JSON(http.StatusOK, result)
}
func FormatTime(s time.Time) string {
	if s.IsZero() {
		return ""
	}
	return util.TimeToFullTimeFmtStr(s)
}
func (a *AlarmRecordController) GetAlarmContactInfo(c *gin.Context) {
	bizId := c.Param("AlarmBizId")
	tenantId, err := util.GetTenantId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, openapi.NewRespError(openapi.InvalidParameter, c))
		return
	}
	c.Set(global.ResourceName, bizId)
	r := dao.AlarmRecord.GetByBizIdAndTenantId(global.DB, bizId, tenantId)
	if strutil.IsBlank(r.ContactInfo) {
		c.JSON(http.StatusOK, openapi.NewResSuccess(c))
		return
	}
	var contactGroupsVo []*commonDtos.ContactGroupInfo
	jsonutil.ToObject(r.ContactInfo, &contactGroupsVo)
	var result []ContactGroupInfo
	for _, item := range contactGroupsVo {
		group := ContactGroupInfo{
			RequestId: openapi.GetRequestId(c),
			GroupId:   item.GroupId,
			GroupName: item.GroupName,
		}
		for _, contactVo := range item.Contacts {
			user := UserContactInfo{
				ContactId:   contactVo.ContactId,
				ContactName: contactVo.ContactName,
			}
			if len(contactVo.Mail) != 0 {
				user.Channels = append(user.Channels, struct {
					Channel channel_type.ChannelType
					Address string
				}{Channel: channel_type.Email, Address: contactVo.Mail})
			}
			if len(contactVo.Phone) != 0 {
				user.Channels = append(user.Channels, struct {
					Channel channel_type.ChannelType
					Address string
				}{Channel: channel_type.Phone, Address: contactVo.Phone})
			}
			group.Contacts = append(group.Contacts, user)
		}
		result = append(result, group)
	}

	c.JSON(http.StatusOK, result)
}

type AlarmRecordPageQueryForm struct {
	PageNumber int
	PageSize   int
	StartTime  time.Time `time_format:"2006-01-02 15:04:05"`
	EndTime    time.Time `binding:"gtefield=StartTime" time_format:"2006-01-02 15:04:05"`
	AlarmLevel string    `binding:"oneof=1  2 3 4 '' "`
	ResourceId string
	RuleName   string
	Status     string `binding:"oneof=resolved  firing '' "`
}

type AlarmHistoryPage struct {
	openapi.ResCommonPage
	Histories []AlarmHistory
}

type AlarmHistory struct {
	BizId        string
	Status       string
	RuleId       string
	RuleName     string
	MonitorType  string
	SourceType   string
	SourceId     string
	CurrentValue string
	StartTime    string
	EndTime      string
	TargetValue  string
	Expression   string
	Duration     string
	Level        int
	MetricCode   string
	CreateTime   string
}

type ContactGroupInfo struct {
	RequestId string
	GroupId   string
	GroupName string
	Contacts  []UserContactInfo
}

type UserContactInfo struct {
	ContactId   string
	ContactName string
	Channels    []struct {
		Channel channel_type.ChannelType
		Address string
	}
}
