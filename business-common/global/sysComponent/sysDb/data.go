package sysDb

import (
	commonModels "code.cestc.cn/ccos-ops/cloud-monitor/business-common/models"
	"gorm.io/gorm"
)

type DBInitializer struct {
	DB      *gorm.DB
	Fetches []InitializerFetch
}

func (i *DBInitializer) Initnitialization() error {
	for _, b := range i.Fetches {
		t, s := b.Fetch(i.DB)
		if err := i.DB.AutoMigrate(t...); err != nil {
			return err
		}
		for _, sql := range s {
			i.DB.Exec(sql)
		}
	}
	return nil
}

type InitializerFetch interface {
	Fetch(db *gorm.DB) ([]interface{}, []string)
}

type CommonInitializerFetch struct {
}

func (c *CommonInitializerFetch) Fetch(db *gorm.DB) ([]interface{}, []string) {
	var tables []interface{}
	var sqls []string

	tables = append(tables, &commonModels.ConfigItem{}, &commonModels.MonitorItem{})
	if !db.Migrator().HasTable(&commonModels.MonitorItem{}) {
		sqls = append(sqls, "INSERT INTO `monitor_product` VALUES ('1', '云服务器ECS', 1, '云服务器ECS', NULL, NULL);",
			"INSERT INTO `monitor_product` (`id`, `name`, `status`, `description`, `create_user`, `create_time`) VALUES ('2', '弹性公网IP', '1', '弹性公网IP', NULL, NULL);",
			"INSERT INTO `monitor_product` (`id`, `name`, `status`, `description`, `create_user`, `create_time`) VALUES ('3', '负载均衡SLB', '1', '负载均衡SLB', NULL, NULL);",

			"INSERT INTO `monitor_item` VALUES ('1', '1', 'CPU使用率', 'ecs_cpu_usage', 'instance', '100 - (100 * (sum(irate(ecs_cpu_seconds_total{mode=\"idle\",$INSTANCE}[3m])) / sum(irate(ecs_cpu_seconds_total{$INSTANCE}[3m]))))', '', '', '%', null, '2', '1', '1', null, null, null);",
			"INSERT INTO `monitor_item` VALUES ('2', '1', 'cpu1分钟平均负载', 'ecs_load1', 'instance', 'ecs_load1{$INSTANCE}', '', null, '', null, '2', '1', '1', null, null, null);",
			"INSERT INTO `monitor_item` VALUES ('3', '1', 'cpu5分钟平均负载', 'ecs_load5', 'instance', 'ecs_load5{$INSTANCE}', '', 'max', '', null, '2', '1', '1', null, null, null);",
			"INSERT INTO `monitor_item` VALUES ('4', '1', 'cpu15分钟平均负载', 'ecs_load15', 'instance', 'ecs_load15{$INSTANCE}', null, null, null, null, '2', '1', '1', null, null, null);",
			"INSERT INTO `monitor_item` VALUES ('5', '1', '内存使用量', 'ecs_memory_used', 'instance', 'ecs_memory_MemTotal_bytes{$INSTANCE} - ecs_memory_MemFree_bytes', null, null, 'B', null, '2', '1', '1', null, null, null);",
			"INSERT INTO `monitor_item` VALUES ('6', '1', '内存使用率', 'ecs_memory_usage', 'instance', '100 * ((ecs_memory_MemTotal_bytes{$INSTANCE} - ecs_memory_MemFree_bytes) / ecs_memory_MemTotal_bytes)', null, null, '%', null, '2', '1', '1', null, null, null);",
			"INSERT INTO `monitor_item` VALUES ('7', '1', '磁盘使用率', 'ecs_disk_usage', 'instance,device', '100 * ((ecs_filesystem_size_bytes{$INSTANCE} - ecs_filesystem_free_bytes) / ecs_filesystem_size_bytes)', null, null, '%', null, '2', '1', '1', null, null, null);",
			"INSERT INTO `monitor_item` VALUES ('8', '1', '磁盘读速率', 'ecs_disk_read_rate', 'instance,device', 'irate(ecs_disk_read_bytes_total{$INSTANCE}[3m])', null, null, 'B/s', null, '2', '1', '1', null, null, null);",
			"INSERT INTO `monitor_item` VALUES ('9', '1', '磁盘写速率', 'ecs_disk_write_rate', 'instance,device', 'irate(ecs_disk_written_bytes_total{$INSTANCE}[3m])', null, null, 'B/s', null, '2', '1', '1', null, null, null);",
			"INSERT INTO `monitor_item` VALUES ('10', '1', '磁盘读IOPS', 'ecs_disk_read_iops', 'instance,device', 'irate(ecs_disk_reads_completed_total{$INSTANCE}[3m])', null, null, '次', null, '2', '1', '1', null, null, null);",
			"INSERT INTO `monitor_item` VALUES ('11', '1', '磁盘写IOPS', 'ecs_disk_write_iops', 'instance,device', 'irate(ecs_disk_writes_completed_total{$INSTANCE}[3m])', null, null, '次', null, '2', '1', '1', null, null, null);",
			"INSERT INTO `monitor_item` VALUES ('12', '1', '流入带宽', 'ecs_network_receive_rate', 'instance,device', 'irate(ecs_network_receive_bytes_total{$INSTANCE}[3m]) / 1024 / 1024 * 8', null, null, '	Mbps', null, '2', '1', '1', null, null, null);",
			"INSERT INTO `monitor_item` VALUES ('13', '1', '流出带宽', 'ecs_network_transmit_rate', 'instance,device', 'irate(ecs_network_transmit_bytes_total{$INSTANCE}[3m]) / 1024 / 1024 * 8', null, null, '	Mbps', null, '2', '1', '1', null, null, null);",
			"INSERT INTO `monitor_item` VALUES ('14', '1', '包接收速率', 'ecs_network_receive_packets_rate', 'instance,device', 'irate(ecs_network_receive_packets_total{$INSTANCE}[3m])', null, null, '个/s', null, '2', '1', '1', null, null, null);",
			"INSERT INTO `monitor_item` VALUES ('15', '1', '包发送速率', 'ecs_network_transmit_packets_rate', 'instance,device', 'irate(ecs_network_transmit_packets_total{$INSTANCE}[3m])', null, null, '个/s', null, '2', '1', '1', null, null, null);",
			"INSERT INTO `monitor_item` VALUES ('16', '1', '(基础)CPU的平均使用率', 'ecs_cpu_base_usage', 'instance', '100 * avg(irate(ecs_base_vcpu_seconds{$INSTANCE}[6m]))', null, null, '%', null, '1', '1', '1', null, null, null);",
			"INSERT INTO `monitor_item` VALUES ('17', '1', '(基础)磁盘读速率', 'ecs_disk_base_read_rate', 'instance,drive', 'irate(ecs_base_storage_traffic_bytes_total{type=\"read\",$INSTANCE}[6m])', null, null, 'B/s', null, '1', '1', '1', null, null, null);",
			"INSERT INTO `monitor_item` VALUES ('18', '1', '(基础)磁盘写速率', 'ecs_disk_base_write_rate', 'instance,drive', 'irate(ecs_base_storage_traffic_bytes_total{type=\"write\",$INSTANCE}[6m])', null, null, 'B/s', null, '1', '1', '1', null, null, null);",
			"INSERT INTO `monitor_item` VALUES ('19', '1', '(基础)网卡下行带宽', 'ecs_network_base_receive_rate', 'instance,interface', 'irate(ecs_base_network_traffic_bytes_total{type=\"rx\",$INSTANCE}[6m])', null, '', 'B/s', null, '1', '1', '1', null, null, null);",
			"INSERT INTO `monitor_item` VALUES ('20', '1', '(基础)网卡上行带宽', 'ecs_network_base_transmit_rate', 'instance,interface', 'irate(ecs_base_network_traffic_bytes_total{type=\"tx\",$INSTANCE}[6m])', null, null, 'B/s', null, '1', '1', '1', null, null, null);",

			"INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('21', '2', '出网带宽', 'eip_upstream_bandwidth', 'instance', 'rate(eip_upstream_total_bits{$INSTANCE}[1m])', NULL, NULL, 'bps', NULL, NULL, '1', '1', NULL, NULL, NULL);",
			"INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('22', '2', '入网带宽', 'eip_downstream_bandwidth', 'instance', 'rate(eip_downstream_total_bits{$INSTANCE}[1m])', NULL, NULL, 'bps', NULL, NULL, '1', '1', NULL, NULL, NULL);",
			"INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('23', '2', '出网流量', 'eip_upstream', 'instance', 'rate(eip_upstream_total_bits{$INSTANCE}[1m]) / 8', NULL, NULL, 'Byte', NULL, NULL, '1', '1', NULL, NULL, NULL);",
			"INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('24', '2', '入网流量', 'eip_downstream', 'instance', 'rate(eip_downstream_total_bits{$INSTANCE}[1m]) / 8', NULL, NULL, 'Byte', NULL, NULL, '1', '1', NULL, NULL, NULL);",
			"INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('25', '2', '出网带宽使用率', 'eip_upstream_bandwidth_usage', 'instance', '(sum by(eip) (rate(eip_upstream_total_bits{$INSTANCE}[1m])) / sum by(eip) (eip_config_upstream_bandwidth{$INSTANCE})) * 100', NULL, NULL, '%', NULL, NULL, '1', '1', NULL, NULL, NULL);",
			"INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('26', '3', '出网带宽', 'slb_out_bandwidth', 'instance,slb_listener_id', 'sum by(slb_listener_id) (rate(Slb_http_byte_out_total_count{$INSTANCE}[1m]) * 8)', NULL, NULL, 'bps', NULL, NULL, '1', '1', NULL, NULL, NULL);",
			"INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('27', '3', '入网带宽', 'slb_in_bandwidth', 'instance,slb_listener_id', 'sum by( slb_listener_id) (rate(Slb_http_byte_in_total_count{$INSTANCE}[1m]) * 8)', NULL, NULL, 'bps', NULL, NULL, '1', '1', NULL, NULL, NULL);",
			"INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('30', '3', '并发连接数', 'slb_max_connection', 'instance,slb_listener_id', 'sum by(slb_listener_id) (Slb_all_connection_count{$INSTANCE})', NULL, NULL, '个', NULL, NULL, '1', '1', NULL, NULL, NULL);",
			"INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('31', '3', '活跃连接数', 'slb_active_connection', 'instance,slb_listener_id', 'sum by(slb_listener_id) (Slb_all_est_connection_count{$INSTANCE})', NULL, NULL, '个', NULL, NULL, '1', '1', NULL, NULL, NULL);",
			"INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('32', '3', '非活跃连接数', 'slb_inactive_connection', 'instance,slb_listener_id', 'sum by(slb_listener_id) (Slb_all_none_est_connection_count{$INSTANCE})', NULL, NULL, '个', NULL, NULL, '1', '1', NULL, NULL, NULL);",
			"INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('33', '3', '新建连接数', 'slb_new_connection', 'instance,slb_listener_id', 'sum by(slb_listener_id) (rate(Slb_accepted_connection_total_count{$INSTANCE}[1m]))', NULL, NULL, '个/s', NULL, NULL, '1', '1', NULL, NULL, NULL);",
			"INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('34', '3', '丢弃连接数', 'slb_drop_connection', 'instance,slb_listener_id', 'sum by(slb_listener_id) (rate(Slb_accepted_connection_total_count{$INSTANCE}[1m]) - rate(Slb_handled_connection_total_count[1m]))', NULL, NULL, '个/s', NULL, NULL, '1', '1', NULL, NULL, NULL);",
			"INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('39', '3', '异常后端云服务器数', 'slb_unhealthyserver', 'instance', 'avg (Slb_unhealthy_server_count{$INSTANCE})', NULL, NULL, '个', NULL, NULL, '1', '1', NULL, NULL, NULL);",
			"INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('40', '3', '正常后端云服务器数', 'slb_healthyserver', 'instance', 'avg (Slb_healthy_server_count{$INSTANCE})', NULL, NULL, '个', NULL, NULL, '1', '1', NULL, NULL, NULL);",
			"INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('41', '3', '7层协议查询速率', 'slb_qps', 'instance,slb_listener_id', 'sum by(slb_listener_id) (rate(Slb_request_total_count{$INSTANCE}[1m]))', NULL, NULL, '个/s', NULL, NULL, '1', '1', NULL, NULL, NULL);",
			"INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('42', '3', '7层协议返回客户端2xx状态码数', 'slb_statuscode2xx', 'instance', 'sum (rate(Slb_http_2xx_total_count{$INSTANCE}[1m]))', NULL, NULL, '个/s', NULL, NULL, '1', '1', NULL, NULL, NULL);",
			"INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('43', '3', '7层协议返回客户端3xx状态码数', 'slb_statuscode3xx', 'instance', 'sum (rate(Slb_http_3xx_total_count{$INSTANCE}[1m]))', NULL, NULL, '个/s', NULL, NULL, '1', '1', NULL, NULL, NULL);",
			"INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('44', '3', '7层协议返回客户端4xx状态码数', 'slb_statuscode4xx', 'instance', 'sum (rate(Slb_http_4xx_total_count{$INSTANCE}[1m]))', NULL, NULL, '个/s', NULL, NULL, '1', '1', NULL, NULL, NULL);",
			"INSERT INTO `monitor_item` (`id`, `product_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`) VALUES ('45', '3', '7层协议返回客户端5xx状态码数', 'slb_statuscode5xx', 'instance', 'sum (rate(Slb_http_5xx_total_count{$INSTANCE}[1m]))', NULL, NULL, '个/s', NULL, NULL, '1', '1', NULL, NULL, NULL);",

			"INSERT INTO `config_item` VALUES ('', '-1', '监控周期', NULL, '', 0, NULL);",
			"INSERT INTO `config_item` VALUES ('1', '-1', '统计周期', NULL, NULL, 0, NULL);",
			"INSERT INTO `config_item` VALUES ('10', '2', '持续3个周期', '3', '3', 1, NULL);",
			"INSERT INTO `config_item` VALUES ('11', '2', '持续5个周期', '5', '5', 2, NULL);",
			"INSERT INTO `config_item` VALUES ('12', '3', '平均值', 'Average', 'avg', 0, NULL);",
			"INSERT INTO `config_item` VALUES ('13', '3', '最大值', 'Maximum', 'max', 1, NULL);",
			"INSERT INTO `config_item` VALUES ('14', '3', '最小值', 'Minimum', 'min', 2, NULL);",
			"INSERT INTO `config_item` VALUES ('15', '4', '大于', 'greater', '>', 0, NULL);",
			"INSERT INTO `config_item` VALUES ('16', '4', '大于等于', 'greaterOrEqual', '>=', 1, NULL);",
			"INSERT INTO `config_item` VALUES ('17', '4', '小于', 'less', '<', 2, NULL);",
			"INSERT INTO `config_item` VALUES ('18', '4', '小于等于', 'lessOrEqual', '<=', 3, NULL);",
			"INSERT INTO `config_item` VALUES ('19', '4', '等于', 'equal', '=', 4, NULL);",
			"INSERT INTO `config_item` VALUES ('2', '-1', '持续周期', NULL, NULL, 0, NULL);",
			"INSERT INTO `config_item` VALUES ('20', '4', '不等于', 'notEqual', '!=', 5, NULL);",
			"INSERT INTO `config_item` VALUES ('21', '-1', '概览监控项', NULL, NULL, 0, NULL);",
			"INSERT INTO `config_item` VALUES ('22', '21', 'CPU使用率（操作系统）', NULL, 'ecs_cpu_usage', 0, NULL);",
			"INSERT INTO `config_item` VALUES ('23', '21', '内存使用率（操作系统）', NULL, 'ecs_memory_usage', 1, NULL);",
			"INSERT INTO `config_item` VALUES ('28', '-1', '监控周期', NULL, NULL, 0, NULL);",
			"INSERT INTO `config_item` VALUES ('29', '28', '紧急', '1', 'MAIN', 0, NULL);",
			"INSERT INTO `config_item` VALUES ('3', '-1', '统计方式', NULL, NULL, 0, NULL);",
			"INSERT INTO `config_item` VALUES ('30', '28', '重要', '2', 'MARJOR', 0, NULL);",
			"INSERT INTO `config_item` VALUES ('31', '28', '次要', '3', 'MINOR', 0, NULL);",
			"INSERT INTO `config_item` VALUES ('32', '28', '提醒', '4', 'WARN', 0, NULL);",
			"INSERT INTO `config_item` VALUES ('33', '-1', '通知方式', NULL, NULL, 0, NULL);",
			"INSERT INTO `config_item` VALUES ('34', '33', '全部', '1', 'all', 0, NULL);",
			"INSERT INTO `config_item` VALUES ('35', '33', '邮件', '2', 'email', 0, NULL);",
			"INSERT INTO `config_item` VALUES ('36', '33', '短信', '3', 'sms', 0, NULL);",
			"INSERT INTO `config_item` VALUES ('4', '-1', '对比方式', NULL, NULL, 0, NULL);",
			"INSERT INTO `config_item` VALUES ('5', '-1', '监控数据', NULL, NULL, 0, NULL);",
			"INSERT INTO `config_item` VALUES ('51', '5', '0-3H', '0,3', '60', 1, NULL);",
			"INSERT INTO `config_item` VALUES ('52', '5', '3H-12H', '3,12', '180', 2, NULL);",
			"INSERT INTO `config_item` VALUES ('53', '5', '12H-3D', '12,72', '900', 3, NULL);",
			"INSERT INTO `config_item` VALUES ('54', '5', '3D-10D', '72,240', '2700', 4, NULL);",
			"INSERT INTO `config_item` VALUES ('6', '1', '5分钟', '300', '5m', 1, NULL);",
			"INSERT INTO `config_item` VALUES ('7', '1', '15分钟', '900', '15m', 2, NULL);",
			"INSERT INTO `config_item` VALUES ('8', '1', '30分钟', '1800', '30m', 3, NULL);",
			"INSERT INTO `config_item` VALUES ('9', '2', '持续1个周期', '1', '1', 0, NULL);",
		)
	}

	migrator := db.Migrator()
	if !migrator.HasTable(&commonModels.AlertContact{}) {
		tables = append(tables, &commonModels.AlertContact{})
	}
	if !migrator.HasTable(&commonModels.AlertContactGroup{}) {
		tables = append(tables, &commonModels.AlertContactGroup{})
	}
	if !migrator.HasTable(&commonModels.AlarmRule{}) {
		tables = append(tables, &commonModels.AlarmRule{})
	}
	if !migrator.HasTable(&commonModels.AlarmNotice{}) {
		tables = append(tables, &commonModels.AlarmNotice{})
	}
	if !migrator.HasTable(&commonModels.AlarmInstance{}) {
		tables = append(tables, &commonModels.AlarmInstance{})
	}
	if !migrator.HasTable(&commonModels.AlertRecord{}) {
		tables = append(tables, &commonModels.AlertRecord{})
	}
	if !migrator.HasTable(&commonModels.AlertContactGroupRel{}) {
		tables = append(tables, &commonModels.AlertContactGroupRel{})
	}
	if !migrator.HasTable(&commonModels.AlertContactInformation{}) {
		tables = append(tables, &commonModels.AlertContactInformation{})
	}
	if !migrator.HasTable(&commonModels.NotificationRecord{}) {
		tables = append(tables, &commonModels.NotificationRecord{})
	}

	return tables, sqls
}
