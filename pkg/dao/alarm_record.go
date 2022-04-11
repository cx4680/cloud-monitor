package dao

import (
	"bytes"
	"code.cestc.cn/ccos-ops/cloud-monitor/business-common/enum/source_type"
	commonModels "code.cestc.cn/ccos-ops/cloud-monitor/business-common/model"
	cvo "code.cestc.cn/ccos-ops/cloud-monitor/business-common/vo"
	"code.cestc.cn/ccos-ops/cloud-monitor/common/util/strutil"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/form"
	"code.cestc.cn/ccos-ops/cloud-monitor/pkg/vo"
	"gorm.io/gorm"
	"strings"
)

type AlarmRecordDao struct {
}

var AlarmRecord = new(AlarmRecordDao)

const (
	alarmPageBaseSql = "SELECT ar.id, ar.biz_id, ar.`status`, ar.tenant_id, ar.rule_id, ar.rule_name, ar.monitor_type,ar.source_type, " +
		" ar.source_id, ai.summary,ar.current_value,ar.start_time,ar.end_time,ar.target_value,ai.expression,ar.duration,ar.`level`, " +
		" ar.alarm_key, ar.region,ar.create_time,ar.update_time  " +
		" FROM t_alarm_record ar join t_alarm_info ai on ar.biz_id=ai.alarm_biz_id " +
		" WHERE ar.tenant_id=? and ar.rule_source_type!=? "

	alarmDetailBaseSql = "SELECT ar.id, ar.biz_id, ar.`status`, ar.tenant_id, ar.rule_id, ar.rule_name, ar.monitor_type,ar.source_type, " +
		" ar.source_id,ar.current_value,ar.start_time,ar.end_time,ar.target_value,ar.duration,ar.`level`, " +
		" ar.alarm_key, ar.region,ar.create_time,ar.update_time, ai.contact_info, ai.summary,ai.expression  " +
		" FROM t_alarm_record ar join t_alarm_info ai on ar.biz_id=ai.alarm_biz_id " +
		" WHERE ar.biz_id=? and ar.tenant_id=? "
)

func (a *AlarmRecordDao) GetPageList(db *gorm.DB, tenantId string, f form.AlarmRecordPageQueryForm) *cvo.PageVO {
	var list []vo.AlarmRecordPageVO
	var c = bytes.Buffer{}
	var cv []interface{}
	c.WriteString(alarmPageBaseSql)
	cv = append(cv, tenantId, source_type.AutoScaling)

	if strutil.IsNotBlank(f.Level) {
		c.WriteString(" and ar.level in (?) ")
		cv = append(cv, strings.Split(f.Level, ","))
	}
	if strutil.IsNotBlank(f.Region) {
		c.WriteString(" and ar.region=? ")
		cv = append(cv, f.Region)
	}
	if strutil.IsNotBlank(f.ResourceId) {
		c.WriteString(" and ar.source_id=? ")
		cv = append(cv, f.ResourceId)
	}
	if strutil.IsNotBlank(f.ResourceType) {
		c.WriteString(" and ar.source_type=? ")
		cv = append(cv, f.ResourceType)
	}
	if strutil.IsNotBlank(f.RuleId) {
		c.WriteString(" and ar.rule_id=? ")
		cv = append(cv, f.RuleId)
	}
	if strutil.IsNotBlank(f.RuleName) {
		c.WriteString(" and ar.rule_name like concat('%', ?, '%') ")
		cv = append(cv, f.RuleName)
	}
	if strutil.IsNotBlank(f.Status) {
		c.WriteString(" and ar.status=? ")
		cv = append(cv, f.Status)
	}
	if strutil.IsNotBlank(f.StartTime) {
		c.WriteString(" and ar.create_time>=? ")
		cv = append(cv, f.StartTime)
	}
	if strutil.IsNotBlank(f.EndTime) {
		c.WriteString(" and ar.create_time<=? ")
		cv = append(cv, f.EndTime)
	}
	var total int64
	db.Raw("select count(1) from ("+c.String()+" ) t", cv...).Scan(&total)
	if total > 0 && int(total) >= (f.PageNum-1)*f.PageSize {
		c.WriteString(" order by ar.create_time desc limit ?,?")
		cv = append(cv, (f.PageNum-1)*f.PageSize)
		cv = append(cv, f.PageSize)
		db.Raw(c.String(), cv...).Find(&list)
	}

	return &cvo.PageVO{
		Records: list,
		Total:   int(total),
		Size:    f.PageSize,
		Current: f.PageNum,
		Pages:   (int(total) / f.PageSize) + 1,
	}
}

func (a *AlarmRecordDao) GetByBizIdAndTenantId(db *gorm.DB, bizId, tenantId string) *vo.AlarmRecordDetailVO {
	var detail vo.AlarmRecordDetailVO
	var c = bytes.Buffer{}
	var cv []interface{}
	c.WriteString(alarmDetailBaseSql)
	cv = append(cv, bizId, tenantId)
	db.Raw(c.String(), cv...).Find(&detail)
	return &detail
}

func (a *AlarmRecordDao) GetAlarmRecordTotal(db *gorm.DB, tenantId string, region string, startTime string, endTime string) int64 {
	var count int64
	if region != "" {
		db.Model(&commonModels.AlarmRecord{}).
			Where("tenant_id = ? AND create_time BETWEEN ? AND ? AND region = ? and rule_source_type != ? AND status = ?",
				tenantId, startTime, endTime, region, source_type.AutoScaling, "firing").
			Count(&count)
	} else {
		db.Model(&commonModels.AlarmRecord{}).
			Where("tenant_id = ? AND create_time BETWEEN ? AND ? and rule_source_type != ? AND status = ?",
				tenantId, startTime, endTime, source_type.AutoScaling, "firing").
			Count(&count)
	}
	return count
}

func (a *AlarmRecordDao) GetRecordNumHistory(db *gorm.DB, tenantId string, region string, startTime string, endTime string) []vo.RecordNumHistory {
	var sql = "SELECT COUNT(t.id) AS number, " +
		"DATE_FORMAT(t.create_time, '%Y-%m-%d') AS DayTime " +
		"FROM t_alarm_record t " +
		"WHERE t.tenant_id=? " +
		" AND t.create_time between ? AND ? " +
		" and t.status='firing' "
	var cv []interface{}
	cv = append(cv, tenantId, startTime, endTime)
	if strutil.IsNotBlank(region) {
		sql += " AND t.region=? "
		cv = append(cv, region)
	}
	sql += " GROUP BY daytime "
	var list []vo.RecordNumHistory
	db.Raw(sql, cv...).Find(&list)
	return list
}
