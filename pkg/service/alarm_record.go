package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/model"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/util"
	cvo "code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/vo"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/vo"
	"context"
	"gorm.io/gorm"
	"time"
)

type AlarmRecordService struct {
	AlarmRecordDao *dao.AlarmRecordDao
	AlarmInfoDao   *dao.AlarmInfoDao
}

func NewAlarmRecordService() *AlarmRecordService {
	return &AlarmRecordService{AlarmRecordDao: dao.AlarmRecord, AlarmInfoDao: dao.AlarmInfo}
}

func (s *AlarmRecordService) InsertAndHandler(ctx *context.Context, recordList []model.AlarmRecord, infoList []model.AlarmInfo, eventList []interface{}) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		if len(recordList) > 0 {
			s.AlarmRecordDao.InsertBatch(tx, recordList)
		}
		if len(infoList) > 0 {
			s.AlarmInfoDao.InsertBatch(tx, infoList)
		}
		//告警处置
		if !AlarmHandlerQueue.PushFrontBatch(eventList) {
			logger.Logger().Error("requestId=", util.GetRequestId(ctx), ", add to alarm handler fail, handlerMap=", jsonutil.ToString(eventList))
		}
		return nil
	})
}

func (s *AlarmRecordService) GetAlarmRecordTotalByIam(param form.AlarmRecordPageQueryForm) (int64, error) {
	resourcesIdList, err := GetIamResourcesIdList(param.IamUserId)
	if err != nil {
		return 0, err
	}
	start, end := getFmtTime(param.StartTime, param.EndTime)
	return dao.AlarmRecord.GetAlarmRecordTotalByIam(global.DB, resourcesIdList, start, end), nil
}

func (s *AlarmRecordService) GetLevelTotal(param form.AlarmRecordPageQueryForm) ([]*form.AlarmRecordNum, error) {
	start, end := getFmtTime(param.StartTime, param.EndTime)
	isIamLogin := CheckIamLogin(param.TenantId, param.IamUserId)
	if isIamLogin {
		resourcesIdList, err := GetIamResourcesIdList(param.IamUserId)
		if err != nil {
			return nil, err
		}
		if len(resourcesIdList) == 0 {
			return nil, nil
		}
		return dao.AlarmRecord.GetLevelTotalByIam(global.DB, resourcesIdList, start, end), nil
	}
	return dao.AlarmRecord.GetLevelTotal(global.DB, param.TenantId, param.Region, start, end), nil
}

func (s *AlarmRecordService) GetRecordNumHistoryByIam(param form.AlarmRecordPageQueryForm) ([]vo.RecordNumHistory, error) {
	resourcesIdList, err := GetIamResourcesIdList(param.IamUserId)
	if err != nil {
		return nil, err
	}
	if len(resourcesIdList) == 0 {
		return nil, nil
	}
	var start, end string
	var startDate, endDate time.Time
	d, _ := time.ParseDuration("24h")
	if strutil.IsBlank(param.StartTime) || strutil.IsBlank(param.EndTime) {
		start, end = getFmtTime(param.StartTime, param.EndTime)
		startDate = util.StrToTime(util.FullTimeFmt, start+" 00:00:00")
		endDate = util.StrToTime(util.FullTimeFmt, end+" 00:00:00")
	} else {
		startDate = util.StrToTime(util.FullTimeFmt, param.StartTime)
		endDate = util.StrToTime(util.FullTimeFmt, param.EndTime).Add(d)
		start = util.TimeToStr(startDate, util.DayTimeFmt)
		end = util.TimeToStr(endDate, util.DayTimeFmt)
	}
	numList := dao.AlarmRecord.GetRecordNumHistoryByIam(global.DB, resourcesIdList, start, end)
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
	return data, nil
}

func (s *AlarmRecordService) GetTotalByProduct(param form.AlarmRecordPageQueryForm) ([]*form.ProductAlarmRecordNum, error) {
	start, end := getFmtTime(param.StartTime, param.EndTime)
	isIamLogin := CheckIamLogin(param.TenantId, param.IamUserId)
	var list []*form.ProductAlarmRecordNum
	if isIamLogin {
		resourcesIdList, err := GetIamResourcesIdList(param.IamUserId)
		if err != nil {
			return nil, err
		}
		if len(resourcesIdList) == 0 {
			return nil, nil
		}
		list = dao.AlarmRecord.GetProductTotalByIam(global.DB, resourcesIdList, start, end)

	} else {
		list = dao.AlarmRecord.GetTotalByProduct(global.DB, param.TenantId, param.Region, start, end)
	}
	if len(list) > 10 {
		num := 0
		for _, v := range list[10:] {
			num += v.Count
		}
		list = append(list[:10], &form.ProductAlarmRecordNum{ProductCode: "other", Count: num})
	}
	return list, nil
}

func (s *AlarmRecordService) GetPageListByProduct(param form.AlarmRecordPageQueryForm) (*cvo.PageVO, error) {
	isIamLogin := CheckIamLogin(param.TenantId, param.IamUserId)
	start, end := getFmtTime(param.StartTime, param.EndTime)
	var page []*form.AlarmRecordPage
	var total int64
	if isIamLogin {
		resourcesIdList, err := GetIamResourcesIdList(param.IamUserId)
		if err != nil {
			return nil, err
		}
		if len(resourcesIdList) == 0 {
			return nil, nil
		}
		page, total = dao.AlarmRecord.GetPageListByProductByIam(global.DB, param.ProductCode, resourcesIdList, start, end, param.PageNum, param.PageSize)
	} else {
		page, total = dao.AlarmRecord.GetPageListByProduct(global.DB, param.ProductCode, param.TenantId, param.Region, start, end, param.PageNum, param.PageSize)
	}
	for i, v := range page {
		page[i].FmtTime = util.TimeToFullTimeFmtStr(v.Time)
	}
	return &cvo.PageVO{
		Records: page,
		Total:   int(total),
		Size:    param.PageSize,
		Current: param.PageNum,
		Pages:   (int(total) / param.PageSize) + 1,
	}, nil
}

func getFmtTime(startTime, endTime string) (string, string) {
	d, _ := time.ParseDuration("24h")
	d7, _ := time.ParseDuration("-168h")
	var start, end string
	//没有传日期则计算7天内的数据
	if startTime == "" || endTime == "" {
		now := util.GetNow()
		end = util.TimeToStr(now.Add(d), util.DayTimeFmt)
		start = util.TimeToStr(now.Add(d7), util.DayTimeFmt)
	} else {
		start = util.TimeToStr(util.StrToTime(util.FullTimeFmt, startTime), util.DayTimeFmt)
		end = util.TimeToStr(util.StrToTime(util.FullTimeFmt, endTime).Add(d), util.DayTimeFmt)
	}
	return start, end
}
