DELETE FROM t_monitor_item WHERE biz_id IN ('149','150','151','152','153','154','220','221','222','223','224','225','226','227','228');

DELETE FROM t_config_item WHERE biz_id = '24' ;

UPDATE t_monitor_product SET status = '0' WHERE abbreviation IN ('cgw');

INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`, `display`) VALUES ('149', '11', 'Broker控制器', 'kafka_controller_kafkacontroller_activecontrollercount', 'instance', 'sum by (instance)(kafka_controller_kafkacontroller_activecontrollercount{$INSTANCE})', null, null, '个', null, '2', '1', '1', null, null, null, null, 'chart,rule');
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`, `display`) VALUES ('150', '11', '在线Broker数', 'kafka_brokers', 'instance', 'sum by (instance)(kafka_brokers{$INSTANCE})', null, null, '个', null, '2', '1', '1', null, null, null, null, 'chart,rule');
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`, `display`) VALUES ('151', '11', '全局分区数', 'kafka_controller_kafkacontroller_globalpartitioncount', 'instance', 'sum by (instance)(kafka_controller_kafkacontroller_globalpartitioncount{$INSTANCE})', null, null, '个', null, '2', '1', '1', null, null, null, null, 'chart,rule');
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`, `display`) VALUES ('152', '11', '全局主题数', 'kafka_controller_kafkacontroller_globaltopiccount', 'instance', 'sum by (instance)(kafka_controller_kafkacontroller_globaltopiccount{$INSTANCE})', null, null, '个', null, '2', '1', '1', null, null, null, null, 'chart,rule');
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`, `display`) VALUES ('153', '11', '离线主题数', 'kafka_controller_kafkacontroller_offlinepartitionscount', 'instance', 'sum by (instance)(kafka_controller_kafkacontroller_offlinepartitionscount{$INSTANCE})', null, null, '个', null, '2', '1', '1', null, null, null, null, 'chart,rule');
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`, `display`) VALUES ('154', '11', '未平衡的副本数', 'kafka_controller_kafkacontroller_preferredreplicaimbalancecount', 'instance', 'sum by (instance)(kafka_controller_kafkacontroller_offlinepartitionscount{$INSTANCE})', null, null, '个', null, '2', '1', '1', null, null, null, null, 'chart,rule');
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`, `display`) VALUES ('155', '11', '删除的主题数', 'kafka_controller_kafkacontroller_topicstodeletecount', 'instance', 'sum by (instance)(kafka_controller_kafkacontroller_topicstodeletecount{$INSTANCE})', null, null, '个', null, '2', '1', '1', null, null, null, null, 'chart,rule');
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`, `display`) VALUES ('156', '11', '未同步的选举', 'kafka_controller_controllerstats_uncleanleaderelectionspersec', 'instance', 'sum by (instance)(kafka_controller_controllerstats_uncleanleaderelectionspersec{$INSTANCE})', null, null, '次/s', null, '2', '1', '1', null, null, null, null, 'chart,rule');
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`, `display`) VALUES ('157', '11', 'Leader数量', 'kafka_server_replicamanager_leadercount', 'instance', 'sum by (instance)(kafka_server_replicamanager_leadercount{$INSTANCE})', null, null, '个', null, '2', '1', '1', null, null, null, null, 'chart,rule');
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`, `display`) VALUES ('158', '11', '分区的数量', 'kafka_server_replicamanager_partitioncount', 'instance', 'sum by (instance)(kafka_server_replicamanager_partitioncount{$INSTANCE})', null, null, '个', null, '2', '1', '1', null, null, null, null, 'chart,rule');
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`, `display`) VALUES ('159', '11', '最小未同步分区的数量', 'kafka_server_replicamanager_underminisrpartitioncount', 'instance', 'sum by (instance)(kafka_server_replicamanager_underminisrpartitioncount{$INSTANCE})', null, null, '个', null, '2', '1', '1', null, null, null, null, 'chart,rule');
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`, `display`) VALUES ('160', '11', '没有副本的分区数', 'kafka_server_replicamanager_underreplicatedpartitions', 'instance', 'sum by (instance)(kafka_server_replicamanager_underreplicatedpartitions{$INSTANCE})', null, null, '个', null, '2', '1', '1', null, null, null, null, 'chart,rule');
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`, `display`) VALUES ('161', '11', '再分配的分区数', 'kafka_server_replicamanager_reassigningpartitions', 'instance', 'sum by (instance)(kafka_server_replicamanager_reassigningpartitions{$INSTANCE})', null, null, '个', null, '2', '1', '1', null, null, null, null, 'chart,rule');
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`, `display`) VALUES ('162', '11', '生产速率', 'kafka_server_brokertopicmetrics_bytesinpersec', 'instance', 'sum by (instance)(kafka_server_brokertopicmetrics_bytesinpersec{$INSTANCE})', null, null, 'Byte/s', null, '2', '1', '1', null, null, null, null, 'chart,rule');
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`, `display`) VALUES ('163', '11', '消费速率', 'kafka_server_brokertopicmetrics_bytesoutpersec', 'instance', 'sum by (instance)(kafka_server_brokertopicmetrics_bytesoutpersec{$INSTANCE})', null, null, 'Byte/s', null, '2', '1', '1', null, null, null, null, 'chart,rule');
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`, `display`) VALUES ('164', '11', '消息生产速率', 'kafka_server_brokertopicmetrics_messagesinpersec', 'instance', 'sum by (instance)(kafka_server_brokertopicmetrics_messagesinpersec{$INSTANCE})', null, null, '个/s', null, '2', '1', '1', null, null, null, null, 'chart,rule');
INSERT INTO `t_monitor_item` (`biz_id`, `product_biz_id`, `name`, `metric_name`, `labels`, `metrics_linux`, `metrics_windows`, `statistics`, `unit`, `frequency`, `type`, `is_display`, `status`, `description`, `create_user`, `create_time`, `show_expression`, `display`) VALUES ('165', '11', '落后的消费量', 'kafka_consumergroup_lag', 'instance', 'sum by (instance)(kafka_consumergroup_lag{$INSTANCE})', null, null, '个', null, '2', '1', '1', null, null, null, null, 'chart,rule');