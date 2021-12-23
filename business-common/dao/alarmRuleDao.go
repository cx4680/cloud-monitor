package dao

import (
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/enums/sourceType"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/forms"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/global"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/pageUtils"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/vo"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/logger"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/utils/snowflake"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
	"strings"
)

type AlarmRuleDao struct {
}

var AlarmRule = new(AlarmRuleDao)

func (dao *AlarmRuleDao) SaveRule(tx *gorm.DB, ruleReqDTO *forms.AlarmRuleAddReqDTO) {
	rule := buildAlarmRule(ruleReqDTO)
	rule.MonitorType = ruleReqDTO.MonitorType
	rule.ProductType = ruleReqDTO.ProductType
	tx.Create(rule)
	dao.saveRuleOthers(tx, ruleReqDTO, rule.ID)
}
func (dao *AlarmRuleDao) UpdateRule(tx *gorm.DB, ruleReqDTO *forms.AlarmRuleAddReqDTO) {
	if !dao.checkRuleExists(tx, ruleReqDTO.Id, ruleReqDTO.TenantId) {
		logger.Logger().Infof("规则不存在 %+v", ruleReqDTO)
		return
	}
	dao.deleteOthers(tx, ruleReqDTO.Id)
	rule := buildAlarmRule(ruleReqDTO)
	tx.Model(&rule).Updates(rule)
	dao.saveRuleOthers(tx, ruleReqDTO, ruleReqDTO.Id)
}

func (dao *AlarmRuleDao) DeleteRule(tx *gorm.DB, ruleReqDTO *forms.RuleReqDTO) {
	if !dao.checkRuleExists(tx, ruleReqDTO.Id, ruleReqDTO.TenantId) {
		logger.Logger().Infof("规则不存在 %+v", ruleReqDTO)
		return
	}
	rule := models.AlarmRule{
		TenantID: ruleReqDTO.TenantId,
		ID:       ruleReqDTO.Id,
	}
	tx.Delete(&rule)
	dao.deleteOthers(tx, ruleReqDTO.Id)
}

func (dao *AlarmRuleDao) UpdateRuleState(tx *gorm.DB, ruleReqDTO *forms.RuleReqDTO) {
	if !dao.checkRuleExists(tx, ruleReqDTO.Id, ruleReqDTO.TenantId) {
		logger.Logger().Infof("规则不存在 %+v", ruleReqDTO)
		return
	}
	rule := models.AlarmRule{ID: ruleReqDTO.Id}
	tx.Model(&rule).Update("enabled", getAlarmStatusTextInt(ruleReqDTO.Status))
}

func (dao *AlarmRuleDao) checkRuleExists(tx *gorm.DB, ruleId string, tenantId string) bool {
	var count int64
	tx.Model(&models.AlarmRule{}).Where("id=?", ruleId).Where("tenant_id=?", tenantId).Count(&count)
	if count == 0 {
		return false
	}
	return true
}

func (dao *AlarmRuleDao) SelectRulePageList(param *forms.AlarmPageReqParam) *vo.PageVO {
	var modelList []forms.AlarmRulePageDTO
	selectList := &strings.Builder{}
	var sqlParam = []interface{}{param.TenantId}
	selectList.WriteString("select * from ( SELECT    count(DISTINCT(t2.resource_id)) AS instanceNum,    t1.NAME  as name,    t1.monitor_type,    t1.product_type,    t1.metric_name,    t1.trigger_condition,    t1.enabled AS 'status',    t1.id AS ruleId,    t1.create_time  FROM    t_alarm_rule t1  LEFT JOIN t_alarm_rule_resource_rel t2 ON t2.alarm_rule_id = t1.id  WHERE    t1.tenant_id = ?  AND t1.deleted = 0  AND t1.source_type = 1  ")
	if len(param.Status) != 0 {
		selectList.WriteString(" and t1.enabled = ?")
		sqlParam = append(sqlParam, param.Status)
	}
	if len(param.RuleName) != 0 {
		selectList.WriteString(" and t1.name like concat('%',?,'%')")
		sqlParam = append(sqlParam, param.RuleName)
	}
	selectList.WriteString(" group by ruleId ) t order by t.create_time  desc ")
	page := pageUtils.Paginate(param.PageSize, param.Current, selectList.String(), sqlParam, &modelList)
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

func (dao *AlarmRuleDao) GetDetail(tx *gorm.DB, id string, tenantId string) *forms.AlarmRuleDetailDTO {
	model := &forms.AlarmRuleDetailDTO{}
	tx.Raw("SELECT    id ,    NAME  as ruleName,  enabled as status,  product_type,  monitor_type,   level as alarmLevel,  dimensions as scope,  trigger_condition as ruleCondition ,  silences_time,   effective_start,  effective_end    FROM t_alarm_rule        WHERE id = ?          AND deleted = 0  and tenant_id=?", id, tenantId).Scan(model)
	model.NoticeGroups = dao.GetNoticeGroupList(tx, id)
	model.InstanceList = dao.GetInstanceList(tx, id)
	model.AlarmHandlerList = dao.GetHandlerList(tx, id)
	model.Status = getAlarmStatusSqlText(model.Status)
	model.Describe = GetExpress(model.RuleCondition)
	scope, _ := strconv.Atoi(model.Scope)
	model.Scope = getResourceScopeText(scope)
	return model
}

func (dao *AlarmRuleDao) GetInstanceList(tx *gorm.DB, ruleId string) []*forms.InstanceInfo {
	var model []*forms.InstanceInfo
	tx.Raw("SELECT DISTINCT t2.instance_id, t2.region_code, t2.region_name, t2.ip, t2.instance_name  FROM t_alarm_rule_resource_rel t1, t_alarm_instance t2  WHERE t1.alarm_rule_id = ?  AND t1.resource_id = t2.instance_id", ruleId).Scan(&model)
	return model
}

func (dao *AlarmRuleDao) GetNoticeGroupList(tx *gorm.DB, ruleId string) []*forms.NoticeGroup {
	var model []*forms.NoticeGroup
	tx.Raw("SELECT t1.contract_group_id as id, t2.`name` as name  FROM t_alarm_notice t1,  alert_contact_group t2   WHERE t1.alarm_rule_id = ?   and t1.contract_group_id = t2.id  ORDER BY name", ruleId).Scan(&model)
	for _, group := range model {
		group.UserList = dao.GetUserList(tx, group.Id)
	}
	return model
}

func (dao *AlarmRuleDao) GetUserList(tx *gorm.DB, groupId string) []*forms.UserInfo {
	var model []*forms.UserInfo
	tx.Raw("select t2.`name` as userName  ,t2.id as id, GROUP_CONCAT(CASE t3.type WHEN 1 THEN t3.no  END) as phone, GROUP_CONCAT(CASE t3.type WHEN 2 THEN t3.no  END) as email from alert_contact_group_rel  t   LEFT JOIN alert_contact t2 on t2.id = t.contact_id   LEFT JOIN alert_contact_information t3 on (t3.contact_id = t2.id and t3.is_certify=1)  where t.group_id=?  and t2.`status`=1  GROUP BY id  order by userName", groupId).Scan(&model)
	return model
}

func (dao *AlarmRuleDao) GetMonitorItem(metricName string) *models.MonitorItem {
	model := &models.MonitorItem{}
	global.DB.Raw("SELECT metrics_linux,metrics_windows,metric_name,unit,name  ,labels FROM monitor_item  where metric_name=? ", metricName).Scan(model)
	return model
}

func buildAlarmRule(ruleReqDTO *forms.AlarmRuleAddReqDTO) *models.AlarmRule {
	return &models.AlarmRule{TenantID: ruleReqDTO.TenantId,
		ID:            ruleReqDTO.Id,
		ProductType:   ruleReqDTO.ProductType,
		Dimensions:    GetResourceScopeInt(ruleReqDTO.Scope),
		Name:          ruleReqDTO.RuleName,
		MetricName:    ruleReqDTO.RuleCondition.MetricName,
		RuleCondition: ruleReqDTO.RuleCondition,
		SilencesTime:  ruleReqDTO.SilencesTime,
		Level:         ruleReqDTO.AlarmLevel,
		CreateUser:    ruleReqDTO.UserId,
		Source:        ruleReqDTO.Source,
		SourceType:    ruleReqDTO.SourceType,
	}
}

func (dao *AlarmRuleDao) saveRuleOthers(tx *gorm.DB, ruleReqDTO *forms.AlarmRuleAddReqDTO, ruleId string) {
	dao.saveAlarmNotice(tx, ruleReqDTO, ruleId)
	dao.saveAlarmRuleResGroup(tx, ruleReqDTO)
	dao.saveAlarmRuleResource(tx, ruleReqDTO)
	dao.saveAlarmHandler(tx, ruleReqDTO)
}

func (dao *AlarmRuleDao) deleteOthers(tx *gorm.DB, ruleId string) {
	rule := &models.AlarmRule{}
	tx.Where("id=?", ruleId).Find(rule)
	if rule.SourceType == sourceType.AutoScaling {
		tx.Exec("delete t2.* from t_alarm_rule_group_rel  t1 INNER JOIN  t_resource_group t2  where  t1.alarm_rule_id=? and t1.resource_group_id=t2.id", ruleId)
		//todo 弹性伸缩删除资源组关联的资源
	}
	//删除规则关联的联系组
	tx.Where("alarm_rule_id=?", ruleId).Delete(&models.AlarmNotice{})
	//删除规则关联与资源组的关系
	tx.Where("alarm_rule_id=?", ruleId).Delete(&models.AlarmRuleGroupRel{})
	//删除规则关联与资源的关系
	tx.Where("alarm_rule_id=?", ruleId).Delete(&models.AlarmRuleResourceRel{})
	//删除规则关联的告警处理
	tx.Where("alarm_rule_id=?", ruleId).Delete(&models.AlarmHandler{})
}

// saveAlarmNotice 保存规则的告警联系组
func (dao *AlarmRuleDao) saveAlarmNotice(tx *gorm.DB, ruleReqDTO *forms.AlarmRuleAddReqDTO, ruleId string) {
	if len(ruleReqDTO.GroupList) == 0 {
		return
	}
	list := make([]models.AlarmNotice, len(ruleReqDTO.GroupList))
	for index, group := range ruleReqDTO.GroupList {
		list[index] = models.AlarmNotice{
			AlarmRuleID:     ruleId,
			ContractGroupID: group,
		}
	}
	tx.Create(&list)
}

// saveAlarmRuleResGroup 保存资源组、资源组与规则的关系
func (dao *AlarmRuleDao) saveAlarmRuleResGroup(tx *gorm.DB, ruleReqDTO *forms.AlarmRuleAddReqDTO) {
	groupSize := len(ruleReqDTO.ResourceGroupList)
	if groupSize == 0 {
		return
	}
	groupRelList := make([]*models.AlarmRuleGroupRel, groupSize)
	groups := make([]*models.ResourceGroup, groupSize)
	for index, info := range ruleReqDTO.ResourceGroupList {
		if len(info.ResGroupId) == 0 {
			info.ResGroupId = strconv.FormatInt(snowflake.GetWorker().NextId(), 10)
		}
		groups[index] = &models.ResourceGroup{
			Id:         info.ResGroupId,
			Name:       info.ResGroupName,
			TenantId:   ruleReqDTO.TenantId,
			ProductId:  ruleReqDTO.ProductId,
			SourceType: ruleReqDTO.SourceType,
		}
		groupRelList[index] = &models.AlarmRuleGroupRel{
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
func (dao *AlarmRuleDao) saveResource(tx *gorm.DB, tenantID string, info *forms.ResGroupInfo, productType string) {
	resSize := len(info.ResourceList)
	if resSize == 0 {
		return
	}
	resourceList := make([]*models.AlarmInstance, resSize)
	resourceRelList := make([]*models.ResourceResourceGroupRel, resSize)
	for index, resInfo := range info.ResourceList {
		resourceList[index] = dao.buildResource(resInfo, tenantID, productType)
		resourceRelList[index] = &models.ResourceResourceGroupRel{
			ResourceGroupId: info.ResGroupId,
			ResourceId:      resInfo.InstanceId,
			TenantId:        tenantID,
		}
	}
	tx.Clauses(clause.OnConflict{DoNothing: false}).Create(&resourceList)
	tx.Create(&resourceRelList)
}

// SaveResource 保存资源 、资源与规则的关系
func (dao *AlarmRuleDao) saveAlarmRuleResource(tx *gorm.DB, ruleReqDTO *forms.AlarmRuleAddReqDTO) {
	resourceSize := len(ruleReqDTO.ResourceList)
	if resourceSize == 0 {
		return
	}
	resourceRelList := make([]*models.AlarmRuleResourceRel, resourceSize)
	resourceList := make([]*models.AlarmInstance, resourceSize)
	for index, info := range ruleReqDTO.ResourceList {
		resourceRelList[index] = &models.AlarmRuleResourceRel{
			AlarmRuleId: ruleReqDTO.Id,
			ResourceId:  info.InstanceId,
			TenantId:    ruleReqDTO.TenantId,
		}
		resourceList[index] = dao.buildResource(info, ruleReqDTO.TenantId, ruleReqDTO.ProductType)
	}
	tx.Clauses(clause.OnConflict{DoNothing: false}).Create(&resourceList)
	tx.Create(&resourceRelList)
}

// saveAlarmHandler 保存规则告警handler
func (dao *AlarmRuleDao) saveAlarmHandler(tx *gorm.DB, ruleReqDTO *forms.AlarmRuleAddReqDTO) {
	handlerSize := len(ruleReqDTO.AlarmHandlerList)
	if handlerSize == 0 {
		return
	}
	handlers := make([]*models.AlarmHandler, handlerSize)
	for index, info := range ruleReqDTO.AlarmHandlerList {
		handlers[index] = &models.AlarmHandler{
			Id:           strconv.FormatInt(snowflake.GetWorker().NextId(), 10),
			AlarmRuleId:  ruleReqDTO.Id,
			HandleType:   info.HandleType,
			HandleParams: info.HandleParams,
			TenantId:     ruleReqDTO.TenantId,
		}
	}
	tx.Create(&handlers)
}

func (dao *AlarmRuleDao) buildResource(info *forms.InstanceInfo, tenantId string, productType string) *models.AlarmInstance {
	return &models.AlarmInstance{
		Ip:           info.Ip,
		InstanceID:   info.InstanceId,
		RegionCode:   info.RegionCode,
		ZoneCode:     info.ZoneCode,
		ZoneName:     info.ZoneName,
		RegionName:   info.RegionName,
		InstanceName: info.InstanceName,
		TenantID:     tenantId,
		ProductType:  productType,
	}
}

func (dao *AlarmRuleDao) GetHandlerList(tx *gorm.DB, ruleId string) []*forms.Handler {
	var model []*forms.Handler
	tx.Raw("SELECT handle_type,handle_params FROM t_alarm_handler where alarm_rule_id=?", ruleId).Scan(&model)
	return model
}

func (dao *AlarmRuleDao) GetResourceListByGroup(tx *gorm.DB, resGroup string) []*forms.InstanceInfo {
	var model []*forms.InstanceInfo
	tx.Raw(" SELECT  DISTINCT t2.instance_id  FROM   t_alarm_instance t2,   t_resource_resource_group_rel t1 WHERE   t1.resource_group_id = ? AND t1.resource_id = t2.instance_id", resGroup).Scan(&model)
	return model
}

const (
	ALL      = "ALL"
	INSTANCE = "INSTANCE"
)

var ResourceScopeText = map[string]int{
	ALL:      1,
	INSTANCE: 2,
}

func GetResourceScopeInt(code string) int {
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

func getAlarmStatisticsText(s string) string {
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

func getComparisonOperator(s string) string {
	return comparisonOperatorText[s]
}

func GetExpress(form *forms.RuleCondition) string {
	return fmt.Sprintf("%s%s%s%s%s 统计周期%s分钟 持续%s个周期", form.MonitorItemName, getAlarmStatisticsText(form.Statistics), getComparisonOperator(form.ComparisonOperator), strconv.FormatFloat(form.Threshold, 'g', 5, 32), form.Unit, strconv.Itoa(form.Period/60), strconv.Itoa(form.Times))
}

var ResourceScopeInt = map[int]string{
	1: ALL,
	2: INSTANCE,
}

func getResourceScopeText(code int) string {
	return ResourceScopeInt[code]
}
