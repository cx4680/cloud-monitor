
-- ----------------------------
-- Records of monitor_product
-- ----------------------------
INSERT INTO `monitor_product` (`id`, `name`, `status`, `description`, `create_user`, `create_time`, `route`) VALUES ('1', '云服务器ECS', 1, 'ecs', NULL, NULL, '/productmonitoring/ecs');
INSERT INTO `monitor_product` (`id`, `name`, `status`, `description`, `create_user`, `create_time`, `route`) VALUES ('2', '弹性公网IP', '1', 'eip', NULL, NULL, '/productmonitoring/eip');
INSERT INTO `monitor_product` (`id`, `name`, `status`, `description`, `create_user`, `create_time`, `route`) VALUES ('3', '负载均衡SLB', '1', 'slb', NULL, NULL, '/productmonitoring/slb');

INSERT INTO `monitor_product` (`id`, `name`, `status`, `description`, `create_user`, `create_time`, `route`) VALUES ('5', '云备份CBR', '1', 'cbr', NULL, NULL, '/productmonitoring/cbr');
INSERT INTO `monitor_product` (`id`, `name`, `status`, `description`, `create_user`, `create_time`, `route`) VALUES ('6', 'NAT网关', '1', 'nat', NULL, NULL, '/productmonitoring/nat');

-- ----------------------------
-- Records of monitor_item
-- ----------------------------
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('1', '1', 'CPU使用率', 'ecs_cpu_usage', 'instance', '100 - (100 * (sum by(instance) (irate(ecs_cpu_seconds_total{mode="idle",$INSTANCE}[3m])) / sum by(instance) (irate(ecs_cpu_seconds_total{$INSTANCE}[3m]))))', '', '', '%', null, '2', '1', '1', null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('2', '1', 'CPU1分钟平均负载', 'ecs_load1', 'instance', 'ecs_load1{$INSTANCE}', '', null, '', null, '2', '1', '1', null, null, null, "'windows'!='$OS_TYPE'");
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('3', '1', 'CPU5分钟平均负载', 'ecs_load5', 'instance', 'ecs_load5{$INSTANCE}', '', 'max', '', null, '2', '1', '1', null, null, null, "'windows'!='$OS_TYPE'");
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`) VALUES ('4', '1', 'CPU15分钟平均负载', 'ecs_load15', 'instance', 'ecs_load15{$INSTANCE}', null, null, null, null, '2', '1', '1', null, null, null, "'windows'!='$OS_TYPE'");
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('5', '1', '内存使用量', 'ecs_memory_used', 'instance', 'ecs_memory_MemTotal_bytes{$INSTANCE} - ecs_memory_MemFree_bytes', null, null, 'B', null, '2', '1', '1', null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('6', '1', '内存使用率', 'ecs_memory_usage', 'instance', '100 * ((ecs_memory_MemTotal_bytes{$INSTANCE} - ecs_memory_MemFree_bytes) / ecs_memory_MemTotal_bytes)', null, null, '%', null, '2', '1', '1', null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('7', '1', '磁盘使用率', 'ecs_disk_usage', 'instance,device', '100 * ((ecs_filesystem_size_bytes{$INSTANCE} - ecs_filesystem_free_bytes) / ecs_filesystem_size_bytes)', null, null, '%', null, '2', '1', '1', null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('8', '1', '磁盘读速率', 'ecs_disk_read_rate', 'instance,device', 'irate(ecs_disk_read_bytes_total{$INSTANCE}[3m])', null, null, 'B/s', null, '2', '1', '1', null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('9', '1', '磁盘写速率', 'ecs_disk_write_rate', 'instance,device', 'irate(ecs_disk_written_bytes_total{$INSTANCE}[3m])', null, null, 'B/s', null, '2', '1', '1', null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('10', '1', '磁盘读IOPS', 'ecs_disk_read_iops', 'instance,device', 'irate(ecs_disk_reads_completed_total{$INSTANCE}[3m])', null, null, '次', null, '2', '1', '1', null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('11', '1', '磁盘写IOPS', 'ecs_disk_write_iops', 'instance,device', 'irate(ecs_disk_writes_completed_total{$INSTANCE}[3m])', null, null, '次', null, '2', '1', '1', null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('12', '1', '流入带宽', 'ecs_network_receive_rate', 'instance,device', 'irate(ecs_network_receive_bytes_total{$INSTANCE}[3m]) / 1024 / 1024 * 8', null, null, '	Mbps', null, '2', '1', '1', null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('13', '1', '流出带宽', 'ecs_network_transmit_rate', 'instance,device', 'irate(ecs_network_transmit_bytes_total{$INSTANCE}[3m]) / 1024 / 1024 * 8', null, null, '	Mbps', null, '2', '1', '1', null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('14', '1', '包接收速率', 'ecs_network_receive_packets_rate', 'instance,device', 'irate(ecs_network_receive_packets_total{$INSTANCE}[3m])', null, null, '个/s', null, '2', '1', '1', null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('15', '1', '包发送速率', 'ecs_network_transmit_packets_rate', 'instance,device', 'irate(ecs_network_transmit_packets_total{$INSTANCE}[3m])', null, null, '个/s', null, '2', '1', '1', null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('16', '1', '(基础)CPU的平均使用率', 'ecs_cpu_base_usage', 'instance', '100 * avg(irate(ecs_base_vcpu_seconds{$INSTANCE}[6m]))', null, null, '%', null, '1', '1', '1', null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('17', '1', '(基础)磁盘读速率', 'ecs_disk_base_read_rate', 'instance,drive', 'irate(ecs_base_storage_traffic_bytes_total{type="read",$INSTANCE}[6m])', null, null, 'B/s', null, '1', '1', '1', null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('18', '1', '(基础)磁盘写速率', 'ecs_disk_base_write_rate', 'instance,drive', 'irate(ecs_base_storage_traffic_bytes_total{type="write",$INSTANCE}[6m])', null, null, 'B/s', null, '1', '1', '1', null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('19', '1', '(基础)网卡下行带宽', 'ecs_network_base_receive_rate', 'instance,interface', 'irate(ecs_base_network_traffic_bytes_total{type="rx",$INSTANCE}[6m])', null, '', 'B/s', null, '1', '1', '1', null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('20', '1', '(基础)网卡上行带宽', 'ecs_network_base_transmit_rate', 'instance,interface', 'irate(ecs_base_network_traffic_bytes_total{type="tx",$INSTANCE}[6m])', null, null, 'B/s', null, '1', '1', '1', null, null, null);

INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('21', '2', '出网带宽', 'eip_upstream_bandwidth', 'instance', 'sum(eip_upstream_bits_rate{$INSTANCE}) by (instance,eip)', NULL, NULL, 'bps', NULL, NULL, '1', '1', NULL, NULL, NULL);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('22', '2', '入网带宽', 'eip_downstream_bandwidth', 'instance', 'sum(eip_downstream_bits_rate{$INSTANCE}) by (instance,eip)', NULL, NULL, 'bps', NULL, NULL, '1', '1', NULL, NULL, NULL);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('23', '2', '出网流量', 'eip_upstream', 'instance', '((sum(eip_upstream_bits_rate{$INSTANCE}) by (instance,eip))/8)*60', NULL, NULL, 'Byte', NULL, NULL, '1', '1', NULL, NULL, NULL);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('24', '2', '入网流量', 'eip_downstream', 'instance', '((sum(eip_downstream_bits_rate{$INSTANCE}) by (instance,eip))/8)*60', NULL, NULL, 'Byte', NULL, NULL, '1', '1', NULL, NULL, NULL);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('25', '2', '出网带宽使用率', 'eip_upstream_bandwidth_usage', 'instance', '(sum(eip_upstream_bits_rate{$INSTANCE}) by (instance,eip) / avg(eip_config_upstream_bandwidth{$INSTANCE}) by (instance,eip)) * 100', NULL, NULL, '%', NULL, NULL, '1', '1', NULL, NULL, NULL);

INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('26', '3', '出网带宽', 'slb_out_bandwidth', 'instance,slb_listener_id', 'sum by(instance) (Slb_http_bps_out_rate{$INSTANCE})', NULL, NULL, 'bps', NULL, NULL, '1', '1', NULL, NULL, NULL);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('27', '3', '入网带宽', 'slb_in_bandwidth', 'instance,slb_listener_id', 'sum by(instance) (Slb_http_bps_in_rate{$INSTANCE})', NULL, NULL, 'bps', NULL, NULL, '1', '1', NULL, NULL, NULL);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('30', '3', '并发连接数', 'slb_max_connection', 'instance,slb_listener_id', 'sum by(instance) (Slb_all_connection_count{$INSTANCE})', NULL, NULL, '个', NULL, NULL, '1', '1', NULL, NULL, NULL);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('31', '3', '活跃连接数', 'slb_active_connection', 'instance,slb_listener_id', 'sum by (instance)(Slb_all_est_connection_count{$INSTANCE})', NULL, NULL, '个', NULL, NULL, '1', '1', NULL, NULL, NULL);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('32', '3', '非活跃连接数', 'slb_inactive_connection', 'instance,slb_listener_id', 'sum by (instance) (Slb_all_none_est_connection_count{$INSTANCE})', NULL, NULL, '个', NULL, NULL, '1', '1', NULL, NULL, NULL);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('33', '3', '新建连接数', 'slb_new_connection', 'instance,slb_listener_id', 'sum by(instance) (Slb_new_connection_rate{$INSTANCE})', NULL, NULL, '个/s', NULL, NULL, '1', '1', NULL, NULL, NULL);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('34', '3', '丢弃连接数', 'slb_drop_connection', 'instance,slb_listener_id', 'sum by(instance)(Slb_drop_connection_rate{$INSTANCE})', NULL, NULL, '个/s', NULL, NULL, '1', '1', NULL, NULL, NULL);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('39', '3', '异常后端云服务器数', 'slb_unhealthyserver', 'instance', 'avg by(instance) (Slb_unhealthy_server_count{$INSTANCE})', NULL, NULL, '个', NULL, NULL, '1', '1', NULL, NULL, NULL);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('40', '3', '正常后端云服务器数', 'slb_healthyserver', 'instance', 'avg by(instance) (Slb_healthy_server_count{$INSTANCE})', NULL, NULL, '个', NULL, NULL, '1', '1', NULL, NULL, NULL);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('41', '3', '7层协议查询速率', 'slb_qps', 'instance,slb_listener_id', 'sum by(instance)(Slb_request_rate{$INSTANCE})', NULL, NULL, '个/s', NULL, NULL, '1', '1', NULL, NULL, NULL);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('42', '3', '7层协议返回客户端2xx状态码数', 'slb_statuscode2xx', 'instance', 'sum by(instance) (Slb_http_2xx_rate{$INSTANCE})', NULL, NULL, '个/s', NULL, NULL, '1', '1', NULL, NULL, NULL);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('43', '3', '7层协议返回客户端3xx状态码数', 'slb_statuscode3xx', 'instance', 'sum by(instance) (Slb_http_3xx_rate{$INSTANCE})', NULL, NULL, '个/s', NULL, NULL, '1', '1', NULL, NULL, NULL);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('44', '3', '7层协议返回客户端4xx状态码数', 'slb_statuscode4xx', 'instance', 'sum by(instance) (Slb_http_4xx_rate{$INSTANCE})', NULL, NULL, '个/s', NULL, NULL, '1', '1', NULL, NULL, NULL);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('45', '3', '7层协议返回客户端5xx状态码数', 'slb_statuscode5xx', 'instance', 'sum by(instance) (Slb_http_5xx_rate{$INSTANCE})', NULL, NULL, '个/s', NULL, NULL, '1', '1', NULL, NULL, NULL);

INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('52', '5', '存储库总量', 'cbr_vault_size', 'instance', 'cbr_vault_size{$INSTANCE}', NULL, NULL, 'B', NULL, NULL, '1', '1', NULL, NULL, NULL);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('53', '5', '存储库使用量', 'cbr_vault_used', 'instance', 'cbr_vault_used{$INSTANCE}', NULL, NULL, 'B', NULL, NULL, '1', '1', NULL, NULL, NULL);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('54', '5', '存储库使用率', 'cbr_vault_usage_rate', 'instance', 'cbr_vault_used{$INSTANCE} / cbr_vault_size{$INSTANCE} * 100', NULL, NULL, '%', NULL, NULL, '1', '1', NULL, NULL, NULL);

INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('55', '6', 'SNAT连接数', 'snat_connection', 'instance', 'sum by (instance)(Nat_snat_total_connection_count{$INSTANCE})', NULL, NULL, '个', NULL, NULL, '1', '1', NULL, NULL, NULL);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('56', '6', '入方向带宽', 'inbound_bandwidth', 'instance', 'sum by (instance)(rate(Nat_recv_bytes_total_count{$INSTANCE}[1m])*8)', NULL, NULL, 'bps', NULL, NULL, '1', '1', NULL, NULL, NULL);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('57', '6', '出方向带宽', 'outbound_bandwidth', 'instance', 'sum by (instance)(rate(Nat_send_bytes_total_count{$INSTANCE}[1m])*8)', NULL, NULL, 'bps', NULL, NULL, '1', '1', NULL, NULL, NULL);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('58', '6', '入方向流量', 'inbound_traffic', 'instance', 'sum by (instance)(Nat_recv_bytes_total_count{$INSTANCE})', NULL, NULL, 'byte', NULL, NULL, '1', '1', NULL, NULL, NULL);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('59', '6', '出方向流量', 'outbound_traffic', 'instance', 'sum by (instance)(Nat_send_bytes_total_count{$INSTANCE})', NULL, NULL, 'byte', NULL, NULL, '1', '1', NULL, NULL, NULL);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('60', '6', '入方向PPS', 'inbound_pps', 'instance', 'sum by (instance)(rate(Nat_recv_packets_total_count{$INSTANCE}[1m]))', NULL, NULL, 'pps', NULL, NULL, '1', '1', NULL, NULL, NULL);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('61', '6', '出方向PPS', 'outbound_pps', 'instance', 'sum by (instance)(rate(Nat_send_packets_total_count{$INSTANCE}[1m]))', NULL, NULL, 'pps', NULL, NULL, '1', '1', NULL, NULL, NULL);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('62', '6', 'SNAT连接数使用率', 'snat_connection_ratio', 'instance', 'sum by (instance)(Nat_snat_total_connection_count{$INSTANCE}) / avg by (instance)(Nat_nat_max_connection_count{$INSTANCE}) *100', NULL, NULL, '%', NULL, NULL, '1', '1', NULL, NULL, NULL);

INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('63', '1', '(基础)磁盘读iops', 'ecs_disk_base_read_iops', 'instance,drive', 'sum by(instance) (irate(ecs_base_storage_iops_total{type="read",$INSTANCE}[15m])) by (drive)', null, null, '次', null, '1', '1', '1', null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('64', '1', '(基础)磁盘写iops', 'ecs_disk_base_write_iops', 'instance,drive', 'sum by(instance) (irate(ecs_base_storage_iops_total{type="write",$INSTANCE}[15m])) by (drive)', null, null, '次', null, '1', '1', '1', null, null, null);

INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('65', '1', '磁盘剩余存储量', 'ecs_filesystem_free_bytes', 'instance,device', 'ecs_filesystem_free_bytes{$INSTANCE}', null, null, 'GB', null, '2', '1', '1', null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('66', '1', '磁盘已用存储量', 'ecs_disk_used', 'instance,device', 'ecs_filesystem_size_bytes{$INSTANCE} - ecs_filesystem_free_bytes', null, null, 'GB', null, '2', '1', '1', null, null, null);
INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('67', '1', '磁盘存储总量', 'ecs_filesystem_size_bytes', 'instance,device', 'ecs_filesystem_size_bytes{$INSTANCE}', null, null, 'GB', null, '2', '1', '1', null, null, null);

-- ----------------------------
-- Records of config_item
-- ----------------------------
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('', '-1', '监控周期', NULL, '', 0, NULL);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('1', '-1', '统计周期', NULL, NULL, 0, NULL);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('10', '2', '持续3个周期', '3', '3', 1, NULL);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('11', '2', '持续5个周期', '5', '5', 2, NULL);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('12', '3', '平均值', 'Average', 'avg', 0, NULL);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('13', '3', '最大值', 'Maximum', 'max', 1, NULL);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('14', '3', '最小值', 'Minimum', 'min', 2, NULL);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('15', '4', '大于', 'greater', '>', 0, NULL);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('16', '4', '大于等于', 'greaterOrEqual', '>=', 1, NULL);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('17', '4', '小于', 'less', '<', 2, NULL);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('18', '4', '小于等于', 'lessOrEqual', '<=', 3, NULL);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('19', '4', '等于', 'equal', '=', 4, NULL);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('2', '-1', '持续周期', NULL, NULL, 0, NULL);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('20', '4', '不等于', 'notEqual', '!=', 5, NULL);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('21', '-1', '概览监控项', NULL, NULL, 0, NULL);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('22', '21', 'CPU使用率（操作系统）', NULL, 'ecs_cpu_usage', 0, NULL);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('23', '21', '内存使用率（操作系统）', NULL, 'ecs_memory_usage', 1, NULL);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('28', '-1', '监控周期', NULL, NULL, 0, NULL);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('29', '28', '紧急', '1', 'MAIN', 0, NULL);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('3', '-1', '统计方式', NULL, NULL, 0, NULL);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('30', '28', '重要', '2', 'MARJOR', 0, NULL);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('31', '28', '次要', '3', 'MINOR', 0, NULL);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('32', '28', '提醒', '4', 'WARN', 0, NULL);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('33', '-1', '通知方式', NULL, NULL, 0, NULL);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('34', '33', '全部', '1', 'all', 0, NULL);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('35', '33', '邮件', '2', 'email', 0, NULL);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('36', '33', '短信', '3', 'sms', 0, NULL);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('4', '-1', '对比方式', NULL, NULL, 0, NULL);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('5', '-1', '监控数据', NULL, NULL, 0, NULL);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('51', '5', '0-3H', '0,3', '60', 1, NULL);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('52', '5', '3H-12H', '3,12', '180', 2, NULL);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('53', '5', '12H-3D', '12,72', '900', 3, NULL);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('54', '5', '3D-10D', '72,240', '2700', 4, NULL);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('6', '1', '5分钟', '300', '5m', 1, NULL);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('7', '1', '15分钟', '900', '15m', 2, NULL);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('8', '1', '30分钟', '1800', '30m', 3, NULL);
INSERT INTO `config_item` (`id`, `pid`, `name`, `code`, `data`, `sort_id`,`remark`) VALUES ('9', '2', '持续1个周期', '1', '1', 0, NULL);
