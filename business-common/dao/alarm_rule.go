package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/enum/source_type"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/model"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/util"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/vo"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/snowflake"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
	"strings"
)

type AlarmRuleDao struct {
}

var AlarmRule = new(AlarmRuleDao)

func (dao *AlarmRuleDao) SaveRule(tx *gorm.DB, ruleReqDTO *form.AlarmRuleAddReqDTO) {
	rule := buildAlarmRule(ruleReqDTO)
	rule.MonitorType = ruleReqDTO.MonitorType
	rule.ProductName = ruleReqDTO.ProductType
	tx.Create(rule)
	dao.saveRuleOthers(tx, ruleReqDTO, rule.BizId)
}
func (dao *AlarmRuleDao) UpdateRule(tx *gorm.DB, ruleReqDTO *form.AlarmRuleAddReqDTO) {
	if !dao.checkRuleExists(tx, ruleReqDTO.Id, ruleReqDTO.TenantId) {
		logger.Logger().Infof("规则不存在 %+v", ruleReqDTO)
		return
	}
	dao.deleteOthers(tx, ruleReqDTO.Id)
	rule := buildAlarmRule(ruleReqDTO)
	tx.Model(&rule).Where("biz_id=?", ruleReqDTO.Id).Updates(rule)
	dao.saveRuleOthers(tx, ruleReqDTO, ruleReqDTO.Id)
}

func (dao *AlarmRuleDao) DeleteRule(tx *gorm.DB, ruleReqDTO *form.RuleReqDTO) {
	if !dao.checkRuleExists(tx, ruleReqDTO.Id, ruleReqDTO.TenantId) {
		logger.Logger().Infof("规则不存在 %+v", ruleReqDTO)
		return
	}
	rule := model.AlarmRule{
		TenantID: ruleReqDTO.TenantId,
		BizId:    ruleReqDTO.Id,
		Deleted:  1,
	}
	tx.Where("biz_id=?", ruleReqDTO.Id).Updates(&rule)
	dao.deleteOthers(tx, ruleReqDTO.Id)
}

func (dao *AlarmRuleDao) UpdateRuleState(tx *gorm.DB, ruleReqDTO *form.RuleReqDTO) {
	if !dao.checkRuleExists(tx, ruleReqDTO.Id, ruleReqDTO.TenantId) {
		logger.Logger().Infof("规则不存在 %+v", ruleReqDTO)
		return
	}
	rule := model.AlarmRule{}
	tx.Model(&rule).Where("biz_id=?", ruleReqDTO.Id).Update("enabled", getAlarmStatusTextInt(ruleReqDTO.Status))
}

func (dao *AlarmRuleDao) checkRuleExists(tx *gorm.DB, ruleId string, tenantId string) bool {
	var count int64
	tx.Model(&model.AlarmRule{}).Where("biz_id=?", ruleId).Where("tenant_id=?", tenantId).Count(&count)
	return count != 0
}

func (dao *AlarmRuleDao) SelectRulePageList(param *form.AlarmPageReqParam) *vo.PageVO {
	var modelList []form.AlarmRulePageDTO
	selectList := &strings.Builder{}
	var sqlParam = []interface{}{param.TenantId}
	selectList.WriteString("select * from ( SELECT    count(DISTINCT(t2.resource_id)) AS instanceNum,    t1.NAME  as name,    t1.monitor_type,    t1.product_name as product_type,    t1.metric_name,    t1.trigger_condition,    t1.enabled AS 'status',    t1.biz_id AS ruleId,    t1.create_time  FROM    t_alarm_rule t1  LEFT JOIN t_alarm_rule_resource_rel t2 ON t2.alarm_rule_id = t1.biz_id  WHERE    t1.tenant_id = ?  AND t1.deleted = 0  AND t1.source_type = 1  ")
	if len(param.Status) != 0 {
		selectList.WriteString(" and t1.enabled = ?")
		sqlParam = append(sqlParam, getAlarmStatusTextInt(param.Status))
	}
	if len(param.RuleName) != 0 {
		selectList.WriteString(" and t1.name like concat('%',?,'%')")
		sqlParam = append(sqlParam, param.RuleName)
	}
	selectList.WriteString(" group by ruleId ) t order by t.create_time  desc ")
	page := util.Paginate(param.PageSize, param.Current, selectList.String(), sqlParam, &modelList)
	for i, v := range modelList {
		modelList[i].MonitorItem = v.RuleCondition.MonitorItemName
		modelList[i].Express = GetExpress(v.RuleCondition)
		modelList[i].Status = getAlarmStatusSqlText(v.Status)
	}
	return &vo.PageVO{
		Records: modelList,
		Current: page.Current,
		Size:    page.Size,
		Total:   page.Total,
		Pages:   page.Pages,
	}

}

func (dao *AlarmRuleDao) GetDetail(tx *gorm.DB, id string, tenantId string) *form.AlarmRuleDetailDTO {
	ruleDto := &form.AlarmRuleDetailDTO{}
	tx.Raw("SELECT    biz_id as id ,    NAME  as ruleName,  enabled as status,  product_name as product_type,  monitor_type,   level as alarmLevel,  dimensions as scope,  trigger_condition as ruleCondition ,  silences_time,   effective_start,  effective_end    FROM t_alarm_rule        WHERE biz_id = ?          AND deleted = 0  and tenant_id=?", id, tenantId).Scan(ruleDto)
	ruleDto.NoticeGroups = dao.GetNoticeGroupList(tx, id)
	ruleDto.InstanceList = dao.GetInstanceList(tx, id)
	ruleDto.AlarmHandlerList = dao.GetHandlerList(tx, id)
	ruleDto.Status = getAlarmStatusSqlText(ruleDto.Status)
	ruleDto.Describe = GetExpress(ruleDto.RuleCondition)
	scope, _ := strconv.Atoi(ruleDto.Scope)
	ruleDto.Scope = getResourceScopeText(scope)
	return ruleDto
}

func (dao *AlarmRuleDao) GetById(db *gorm.DB, id string) *model.AlarmRule {
	var alarmRule = &model.AlarmRule{}
	db.Where("biz_id=?", id).Find(&alarmRule)
	return alarmRule
}

func (dao *AlarmRuleDao) GetInstanceList(tx *gorm.DB, ruleId string) []*form.InstanceInfo {
	var instanceInfo []*form.InstanceInfo
	tx.Raw("select * from (SELECT   DISTINCT   t2.instance_id,   t2.region_code,   t2.region_name,   t2.instance_name    FROM   t_alarm_rule_resource_rel t1,   t_resource t2    WHERE   t1.alarm_rule_id = ?    AND t1.resource_id = t2.instance_id and t2.region_code!='' ) t group by   instance_id", ruleId).Scan(&instanceInfo)
	return instanceInfo
}

func (dao *AlarmRuleDao) GetNoticeGroupList(tx *gorm.DB, ruleId string) []*form.NoticeGroup {
	var noticeGroup []*form.NoticeGroup
	tx.Raw("SELECT t1.contract_group_id as id, t2.`name` as name  FROM t_alarm_notice t1,  t_contact_group t2   WHERE t1.alarm_rule_id = ?   and t1.contract_group_id = t2.biz_id  ORDER BY name", ruleId).Scan(&noticeGroup)
	for _, group := range noticeGroup {
		group.UserList = dao.GetUserList(tx, group.Id)
	}
	return noticeGroup
}

func (dao *AlarmRuleDao) GetUserList(tx *gorm.DB, groupId string) []*form.UserInfo {
	var userInfo []*form.UserInfo
	tx.Raw("select t2.`name` as userName  ,t2.biz_id as id, GROUP_CONCAT(CASE t3.type WHEN 1 THEN t3.address  END) as phone, GROUP_CONCAT(CASE t3.type WHEN 2 THEN t3.address  END) as email from t_contact_group_rel  t   LEFT JOIN t_contact t2 on t2.biz_id = t.contact_biz_id   LEFT JOIN t_contact_information t3 on (t3.contact_biz_id = t2.biz_id and t3.state=1)  where t.group_biz_id=?  and t2.state =1  GROUP BY id  order by userName", groupId).Scan(&userInfo)
	return userInfo
}

func (dao *AlarmRuleDao) GetMonitorItem(metricName string) *model.MonitorItem {
	item := &model.MonitorItem{}
	global.DB.Raw("SELECT metrics_linux,metrics_windows,metric_name,unit,name  ,labels FROM t_monitor_item  where metric_name=? ", metricName).Scan(item)
	return item
}

func buildAlarmRule(ruleReqDTO *form.AlarmRuleAddReqDTO) *model.AlarmRule {
	return &model.AlarmRule{TenantID: ruleReqDTO.TenantId,
		BizId:         ruleReqDTO.Id,
		ProductName:   ruleReqDTO.ProductType,
		Dimensions:    GetResourceScopeInt(ruleReqDTO.Scope),
		Name:          ruleReqDTO.RuleName,
		MetricName:    ruleReqDTO.RuleCondition.MetricName,
		RuleCondition: ruleReqDTO.RuleCondition,
		SilencesTime:  ruleReqDTO.SilencesTime,
		Level:         ruleReqDTO.AlarmLevel,
		CreateUser:    ruleReqDTO.UserId,
		Source:        ruleReqDTO.Source,
		SourceType:    ruleReqDTO.SourceType,
		ProductBizId:  strconv.Itoa(ruleReqDTO.ProductId),
	}
}

func (dao *AlarmRuleDao) saveRuleOthers(tx *gorm.DB, ruleReqDTO *form.AlarmRuleAddReqDTO, ruleId string) {
	dao.saveAlarmNotice(tx, ruleReqDTO, ruleId)
	dao.saveAlarmRuleResGroup(tx, ruleReqDTO)
	dao.saveAlarmRuleResource(tx, ruleReqDTO)
	dao.saveAlarmHandler(tx, ruleReqDTO)
}

func (dao *AlarmRuleDao) deleteOthers(tx *gorm.DB, ruleId string) {
	rule := &model.AlarmRule{}
	tx.Where("biz_id=?", ruleId).Find(rule)
	if rule.SourceType == source_type.AutoScaling {
		tx.Exec("delete t2.* from t_alarm_rule_group_rel  t1 INNER JOIN  t_resource_group t2  where  t1.alarm_rule_id=? and t1.resource_group_id=t2.id", ruleId)
		//todo 弹性伸缩删除资源组关联的资源
	}
	//删除规则关联的联系组
	tx.Where("alarm_rule_id=?", ruleId).Delete(&model.AlarmNotice{})
	//删除规则关联与资源组的关系
	tx.Where("alarm_rule_id=?", ruleId).Delete(&model.AlarmRuleGroupRel{})
	//删除规则关联与资源的关系
	tx.Where("alarm_rule_id=?", ruleId).Delete(&model.AlarmRuleResourceRel{})
	//删除规则关联的告警处理
	tx.Where("alarm_rule_id=?", ruleId).Delete(&model.AlarmHandler{})
}

// saveAlarmNotice 保存规则的告警联系组
func (dao *AlarmRuleDao) saveAlarmNotice(tx *gorm.DB, ruleReqDTO *form.AlarmRuleAddReqDTO, ruleId string) {
	if len(ruleReqDTO.GroupList) == 0 {
		return
	}
	list := make([]model.AlarmNotice, len(ruleReqDTO.GroupList))
	for index, group := range ruleReqDTO.GroupList {
		list[index] = model.AlarmNotice{
			AlarmRuleID:     ruleId,
			ContractGroupID: group,
		}
	}
	tx.Create(&list)
}

// saveAlarmRuleResGroup 保存资源组、资源组与规则的关系
func (dao *AlarmRuleDao) saveAlarmRuleResGroup(tx *gorm.DB, ruleReqDTO *form.AlarmRuleAddReqDTO) {
	groupSize := len(ruleReqDTO.ResourceGroupList)
	if groupSize == 0 {
		return
	}
	groupRelList := make([]*model.AlarmRuleGroupRel, groupSize)
	groups := make([]*model.ResourceGroup, groupSize)
	for index, info := range ruleReqDTO.ResourceGroupList {
		if len(info.ResGroupId) == 0 {
			info.ResGroupId = strconv.FormatInt(snowflake.GetWorker().NextId(), 10)
		}
		groups[index] = &model.ResourceGroup{
			BizId:      info.ResGroupId,
			Name:       info.ResGroupName,
			TenantId:   ruleReqDTO.TenantId,
			ProductId:  ruleReqDTO.ProductId,
			SourceType: ruleReqDTO.SourceType,
		}
		groupRelList[index] = &model.AlarmRuleGroupRel{
			AlarmRuleId:     ruleReqDTO.Id,
			ResourceGroupId: info.ResGroupId,
			CalcMode:        info.CalcMode,
			TenantId:        ruleReqDTO.TenantId,
		}
		dao.saveResource(tx, ruleReqDTO.TenantId, info, ruleReqDTO.ProductType)
		tx.Create(&groups)
		tx.Create(&groupRelList)
	}
}

// saveResource 保存资源 、资源与组的关系
func (dao *AlarmRuleDao) saveResource(tx *gorm.DB, tenantID string, info *form.ResGroupInfo, productType string) {
	resSize := len(info.ResourceList)
	if resSize == 0 {
		return
	}
	var resourceList []*model.AlarmInstance
	var resourceRelList []*model.ResourceResourceGroupRel
	resourceMap := map[string]byte{}
	for _, resInfo := range info.ResourceList {
		_, ok := resourceMap[resInfo.InstanceId]
		if ok {
			continue
		}
		resourceMap[resInfo.InstanceId] = 0
		resource := dao.buildResource(resInfo, tenantID, productType)
		resourceList = append(resourceList, resource)
		resourceRel := &model.ResourceResourceGroupRel{
			ResourceGroupId: info.ResGroupId,
			ResourceId:      resInfo.InstanceId,
			TenantId:        tenantID,
		}
		resourceRelList = append(resourceRelList, resourceRel)
	}
	tx.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "instance_id"}}, DoNothing: false}).Create(&resourceList)
	tx.Create(&resourceRelList)
}

// SaveResource 保存资源 、资源与规则的关系
func (dao *AlarmRuleDao) saveAlarmRuleResource(tx *gorm.DB, ruleReqDTO *form.AlarmRuleAddReqDTO) {
	resourceSize := len(ruleReqDTO.ResourceList)
	if resourceSize == 0 {
		return
	}
	resourceRelList := make([]*model.AlarmRuleResourceRel, resourceSize)
	resourceList := make([]*model.AlarmInstance, resourceSize)
	resourceMap := map[string]byte{}
	for _, info := range ruleReqDTO.ResourceList {
		_, ok := resourceMap[info.InstanceId]
		if ok {
			continue
		}
		resourceMap[info.InstanceId] = 0
		resource := dao.buildResource(info, ruleReqDTO.TenantId, ruleReqDTO.ProductType)
		resourceRel := &model.AlarmRuleResourceRel{
			AlarmRuleId: ruleReqDTO.Id,
			ResourceId:  info.InstanceId,
			TenantId:    ruleReqDTO.TenantId,
		}
		resourceRelList = append(resourceRelList, resourceRel)
		resourceList = append(resourceList, resource)
	}
	tx.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "instance_id"}}, DoNothing: false}).Create(&resourceList)
	tx.Create(&resourceRelList)
}

// saveAlarmHandler 保存规则告警handler
func (dao *AlarmRuleDao) saveAlarmHandler(tx *gorm.DB, ruleReqDTO *form.AlarmRuleAddReqDTO) {
	handlerSize := len(ruleReqDTO.AlarmHandlerList)
	if handlerSize == 0 {
		return
	}
	handlers := make([]*model.AlarmHandler, handlerSize)
	for index, info := range ruleReqDTO.AlarmHandlerList {
		handlers[index] = &model.AlarmHandler{
			AlarmRuleId:  ruleReqDTO.Id,
			HandleType:   info.HandleType,
			HandleParams: info.HandleParams,
			TenantId:     ruleReqDTO.TenantId,
		}
	}
	tx.Create(&handlers)
}

func (dao *AlarmRuleDao) buildResource(info *form.InstanceInfo, tenantId string, productType string) *model.AlarmInstance {
	return &model.AlarmInstance{
		Ip:           info.Ip,
		InstanceID:   info.InstanceId,
		RegionCode:   info.RegionCode,
		ZoneCode:     info.ZoneCode,
		ZoneName:     info.ZoneName,
		RegionName:   info.RegionName,
		InstanceName: info.InstanceName,
		TenantID:     tenantId,
		ProductName:  productType,
		ProductBizId: "",
	}
}

func (dao *AlarmRuleDao) GetHandlerList(tx *gorm.DB, ruleId string) []*form.Handler {
	var handler []*form.Handler
	tx.Raw("SELECT handle_type,handle_params FROM t_alarm_handler where alarm_rule_id=?", ruleId).Scan(&handler)
	return handler
}

func (dao *AlarmRuleDao) GetResourceListByGroup(tx *gorm.DB, resGroup string) []*form.InstanceInfo {
	var instanceInfo []*form.InstanceInfo
	tx.Raw(" SELECT  DISTINCT t2.instance_id  FROM   t_resource t2,   t_resource_resource_group_rel t1 WHERE   t1.resource_group_id = ? AND t1.resource_id = t2.instance_id", resGroup).Scan(&instanceInfo)
	return instanceInfo
}

const (
	ALL      = "ALL"
	INSTANCE = "INSTANCE"
)

var ResourceScopeText = map[string]uint8{
	ALL:      1,
	INSTANCE: 2,
}

func GetResourceScopeInt(code string) uint8 {
	return ResourceScopeText[code]
}

const (
	ENABLE  = "enabled"
	DISABLE = "disabled"

	Maximum = "Maximum"
	Minimum = "Minimum"
	Average = "Average"

	Greater        = "greater"
	GreaterOrEqual = "greaterOrEqual"
	Less           = "less"
	lessOrEqual    = "lessOrEqual"
	Equal          = "equal"
	NotEqual       = "notEqual"
)

var alarmStatusText = map[string]int{
	ENABLE:  1,
	DISABLE: 0,
}

func getAlarmStatusTextInt(code string) int {
	return alarmStatusText[code]
}

var alarmStatusSqlText = map[string]string{
	"1": "enabled",
	"0": "disabled",
}

func getAlarmStatusSqlText(code string) string {
	return alarmStatusSqlText[code]
}

var alarmStatisticsText = map[string]string{
	Maximum: "最大值",
	Minimum: "最小值",
	Average: "平均值",
}

func GetAlarmStatisticsText(s string) string {
	return alarmStatisticsText[s]
}

var comparisonOperatorText = map[string]string{
	Greater:        ">",
	GreaterOrEqual: ">=",
	Less:           "<",
	lessOrEqual:    "<=",
	Equal:          "==",
	NotEqual:       "!=",
}

func GetComparisonOperator(s string) string {
	return comparisonOperatorText[s]
}

func GetExpress(form *form.RuleCondition) string {
	return fmt.Sprintf("%s%s%s%s%s 统计周期%s分钟 持续%s个周期", form.MonitorItemName, GetAlarmStatisticsText(form.Statistics), GetComparisonOperator(form.ComparisonOperator), fmt.Sprintf("%.f", form.Threshold), form.Unit, strconv.Itoa(form.Period/60), strconv.Itoa(form.Times))
}

var ResourceScopeInt = map[int]string{
	1: ALL,
	2: INSTANCE,
}

func getResourceScopeText(code int) string {
	return ResourceScopeInt[code]
}
