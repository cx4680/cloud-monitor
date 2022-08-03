

UPDATE t_monitor_item SET metrics_linux = '100 - (100 * (sum by(instance,instanceType) (irate(ecs_cpu_seconds_total{mode="idle",$INSTANCE}[3m])) / sum by(instance,instanceType) (irate(ecs_cpu_seconds_total{$INSTANCE}[3m]))))' WHERE biz_id = '1';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(ecs_load1{$INSTANCE})' WHERE biz_id = '2';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(ecs_load5{$INSTANCE})' WHERE biz_id = '3';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(ecs_load15{$INSTANCE})' WHERE biz_id = '4';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(ecs_memory_MemTotal_bytes{$INSTANCE}) - sum by(instance,instanceType)(ecs_memory_MemFree_bytes{$INSTANCE})' WHERE biz_id = '5';
UPDATE t_monitor_item SET metrics_linux = '100 * ((sum by(instance,instanceType)(ecs_memory_MemTotal_bytes{$INSTANCE}) - sum by(instance,instanceType)(ecs_memory_MemFree_bytes{$INSTANCE})) / sum by(instance,instanceType)(ecs_memory_MemTotal_bytes{$INSTANCE}))' WHERE biz_id = '6';
UPDATE t_monitor_item SET metrics_linux = '100 * ((sum by(instance,instanceType,device)(ecs_filesystem_size_bytes{$INSTANCE}) - sum by(instance,instanceType,device)(ecs_filesystem_free_bytes{$INSTANCE})) / sum by(instance,instanceType,device)(ecs_filesystem_size_bytes{$INSTANCE}))' WHERE biz_id = '7';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType,device)(irate(ecs_disk_read_bytes_total{$INSTANCE}[3m]))' WHERE biz_id = '8';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType,device)(irate(ecs_disk_written_bytes_total{$INSTANCE}[3m]))' WHERE biz_id = '9';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType,device)(irate(ecs_disk_reads_completed_total{$INSTANCE}[3m]))' WHERE biz_id = '10';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType,device)(irate(ecs_disk_writes_completed_total{$INSTANCE}[3m]))' WHERE biz_id = '11';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType,device)(irate(ecs_network_receive_bytes_total{$INSTANCE}[3m])) / 1024 / 1024 * 8' WHERE biz_id = '12';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType,device)(irate(ecs_network_transmit_bytes_total{$INSTANCE}[3m])) / 1024 / 1024 * 8' WHERE biz_id = '13';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType,device)(irate(ecs_network_receive_packets_total{$INSTANCE}[3m]))' WHERE biz_id = '14';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType,device)(irate(ecs_network_transmit_packets_total{$INSTANCE}[3m]))' WHERE biz_id = '15';
UPDATE t_monitor_item SET metrics_linux = '100 * avg by(instance,instanceType)(irate(ecs_base_vcpu_seconds{$INSTANCE}[6m]))' WHERE biz_id = '16';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType,drive)(irate(ecs_base_storage_traffic_bytes_total{type="read",$INSTANCE}[6m]))' WHERE biz_id = '17';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType,drive)(irate(ecs_base_storage_traffic_bytes_total{type="write",$INSTANCE}[6m]))' WHERE biz_id = '18';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType,interface)(irate(ecs_base_network_traffic_bytes_total{type="rx",$INSTANCE}[6m]))' WHERE biz_id = '19';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType,interface)(irate(ecs_base_network_traffic_bytes_total{type="tx",$INSTANCE}[6m]))' WHERE biz_id = '20';

UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(cbr_vault_size{$INSTANCE})' WHERE biz_id = '52';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(cbr_vault_used{$INSTANCE})' WHERE biz_id = '53';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(cbr_vault_used{$INSTANCE}) / sum by(instance,instanceType)(cbr_vault_size{$INSTANCE}) * 100' WHERE biz_id = '54';

UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType,device)(ecs_filesystem_free_bytes{$INSTANCE})' WHERE biz_id = '65';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType,device)(ecs_filesystem_size_bytes{$INSTANCE}) - sum by(instance,instanceType,device)(ecs_filesystem_free_bytes{$INSTANCE})' WHERE biz_id = '66';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType,device)(ecs_filesystem_size_bytes{$INSTANCE})' WHERE biz_id = '67';

UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(ecs_load1{$INSTANCE})' WHERE biz_id = '69';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(ecs_load5{$INSTANCE})' WHERE biz_id = '70';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(ecs_load15{$INSTANCE})' WHERE biz_id = '71';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(ecs_memory_MemTotal_bytes{$INSTANCE}) - sum by(instance,instanceType)(ecs_memory_MemFree_bytes{$INSTANCE})' WHERE biz_id = '72';
UPDATE t_monitor_item SET metrics_linux = '100 * ((sum by(instance,instanceType)(ecs_memory_MemTotal_bytes{$INSTANCE}) - sum by(instance,instanceType)(ecs_memory_MemFree_bytes{$INSTANCE})) / sum by(instance,instanceType)(ecs_memory_MemTotal_bytes{$INSTANCE}))' WHERE biz_id = '73';
UPDATE t_monitor_item SET metrics_linux = '100 * ((sum by(instance,instanceType,device)(ecs_filesystem_size_bytes{$INSTANCE}) - sum by(instance,instanceType,device)(ecs_filesystem_free_bytes{$INSTANCE})) / sum by(instance,instanceType,device)(ecs_filesystem_size_bytes{$INSTANCE}))' WHERE biz_id = '74';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType,device)(irate(ecs_disk_read_bytes_total{$INSTANCE}[3m]))' WHERE biz_id = '75';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType,device)(irate(ecs_disk_written_bytes_total{$INSTANCE}[3m]))' WHERE biz_id = '76';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType,device)(irate(ecs_disk_reads_completed_total{$INSTANCE}[3m]))' WHERE biz_id = '77';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType,device)(irate(ecs_disk_writes_completed_total{$INSTANCE}[3m]))' WHERE biz_id = '78';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType,device)(irate(ecs_network_receive_bytes_total{$INSTANCE}[3m])) / 1024 / 1024 * 8' WHERE biz_id = '79';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType,device)(irate(ecs_network_transmit_bytes_total{$INSTANCE}[3m])) / 1024 / 1024 * 8' WHERE biz_id = '80';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType,device)(irate(ecs_network_receive_packets_total{$INSTANCE}[3m]))' WHERE biz_id = '81';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType,device)(irate(ecs_network_transmit_packets_total{$INSTANCE}[3m]))' WHERE biz_id = '82';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType,device)(ecs_filesystem_free_bytes{$INSTANCE})' WHERE biz_id = '83';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType,device)(ecs_filesystem_size_bytes{$INSTANCE}) - sum by(instance,instanceType,device)(ecs_filesystem_free_bytes{$INSTANCE})' WHERE biz_id = '84';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType,device)(ecs_filesystem_size_bytes{$INSTANCE})' WHERE biz_id = '85';

UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(mysql_slave_io{$INSTANCE})' WHERE biz_id = '86';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(mysql_slave_sql{$INSTANCE})' WHERE biz_id = '87';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(mysql_slave_seconds_behind_master{$INSTANCE})' WHERE biz_id = '88';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(mysql_active_connections{$INSTANCE})' WHERE biz_id = '89';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(mysql_current_connection_percent{$INSTANCE})' WHERE biz_id = '90';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(mysql_qps{$INSTANCE})' WHERE biz_id = '91';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(mysql_tps{$INSTANCE})' WHERE biz_id = '92';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(mysql_select_ps{$INSTANCE})' WHERE biz_id = '93';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(mysql_update_ps{$INSTANCE})' WHERE biz_id = '94';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(mysql_insert_ps{$INSTANCE})' WHERE biz_id = '95';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(mysql_delete_ps{$INSTANCE})' WHERE biz_id = '96';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(mysql_cpu_usage{$INSTANCE})' WHERE biz_id = '97';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(mysql_mem_usage{$INSTANCE})' WHERE biz_id = '98';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(mysql_disk_usage{$INSTANCE})' WHERE biz_id = '99';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(mysql_innodb_select_ps{$INSTANCE})' WHERE biz_id = '100';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(mysql_innodb_update_ps{$INSTANCE})' WHERE biz_id = '101';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(mysql_innodb_insert_ps{$INSTANCE})' WHERE biz_id = '102';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(mysql_innodb_delete_ps{$INSTANCE})' WHERE biz_id = '103';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(mysql_innodb_cache_hit_rate{$INSTANCE})' WHERE biz_id = '104';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(mysql_innodb_reads_ps{$INSTANCE})' WHERE biz_id = '105';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(mysql_innodb_writes_ps{$INSTANCE})' WHERE biz_id = '106';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(mysql_innodb_buffer_pool_pages_dirty{$INSTANCE})' WHERE biz_id = '107';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(mysql_innodb_buffer_pool_bytes_dirty{$INSTANCE})' WHERE biz_id = '108';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(mysql_innodb_log_waits{$INSTANCE})' WHERE biz_id = '109';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(mysql_binlog_cache_disk_use{$INSTANCE})' WHERE biz_id = '110';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(mysql_slow_queries_per_min{$INSTANCE})' WHERE biz_id = '111';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(mysql_long_query_count{$INSTANCE})' WHERE biz_id = '112';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(mysql_long_query_alert_count{$INSTANCE})' WHERE biz_id = '113';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(mysql_exec_statememt_frequency{$INSTANCE})' WHERE biz_id = '114';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(mysql_exec_statememt_frequency{$INSTANCE})' WHERE biz_id = '115';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(mysql_write_frequency{$INSTANCE})' WHERE biz_id = '116';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(mysql_top_statememt_avg_exec_time{$INSTANCE})' WHERE biz_id = '117';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(mysql_top_statememt_exec_err_rate{$INSTANCE})' WHERE biz_id = '118';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(mysql_current_cons_num{$INSTANCE})' WHERE biz_id = '119';

UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(rate(dm_global_status_tps{$INSTANCE}[1m]))' WHERE biz_id = '120';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(rate(dm_global_status_qps{$INSTANCE}[1m]))' WHERE biz_id = '121';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(rate(dm_global_status_ips{$INSTANCE}[1m]))' WHERE biz_id = '122';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(rate(dm_global_status_dps{$INSTANCE}[1m]))' WHERE biz_id = '123';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(rate(dm_global_status_ups{$INSTANCE}[1m]))' WHERE biz_id = '124';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(rate(dm_global_status_ddlps{$INSTANCE}[1m]))' WHERE biz_id = '125';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(rate(dm_global_status_nioips{$INSTANCE}[1m]))' WHERE biz_id = '126';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(rate(dm_global_status_nio_ops{$INSTANCE}[1m]))' WHERE biz_id = '127';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(rate(dm_global_status_fio_ips{$INSTANCE}[1m]))' WHERE biz_id = '128';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(rate(dm_global_status_fio_ops{$INSTANCE}[1m]))' WHERE biz_id = '129';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(dm_global_status_mem_used{$INSTANCE})' WHERE biz_id = '130';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(dm_global_status_cpu_use_rate{$INSTANCE})' WHERE biz_id = '131';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(dm_global_status_sessions{$INSTANCE})' WHERE biz_id = '132';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(dm_global_status_active_sessions{$INSTANCE})' WHERE biz_id = '133';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(dm_global_status_task_waiting{$INSTANCE})' WHERE biz_id = '134';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(dm_global_status_task_ready{$INSTANCE})' WHERE biz_id = '135';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(dm_global_status_task_total_wait_time{$INSTANCE})' WHERE biz_id = '136';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(dm_global_status_avg_wait_time{$INSTANCE})' WHERE biz_id = '137';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(dm_global_status_threads{$INSTANCE})' WHERE biz_id = '138';

UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(pg_cpu_usage{$INSTANCE})' WHERE biz_id = '139';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(pg_mem_usage{$INSTANCE})' WHERE biz_id = '140';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(pg_disk_usage{$INSTANCE})' WHERE biz_id = '141';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(pg_qps{$INSTANCE})' WHERE biz_id = '142';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(pg_rqps{$INSTANCE})' WHERE biz_id = '143';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(pg_wqps{$INSTANCE})' WHERE biz_id = '144';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(pg_tps{$INSTANCE})' WHERE biz_id = '145';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(pg_mean_exec_time{$INSTANCE})' WHERE biz_id = '146';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(pg_open_ct_num{$INSTANCE})' WHERE biz_id = '147';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(pg_active_ct_num{$INSTANCE})' WHERE biz_id = '148';

UPDATE t_monitor_item SET metrics_linux = 'sum by (instance,instanceType)(kafka_controller_kafkacontroller_activecontrollercount{$INSTANCE})' WHERE biz_id = '149';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance,instanceType)(kafka_brokers{$INSTANCE})' WHERE biz_id = '150';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance,instanceType)(kafka_controller_kafkacontroller_globalpartitioncount{$INSTANCE})' WHERE biz_id = '151';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance,instanceType)(kafka_controller_kafkacontroller_globaltopiccount{$INSTANCE})' WHERE biz_id = '152';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance,instanceType)(kafka_controller_kafkacontroller_offlinepartitionscount{$INSTANCE})' WHERE biz_id = '153';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance,instanceType)(kafka_controller_kafkacontroller_preferredreplicaimbalancecount{$INSTANCE})' WHERE biz_id = '154';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance,instanceType)(kafka_controller_kafkacontroller_topicstodeletecount{$INSTANCE})' WHERE biz_id = '155';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance,instanceType)(kafka_controller_controllerstats_uncleanleaderelectionspersec{$INSTANCE})' WHERE biz_id = '156';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance,instanceType)(kafka_server_replicamanager_leadercount{$INSTANCE})' WHERE biz_id = '157';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance,instanceType)(kafka_server_replicamanager_partitioncount{$INSTANCE})' WHERE biz_id = '158';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance,instanceType)(kafka_server_replicamanager_underminisrpartitioncount{$INSTANCE})' WHERE biz_id = '159';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance,instanceType)(kafka_server_replicamanager_underreplicatedpartitions{$INSTANCE})' WHERE biz_id = '160';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance,instanceType)(kafka_server_replicamanager_reassigningpartitions{$INSTANCE})' WHERE biz_id = '161';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance,instanceType)(kafka_server_brokertopicmetrics_bytesinpersec{$INSTANCE})' WHERE biz_id = '162';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance,instanceType)(kafka_server_brokertopicmetrics_bytesoutpersec{$INSTANCE})' WHERE biz_id = '163';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance,instanceType)(kafka_server_brokertopicmetrics_messagesinpersec{$INSTANCE})' WHERE biz_id = '164';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance,instanceType)(kafka_consumergroup_lag{$INSTANCE})' WHERE biz_id = '165';

UPDATE t_monitor_item SET metrics_linux = 'sum by (instance,instanceType)(dm_global_status_mem_use_rate{$INSTANCE})' WHERE biz_id = '166';

UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(ecs_load1{$INSTANCE})' WHERE biz_id = '168';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(ecs_load5{$INSTANCE})' WHERE biz_id = '169';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(ecs_load15{$INSTANCE})' WHERE biz_id = '170';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(ecs_memory_MemTotal_bytes{$INSTANCE}) - sum by(instance,instanceType)(ecs_memory_MemFree_bytes{$INSTANCE})' WHERE biz_id = '171';
UPDATE t_monitor_item SET metrics_linux = '100 * ((sum by(instance,instanceType)(ecs_memory_MemTotal_bytes{$INSTANCE}) - sum by(instance,instanceType)(ecs_memory_MemFree_bytes{$INSTANCE})) / sum by(instance,instanceType)(ecs_memory_MemTotal_bytes{$INSTANCE}))' WHERE biz_id = '172';
UPDATE t_monitor_item SET metrics_linux = '100 * ((sum by(instance,instanceType,device)(ecs_filesystem_size_bytes{$INSTANCE}) - sum by(instance,instanceType,device)(ecs_filesystem_free_bytes{$INSTANCE})) / sum by(instance,instanceType,device)(ecs_filesystem_size_bytes{$INSTANCE}))' WHERE biz_id = '173';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType,device)(irate(ecs_disk_read_bytes_total{$INSTANCE}[3m]))' WHERE biz_id = '174';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType,device)(irate(ecs_disk_written_bytes_total{$INSTANCE}[3m]))' WHERE biz_id = '175';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType,device)(irate(ecs_disk_reads_completed_total{$INSTANCE}[3m]))' WHERE biz_id = '176';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType,device)(irate(ecs_disk_writes_completed_total{$INSTANCE}[3m]))' WHERE biz_id = '177';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType,device)(irate(ecs_network_receive_bytes_total{$INSTANCE}[3m])) / 1024 / 1024 * 8' WHERE biz_id = '178';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType,device)(irate(ecs_network_transmit_bytes_total{$INSTANCE}[3m])) / 1024 / 1024 * 8' WHERE biz_id = '179';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType,device)(irate(ecs_network_receive_packets_total{$INSTANCE}[3m]))' WHERE biz_id = '180';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType,device)(irate(ecs_network_transmit_packets_total{$INSTANCE}[3m]))' WHERE biz_id = '181';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType,device)(ecs_filesystem_free_bytes{$INSTANCE})' WHERE biz_id = '182';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType,device)(ecs_filesystem_size_bytes{$INSTANCE}) - sum by(instance,instanceType,device)(ecs_filesystem_free_bytes{$INSTANCE})' WHERE biz_id = '183';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType,device)(ecs_filesystem_size_bytes{$INSTANCE})' WHERE biz_id = '184';

UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(redis_cpu_usage{$INSTANCE})' WHERE biz_id = '185';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(redis_mem_usage{$INSTANCE})' WHERE biz_id = '186';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(redis_connected_clients{$INSTANCE})' WHERE biz_id = '187';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(redis_tps{$INSTANCE})' WHERE biz_id = '187';

UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType,pod)(mongo_mongos_current_connections{$INSTANCE})' WHERE biz_id = '195';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType,pod)(mongo_shard_current_connections{$INSTANCE})' WHERE biz_id = '196';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType,pod)(mongo_config_current_connections{$INSTANCE})' WHERE biz_id = '197';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(mongo_total_current_connections{$INSTANCE})' WHERE biz_id = '198';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType,pod)(mongo_mongos_memory_ratio{$INSTANCE})' WHERE biz_id = '199';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(mongo_config_memory_ratio{$INSTANCE})' WHERE biz_id = '200';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType,pod)(mongo_shard_memory_ratio{$INSTANCE})' WHERE biz_id = '201';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType,pod)(mongo_mongos_cpu_ratio{$INSTANCE})' WHERE biz_id = '202';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType,pod)(mongo_shard_cpu_ratio{$INSTANCE})' WHERE biz_id = '203';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType)(mongo_config_cpu_ratio{$INSTANCE})' WHERE biz_id = '204';


UPDATE t_monitor_item SET metrics_linux = 'sum by (instance,instanceType)(mysql_slave_io{$INSTANCE})' WHERE biz_id = '86';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance,instanceType)(mysql_slave_sql{$INSTANCE})' WHERE biz_id = '87';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance,instanceType)(mysql_slave_seconds_behind_master{$INSTANCE})' WHERE biz_id = '88';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance,instanceType)(mysql_active_connections{$INSTANCE})' WHERE biz_id = '89';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance,instanceType)(mysql_current_connection_percent{$INSTANCE})' WHERE biz_id = '90';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance,instanceType)(mysql_qps{$INSTANCE})' WHERE biz_id = '91';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance,instanceType)(mysql_tps{$INSTANCE})' WHERE biz_id = '92';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance,instanceType)(mysql_select_ps{$INSTANCE})' WHERE biz_id = '93';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance,instanceType)(mysql_update_ps{$INSTANCE})' WHERE biz_id = '94';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance,instanceType)(mysql_insert_ps{$INSTANCE})' WHERE biz_id = '95';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance,instanceType)(mysql_delete_ps{$INSTANCE})' WHERE biz_id = '96';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance,instanceType)(mysql_cpu_usage{$INSTANCE})' WHERE biz_id = '97';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance,instanceType)(mysql_mem_usage{$INSTANCE})' WHERE biz_id = '98';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance,instanceType)(mysql_disk_usage{$INSTANCE})' WHERE biz_id = '99';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance,instanceType)(mysql_innodb_select_ps{$INSTANCE})' WHERE biz_id = '100';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance,instanceType)(mysql_innodb_update_ps{$INSTANCE})' WHERE biz_id = '101';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance,instanceType)(mysql_innodb_insert_ps{$INSTANCE})' WHERE biz_id = '102';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance,instanceType)(mysql_innodb_delete_ps{$INSTANCE})' WHERE biz_id = '103';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance,instanceType)(mysql_innodb_cache_hit_rate{$INSTANCE})' WHERE biz_id = '104';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance,instanceType)(mysql_innodb_reads_ps{$INSTANCE})' WHERE biz_id = '105';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance,instanceType)(mysql_innodb_writes_ps{$INSTANCE})' WHERE biz_id = '106';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance,instanceType)(mysql_innodb_buffer_pool_pages_dirty{$INSTANCE})' WHERE biz_id = '107';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance,instanceType)(mysql_innodb_buffer_pool_bytes_dirty{$INSTANCE})' WHERE biz_id = '108';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance,instanceType)(mysql_innodb_log_waits{$INSTANCE})' WHERE biz_id = '109';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance,instanceType)(mysql_binlog_cache_disk_use{$INSTANCE})' WHERE biz_id = '110';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance,instanceType)(mysql_slow_queries_per_min{$INSTANCE})' WHERE biz_id = '111';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance,instanceType)(mysql_long_query_count{$INSTANCE})' WHERE biz_id = '112';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance,instanceType)(mysql_long_query_alert_count{$INSTANCE})' WHERE biz_id = '113';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance,instanceType)(mysql_exec_statememt_frequency{$INSTANCE})' WHERE biz_id = '114';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance,instanceType)(mysql_read_frequency{$INSTANCE})' WHERE biz_id = '115';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance,instanceType)(mysql_write_frequency{$INSTANCE})' WHERE biz_id = '116';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance,instanceType)(mysql_top_statememt_avg_exec_time{$INSTANCE})' WHERE biz_id = '117';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance,instanceType)(mysql_top_statememt_exec_err_rate{$INSTANCE})' WHERE biz_id = '118';
UPDATE t_monitor_item SET metrics_linux = 'sum by (instance,instanceType)(mysql_current_cons_num{$INSTANCE})' WHERE biz_id = '119';

UPDATE t_monitor_item SET metrics_linux = 'avg by(instance,instanceType)(dm_global_status_mem_used{$INSTANCE})' WHERE biz_id = '130';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType,service,route)(guard_http_apirequests{$INSTANCE})' WHERE biz_id = '218';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance,instanceType,service)(guard_http_apirequests{$INSTANCE})' WHERE biz_id = '219';