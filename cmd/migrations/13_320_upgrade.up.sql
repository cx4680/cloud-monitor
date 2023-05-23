UPDATE t_monitor_product SET page_url = '/slb/inner/monitor/list' WHERE abbreviation = 'slb';
UPDATE t_monitor_product SET iam_page_url = '/slb/list' WHERE abbreviation IN ('slb');
UPDATE t_monitor_product SET iam_page_url = '/nat-gw/nat/page' WHERE abbreviation IN ('nat');
UPDATE t_monitor_product SET iam_page_url = '/nat-gw/nat/page' WHERE abbreviation IN ('nat-e');

DELETE FROM t_monitor_item WHERE metric_name = 'ecs_memory_base_usage';
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('239', '1', '(基础)内存使用率', 'ecs_memory_base_usage', 'instance', '(1-ecs_base_memory_unused_bytes{$INSTANCE}/ecs_base_memory_available_bytes)*100', null, null, '%', null, '1', '1', '1', null, null, null, null, 'chart,rule,scaling');

DELETE FROM t_monitor_item WHERE metric_name IN ('ecs_disk_read_time','ecs_disk_write_time');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('237', '1', '读延迟', 'ecs_disk_read_time', 'instance,device', 'rate(ecs_disk_read_time_seconds_total{$INSTANCE}[5m]) / clamp_min(rate(ecs_disk_reads_completed_total{}[5m]),1) * 1000', null, null, 'ms', null, '2', '1', '1', null, null, null, null, 'chart,rule');
INSERT INTO t_monitor_item (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('238', '1', '写延迟', 'ecs_disk_write_time', 'instance,device', 'rate(ecs_disk_write_time_seconds_total{$INSTANCE}[5m]) / clamp_min(rate(ecs_disk_writes_completed_total{}[5m]),1) * 1000', null, null, 'ms', null, '2', '1', '1', null, null, null, null, 'chart,rule');

UPDATE t_monitor_item SET metrics_linux = '100 * ((sum by(instance)(ecs_memory_MemTotal_bytes{$INSTANCE}) - (sum by(instance)(ecs_memory_MemFree_bytes{}) + sum by(instance)(ecs_memory_Cached_bytes{}))) / sum by(instance)(ecs_memory_MemTotal_bytes{}))' WHERE metric_name = 'ecs_memory_usage';
UPDATE t_monitor_item SET metrics_linux = '100 * ((sum by(instance)(ecs_memory_MemTotal_bytes{$INSTANCE}) - (sum by(instance)(ecs_memory_MemFree_bytes{}) + sum by(instance)(ecs_memory_Cached_bytes{}))) / sum by(instance)(ecs_memory_MemTotal_bytes{}))' WHERE metric_name = 'bms_memory_usage';
UPDATE t_monitor_item SET metrics_linux = '100 * ((sum by(instance)(ecs_memory_MemTotal_bytes{$INSTANCE}) - (sum by(instance)(ecs_memory_MemFree_bytes{}) + sum by(instance)(ecs_memory_Cached_bytes{}))) / sum by(instance)(ecs_memory_MemTotal_bytes{}))' WHERE metric_name = 'ebms_memory_usage';

UPDATE `t_monitor_product` SET name = '云数据库 RDS DM' WHERE biz_id = '9';
