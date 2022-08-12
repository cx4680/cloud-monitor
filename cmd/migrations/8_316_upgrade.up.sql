UPDATE t_monitor_item SET metrics_linux = 'ecs_processes_top5Cpus{$INSTANCE}', labels = 'instance,pid' WHERE biz_id = '220';
UPDATE t_monitor_item SET metrics_linux = 'ecs_processes_top5Mems{$INSTANCE}', labels = 'instance,pid' WHERE biz_id = '221';
UPDATE t_monitor_item SET metrics_linux = 'ecs_procs_running{$INSTANCE}' WHERE biz_id = '222';
UPDATE t_monitor_item SET metrics_linux = 'ecs_processes_top5Fds{$INSTANCE}', labels = 'instance,pid' WHERE biz_id = '223';

UPDATE t_monitor_product SET page_url = '/inner/cmq/v1/kafka/monitor/listAllCluster', status = '1' WHERE abbreviation IN ('kafka');

UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(redis_connected_clients{$INSTANCE})' WHERE biz_id = '187';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(redis_tps{$INSTANCE})' WHERE biz_id = '188';

UPDATE t_monitor_item SET metric_name = 'nat_snat_connection', show_expression = '{{eq .OSTYPE \"\"}}' WHERE biz_id = '55';
UPDATE t_monitor_item SET metric_name = 'nat_inbound_bandwidth', show_expression = '{{eq .OSTYPE \"\"}}' WHERE biz_id = '56';
UPDATE t_monitor_item SET metric_name = 'nat_outbound_bandwidth', show_expression = '{{eq .OSTYPE \"\"}}' WHERE biz_id = '57';
UPDATE t_monitor_item SET metric_name = 'nat_inbound_traffic', show_expression = '{{eq .OSTYPE \"\"}}' WHERE biz_id = '58';
UPDATE t_monitor_item SET metric_name = 'nat_outbound_traffic', show_expression = '{{eq .OSTYPE \"\"}}' WHERE biz_id = '59';
UPDATE t_monitor_item SET metric_name = 'nat_inbound_pps', show_expression = '{{eq .OSTYPE \"\"}}' WHERE biz_id = '60';
UPDATE t_monitor_item SET metric_name = 'nat_outbound_pps', show_expression = '{{eq .OSTYPE \"\"}}' WHERE biz_id = '61';
UPDATE t_monitor_item SET metric_name = 'nat_snat_connection_ratio', show_expression = '{{eq .OSTYPE \"\"}}' WHERE biz_id = '62';

INSERT INTO `t_monitor_item` (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('226', '6', 'NAT连接数', 'nat_e_total_connection', 'instance', 'sum by (instance)(Nat_snat_total_connection_count{$INSTANCE})', NULL, NULL, '个', NULL, NULL, 1, 1, NULL, NULL, NULL, '{{eq .OSTYPE \"nat-e\"}}', 'chart,rule');
INSERT INTO `t_monitor_item` (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('227', '6', '入方向带宽', 'nat_e_inbound_bandwidth', 'instance', 'sum by (instance)(rate(Nat_recv_bytes_total_count{$INSTANCE}[1m])*8)', NULL, NULL, '个', NULL, NULL, 1, 1, NULL, NULL, NULL, '{{eq .OSTYPE \"nat-e\"}}', 'chart,rule');
INSERT INTO `t_monitor_item` (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('228', '6', '出方向带宽', 'nat_e_outbound_bandwidth', 'instance', 'sum by (instance)(rate(Nat_send_bytes_total_count{$INSTANCE}[1m])*8)', NULL, NULL, '个', NULL, NULL, 1, 1, NULL, NULL, NULL, '{{eq .OSTYPE \"nat-e\"}}', 'chart,rule');
INSERT INTO `t_monitor_item` (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('229', '6', '入方向流量', 'nat_e_inbound_traffic', 'instance', 'sum by (instance)(Nat_recv_bytes_total_count{$INSTANCE})', NULL, NULL, '个', NULL, NULL, 1, 1, NULL, NULL, NULL, '{{eq .OSTYPE \"nat-e\"}}', 'chart,rule');
INSERT INTO `t_monitor_item` (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('230', '6', '出方向流量', 'nat_e_outbound_traffic	', 'instance', 'sum by (instance)(Nat_send_bytes_total_count{$INSTANCE})', NULL, NULL, '个', NULL, NULL, 1, 1, NULL, NULL, NULL, '{{eq .OSTYPE \"nat-e\"}}', 'chart,rule');
INSERT INTO `t_monitor_item` (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('231', '6', '入方向PPS', 'nat_e_inbound_pps', 'instance', 'sum by (instance)(rate(Nat_recv_packets_total_count{$INSTANCE}[1m]))', NULL, NULL, '个', NULL, NULL, 1, 1, NULL, NULL, NULL, '{{eq .OSTYPE \"nat-e\"}}', 'chart,rule');
INSERT INTO `t_monitor_item` (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('232', '6', '出方向PPS', 'nat_e_outbound_pps', 'instance', 'sum by (instance)(rate(Nat_send_packets_total_count{$INSTANCE}[1m]))', NULL, NULL, '个', NULL, NULL, 1, 1, NULL, NULL, NULL, '{{eq .OSTYPE \"nat-e\"}}', 'chart,rule');
INSERT INTO `t_monitor_item` (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('233', '6', 'NAT连接数使用率', 'nat_e_total_connection_ratio', 'instance', 'sum by (instance)(Nat_snat_total_connection_count{$INSTANCE}) / avg by (instance)(Nat_nat_max_connection_count{$INSTANCE}) *100', NULL, NULL, '个', NULL, NULL, 1, 1, NULL, NULL, NULL, '{{eq .OSTYPE \"nat-e\"}}', 'chart,rule');

ALTER TABLE t_monitor_product ADD COLUMN iam_page_url varchar(256) COMMENT 'iam请求路径';
UPDATE t_monitor_product SET iam_page_url = '/compute/ecs/cbc/pageList' WHERE abbreviation IN ('ecs');
UPDATE t_monitor_product SET iam_page_url = '/eip/inner/eipInfoList' WHERE abbreviation IN ('eip');
UPDATE t_monitor_product SET iam_page_url = '/slb/list' WHERE abbreviation IN ('slb');
UPDATE t_monitor_product SET iam_page_url = '/noauth/backup/vault/pageList' WHERE abbreviation IN ('cbr');
UPDATE t_monitor_product SET iam_page_url = '/nat-gw/inner/nat/page' WHERE abbreviation IN ('nat');
UPDATE t_monitor_product SET iam_page_url = '/v1/mysql/instance' WHERE abbreviation IN ('mysql');
UPDATE t_monitor_product SET iam_page_url = '/v1/dm/instance' WHERE abbreviation IN ('dm');
UPDATE t_monitor_product SET iam_page_url = '/v1/pg/instance/' WHERE abbreviation IN ('postgresql');
UPDATE t_monitor_product SET iam_page_url = '/kafka/v1/cluster/listAllCluster' WHERE abbreviation IN ('kafka');
UPDATE t_monitor_product SET iam_page_url = '/compute/bms/ops/v1/tenants/{tenantId}/servers' WHERE abbreviation IN ('bms');
UPDATE t_monitor_product SET iam_page_url = '/compute/ebms/ops/v1/tenants/{tenantId}/servers' WHERE abbreviation IN ('ebms');
UPDATE t_monitor_product SET iam_page_url = '/v1/redis/instance' WHERE abbreviation IN ('redis');
UPDATE t_monitor_product SET iam_page_url = '/v2/mongo/instance' WHERE abbreviation IN ('mongo');
UPDATE t_monitor_product SET iam_page_url = '/gateway/instance/page' WHERE abbreviation IN ('cgw');
