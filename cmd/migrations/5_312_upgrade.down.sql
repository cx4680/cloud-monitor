DROP TABLE IF EXISTS `t_tenant_alarm_template_rel`;
DROP TABLE IF EXISTS `t_alarm_item_template`;
DROP TABLE IF EXISTS `t_alarm_rule_template`;

ALTER TABLE `t_alarm_rule` DROP COLUMN `type`;
ALTER TABLE `t_alarm_rule` DROP COLUMN `template_biz_id`;
ALTER TABLE `t_alarm_rule` DROP COLUMN `combination`;
ALTER TABLE `t_alarm_rule` DROP COLUMN `period`;
ALTER TABLE `t_alarm_rule` DROP COLUMN `times`;

ALTER TABLE `t_alarm_rule`
    ADD `trigger_condition` json comment '条件表达式';

ALTER TABLE t_alarm_rule change COLUMN metric_code metric_name VARCHAR (100);

INSERT INTO t_alarm_rule(id, metric_name, trigger_condition, `level`, silences_time)
SELECT rule_biz_id,
       metric_code,
       trigger_condition,
       `level`,
       silences_time
FROM t_alarm_item;

DROP TABLE IF EXISTS `t_alarm_item`;

DELETE FROM t_monitor_item WHERE biz_id IN ('185','186','187','188','189','190','191','192','193','194','195','196','197','198','199','200','201','202','203','204','205','206','207','208','209','210','211','212','213');

DELETE FROM t_monitor_product WHERE abbreviation IN ('redis','mongo','cgw');
