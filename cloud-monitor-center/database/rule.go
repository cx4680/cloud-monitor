package database

var SelectRuleDetail = "SELECT               id ,               NAME  as ruleName,      " +
	"         enabled as status,               product_type,           " +
	"    monitor_type,               level as alarmLevel,         " +
	"      dimensions as scope,               trigger_condition as ruleCondition ,      " +
	"         silences_time,               effective_start,               effective_end,    " +
	"           notify_channel as noticeChannel        FROM t_alarm_rule        WHERE id = ?     " +
	"     AND deleted = 0  and tenant_id=?"
