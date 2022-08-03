package service

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/common/config"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/httputil"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/jsonutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/dao"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/errors"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/model"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/util"
	cvo "code.cestc.cn/ccos-ops/cloud-monitor/pkg/business-common/vo"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/vo"
	"context"
	"gorm.io/gorm"
	"strconv"
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

func (s *AlarmRecordService) GetLevelTotalByIam(param form.AlarmRecordPageQueryForm) ([]*form.AlarmRecordNum, error) {
	start, end := s.getFmtTime(param.StartTime, param.EndTime)
	directoryIdList, err := GetIamDirectoryIdList(param.IamUserId, param.TenantId)
	if err != nil {
		return nil, err
	}
	resourcesIdList, err := GetIamResourcesIdList(directoryIdList)
	if err != nil {
		return nil, err
	}
	if len(resourcesIdList) == 0 {
		return nil, nil
	}
	return dao.AlarmRecord.GetLevelTotalByIam(global.DB, resourcesIdList, start, end), nil
}

func (s *AlarmRecordService) GetRecordNumHistoryByIam(param form.AlarmRecordPageQueryForm) ([]vo.RecordNumHistory, error) {
	d, _ := time.ParseDuration("24h")
	startDate := util.StrToTime(util.FullTimeFmt, param.StartTime)
	endDate := util.StrToTime(util.FullTimeFmt, param.EndTime).Add(d)
	start := util.TimeToStr(startDate, util.DayTimeFmt)
	end := util.TimeToStr(endDate, util.DayTimeFmt)
	directoryIdList, err := GetIamDirectoryIdList(param.IamUserId, param.TenantId)
	if err != nil {
		return nil, err
	}
	resourcesIdList, err := GetIamResourcesIdList(directoryIdList)
	if err != nil {
		return nil, err
	}
	if len(resourcesIdList) == 0 {
		return nil, nil
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

func (s *AlarmRecordService) GetProductRecordNumHistoryByIam(param form.AlarmRecordPageQueryForm) ([]*form.ProductAlarmRecordNum, error) {
	start, end := s.getFmtTime(param.StartTime, param.EndTime)
	directoryIdList, err := GetIamDirectoryIdList(param.IamUserId, param.TenantId)
	if err != nil {
		return nil, err
	}
	resourcesIdList, err := GetIamResourcesIdList(directoryIdList)
	if err != nil {
		return nil, err
	}
	if len(resourcesIdList) == 0 {
		return nil, nil
	}
	list := dao.AlarmRecord.GetProductRecordNumHistoryByIam(global.DB, resourcesIdList, start, end)
	if len(list) > 10 {
		num := 0
		for _, v := range list[10:] {
			num += v.Count
		}
		list = append(list[:10], &form.ProductAlarmRecordNum{ProductCode: "other", Count: num})
	}
	return list, nil
}

func (s *AlarmRecordService) GetPageListByIam(param form.AlarmRecordPageQueryForm) (*cvo.PageVO, error) {
	start, end := s.getFmtTime(param.StartTime, param.EndTime)
	directoryIdList, err := GetIamDirectoryIdList(param.IamUserId, param.TenantId)
	if err != nil {
		return nil, err
	}
	resourcesIdList, err := GetIamResourcesIdList(directoryIdList)
	if err != nil {
		return nil, err
	}
	if len(resourcesIdList) == 0 {
		return nil, nil
	}
	page, total := dao.AlarmRecord.GetPageListByIam(global.DB, param.ProductCode, resourcesIdList, start, end, param.PageNum, param.PageSize)
	return &cvo.PageVO{
		Records: page,
		Total:   int(total),
		Size:    param.PageSize,
		Current: param.PageNum,
		Pages:   (int(total) / param.PageSize) + 1,
	}, nil
}

func (s *AlarmRecordService) getFmtTime(startTime, endTime string) (string, string) {
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

func GetIamDirectoryIdList(iamUserId, belongLoginId string) ([]string, error) {
	response, err := httputil.HttpPostJson(config.Cfg.Common.IamRetrieveIamUser, map[string]string{"principalId": iamUserId, "belongAccountUid": belongLoginId}, nil)
	if err != nil {
		logger.Logger().Errorf("获取iam部门错误：%v", err)
		return nil, errors.NewBusinessError("获取iam部门错误")
	}
	var iamDirectory form.IamDirectory
	jsonutil.ToObject(response, &iamDirectory)
	var directoryIdList = []string{strconv.Itoa(iamDirectory.Module.DirectoryId)}
	for _, v := range iamDirectory.Module.ChildList {
		directoryIdList = append(directoryIdList, strconv.Itoa(v.DirectoryId))
	}
	return directoryIdList, nil
}

func GetIamResourcesIdList(directoryIds []string) ([]string, error) {
	param := form.InstanceRequest{
		DirectoryIds: directoryIds,
		CurrPage:     "1",
		PageSize:     "99999",
	}
	response, err := httputil.HttpPostJson(config.Cfg.Common.Rc, param, nil)
	if err != nil {
		logger.Logger().Errorf("获取实例列表错误：%v", err)
		return nil, errors.NewBusinessError("获取实例列表错误")
	}
	var result form.InstanceResponse
	jsonutil.ToObject(response, &result)
	var resourcesIdList []string
	for _, v := range result.Data.List {
		resourcesIdList = append(resourcesIdList, v.ResourceId)
	}
	return resourcesIdList, nil
}
