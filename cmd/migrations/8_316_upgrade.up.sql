UPDATE t_monitor_item SET metrics_linux='100 - (100 * (sum by(instance,instanceType) (irate(ecs_cpu_seconds_total{mode="idle",$INSTANCE}[3m])) / sum by(instance,instanceType) (irate(ecs_cpu_seconds_total{$INSTANCE}[3m]))))' WHERE biz_id='1';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType)(ecs_load1{$INSTANCE})' WHERE biz_id='2';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType)(ecs_load5{$INSTANCE})' WHERE biz_id='3';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType)(ecs_load15{$INSTANCE})' WHERE biz_id='4';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType)(ecs_memory_MemTotal_bytes{$INSTANCE}) - sum by(instance,instanceType)(ecs_memory_MemFree_bytes{$INSTANCE})' WHERE biz_id='5';
UPDATE t_monitor_item SET metrics_linux='100 * ((sum by(instance,instanceType)(ecs_memory_MemTotal_bytes{$INSTANCE}) - sum by(instance,instanceType)(ecs_memory_MemFree_bytes{$INSTANCE})) / sum by(instance,instanceType)(ecs_memory_MemTotal_bytes{$INSTANCE}))' WHERE biz_id='6';
UPDATE t_monitor_item SET metrics_linux='100 * ((sum by(instance,instanceType,device)(ecs_filesystem_size_bytes{$INSTANCE}) - sum by(instance,instanceType,device)(ecs_filesystem_free_bytes{$INSTANCE})) / sum by(instance,instanceType,device)(ecs_filesystem_size_bytes{$INSTANCE}))' WHERE biz_id='7';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,device)(irate(ecs_disk_read_bytes_total{$INSTANCE}[3m]))' WHERE biz_id='8';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,device)(irate(ecs_disk_written_bytes_total{$INSTANCE}[3m]))' WHERE biz_id='9';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,device)(irate(ecs_disk_reads_completed_total{$INSTANCE}[3m]))' WHERE biz_id='10';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,device)(irate(ecs_disk_writes_completed_total{$INSTANCE}[3m]))' WHERE biz_id='11';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,device)(irate(ecs_network_receive_bytes_total{$INSTANCE}[3m])) / 1024 / 1024 * 8' WHERE biz_id='12';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,device)(irate(ecs_network_transmit_bytes_total{$INSTANCE}[3m])) / 1024 / 1024 * 8' WHERE biz_id='13';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,device)(irate(ecs_network_receive_packets_total{$INSTANCE}[3m]))' WHERE biz_id='14';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,device)(irate(ecs_network_transmit_packets_total{$INSTANCE}[3m]))' WHERE biz_id='15';
UPDATE t_monitor_item SET metrics_linux='100 * avg by(instance,instanceType)(irate(ecs_base_vcpu_seconds{$INSTANCE}[6m]))' WHERE biz_id='16';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,drive)(irate(ecs_base_storage_traffic_bytes_total{type="read",$INSTANCE}[6m]))' WHERE biz_id='17';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,drive)(irate(ecs_base_storage_traffic_bytes_total{type="write",$INSTANCE}[6m]))' WHERE biz_id='18';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,interface)(irate(ecs_base_network_traffic_bytes_total{type="rx",$INSTANCE}[6m]))' WHERE biz_id='19';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,interface)(irate(ecs_base_network_traffic_bytes_total{type="tx",$INSTANCE}[6m]))' WHERE biz_id='20';
UPDATE t_monitor_item SET metrics_linux='sum(eip_upstream_bits_rate{$INSTANCE}) by (instance,instanceType,eip)' WHERE biz_id='21';
UPDATE t_monitor_item SET metrics_linux='sum(eip_downstream_bits_rate{$INSTANCE}) by (instance,instanceType,eip)' WHERE biz_id='22';
UPDATE t_monitor_item SET metrics_linux='((sum(eip_upstream_bits_rate{$INSTANCE}) by (instance,instanceType,eip))/8)*60' WHERE biz_id='23';
UPDATE t_monitor_item SET metrics_linux='((sum(eip_downstream_bits_rate{$INSTANCE}) by (instance,instanceType,eip))/8)*60' WHERE biz_id='24';
UPDATE t_monitor_item SET metrics_linux='(sum(eip_upstream_bits_rate{$INSTANCE}) by (instance,instanceType,eip) / avg(eip_config_upstream_bandwidth{$INSTANCE}) by (instance,instanceType,eip)) * 100' WHERE biz_id='25';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType) (Slb_http_bps_out_rate{$INSTANCE})' WHERE biz_id='26';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType) (Slb_http_bps_in_rate{$INSTANCE})' WHERE biz_id='27';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType) (Slb_all_connection_count{$INSTANCE})' WHERE biz_id='30';
UPDATE t_monitor_item SET metrics_linux='sum by (instance,instanceType)(Slb_all_est_connection_count{$INSTANCE})' WHERE biz_id='31';
UPDATE t_monitor_item SET metrics_linux='sum by (instance,instanceType) (Slb_all_none_est_connection_count{$INSTANCE})' WHERE biz_id='32';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType) (Slb_new_connection_rate{$INSTANCE})' WHERE biz_id='33';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType)(Slb_drop_connection_rate{$INSTANCE})' WHERE biz_id='34';
UPDATE t_monitor_item SET metrics_linux='avg by(instance,instanceType) (Slb_unhealthy_server_count{$INSTANCE})' WHERE biz_id='39';
UPDATE t_monitor_item SET metrics_linux='avg by(instance,instanceType) (Slb_healthy_server_count{$INSTANCE})' WHERE biz_id='40';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType)(Slb_request_rate{$INSTANCE})' WHERE biz_id='41';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType) (Slb_http_2xx_rate{$INSTANCE})' WHERE biz_id='42';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType) (Slb_http_3xx_rate{$INSTANCE})' WHERE biz_id='43';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType) (Slb_http_4xx_rate{$INSTANCE})' WHERE biz_id='44';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType) (Slb_http_5xx_rate{$INSTANCE})' WHERE biz_id='45';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType)(cbr_vault_size{$INSTANCE})' WHERE biz_id='52';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType)(cbr_vault_used{$INSTANCE})' WHERE biz_id='53';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType)(cbr_vault_used{$INSTANCE}) / sum by(instance,instanceType)(cbr_vault_size{$INSTANCE}) * 100' WHERE biz_id='54';
UPDATE t_monitor_item SET metrics_linux='sum by (instance,instanceType)(Nat_snat_total_connection_count{$INSTANCE})' WHERE biz_id='55';
UPDATE t_monitor_item SET metrics_linux='sum by (instance,instanceType)(rate(Nat_recv_bytes_total_count{$INSTANCE}[1m])*8)' WHERE biz_id='56';
UPDATE t_monitor_item SET metrics_linux='sum by (instance,instanceType)(rate(Nat_send_bytes_total_count{$INSTANCE}[1m])*8)' WHERE biz_id='57';
UPDATE t_monitor_item SET metrics_linux='sum by (instance,instanceType)(Nat_recv_bytes_total_count{$INSTANCE})' WHERE biz_id='58';
UPDATE t_monitor_item SET metrics_linux='sum by (instance,instanceType)(Nat_send_bytes_total_count{$INSTANCE})' WHERE biz_id='59';
UPDATE t_monitor_item SET metrics_linux='sum by (instance,instanceType)(rate(Nat_recv_packets_total_count{$INSTANCE}[1m]))' WHERE biz_id='60';
UPDATE t_monitor_item SET metrics_linux='sum by (instance,instanceType)(rate(Nat_send_packets_total_count{$INSTANCE}[1m]))' WHERE biz_id='61';
UPDATE t_monitor_item SET metrics_linux='sum by (instance,instanceType)(Nat_snat_total_connection_count{$INSTANCE}) / avg by (instance,instanceType)(Nat_nat_max_connection_count{$INSTANCE}) *100' WHERE biz_id='62';
UPDATE t_monitor_item SET metrics_linux='sum(irate(ecs_base_storage_iops_total{type="read",$INSTANCE}[15m])) by (instance,instanceType,drive)' WHERE biz_id='63';
UPDATE t_monitor_item SET metrics_linux='sum(irate(ecs_base_storage_iops_total{type="write",$INSTANCE}[15m])) by (instance,instanceType,drive)' WHERE biz_id='64';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,device)(ecs_filesystem_free_bytes{$INSTANCE})' WHERE biz_id='65';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,device)(ecs_filesystem_size_bytes{$INSTANCE}) - sum by(instance,instanceType,device)(ecs_filesystem_free_bytes{$INSTANCE})' WHERE biz_id='66';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,device)(ecs_filesystem_size_bytes{$INSTANCE})' WHERE biz_id='67';
UPDATE t_monitor_item SET metrics_linux='100 - (100 * (sum by(instance,instanceType) (irate(ecs_cpu_seconds_total{mode="idle",$INSTANCE}[3m])) / sum by(instance,instanceType) (irate(ecs_cpu_seconds_total{$INSTANCE}[3m]))))' WHERE biz_id='68';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType)(ecs_load1{$INSTANCE})' WHERE biz_id='69';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType)(ecs_load5{$INSTANCE})' WHERE biz_id='70';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType)(ecs_load15{$INSTANCE})' WHERE biz_id='71';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType)(ecs_memory_MemTotal_bytes{$INSTANCE}) - sum by(instance,instanceType)(ecs_memory_MemFree_bytes{$INSTANCE})' WHERE biz_id='72';
UPDATE t_monitor_item SET metrics_linux='100 * ((sum by(instance,instanceType)(ecs_memory_MemTotal_bytes{$INSTANCE}) - sum by(instance,instanceType)(ecs_memory_MemFree_bytes{$INSTANCE})) / sum by(instance,instanceType)(ecs_memory_MemTotal_bytes{$INSTANCE}))' WHERE biz_id='73';
UPDATE t_monitor_item SET metrics_linux='100 * ((sum by(instance,instanceType,device)(ecs_filesystem_size_bytes{$INSTANCE}) - sum by(instance,instanceType,device)(ecs_filesystem_free_bytes{$INSTANCE})) / sum by(instance,instanceType,device)(ecs_filesystem_size_bytes{$INSTANCE}))' WHERE biz_id='74';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,device)(irate(ecs_disk_read_bytes_total{$INSTANCE}[3m]))' WHERE biz_id='75';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,device)(irate(ecs_disk_written_bytes_total{$INSTANCE}[3m]))' WHERE biz_id='76';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,device)(irate(ecs_disk_reads_completed_total{$INSTANCE}[3m]))' WHERE biz_id='77';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,device)(irate(ecs_disk_writes_completed_total{$INSTANCE}[3m]))' WHERE biz_id='78';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,device)(irate(ecs_network_receive_bytes_total{$INSTANCE}[3m])) / 1024 / 1024 * 8' WHERE biz_id='79';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,device)(irate(ecs_network_transmit_bytes_total{$INSTANCE}[3m])) / 1024 / 1024 * 8' WHERE biz_id='80';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,device)(irate(ecs_network_receive_packets_total{$INSTANCE}[3m]))' WHERE biz_id='81';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,device)(irate(ecs_network_transmit_packets_total{$INSTANCE}[3m]))' WHERE biz_id='82';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,device)(ecs_filesystem_free_bytes{$INSTANCE})' WHERE biz_id='83';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,device)(ecs_filesystem_size_bytes{$INSTANCE}) - sum by(instance,instanceType,device)(ecs_filesystem_free_bytes{$INSTANCE})' WHERE biz_id='84';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,device)(ecs_filesystem_size_bytes{$INSTANCE})' WHERE biz_id='85';
UPDATE t_monitor_item SET metrics_linux='sum by (instance,instanceType)(mysql_slave_io{$INSTANCE})' WHERE biz_id='86';
UPDATE t_monitor_item SET metrics_linux='sum by (instance,instanceType)(mysql_slave_sql{$INSTANCE})' WHERE biz_id='87';
UPDATE t_monitor_item SET metrics_linux='sum by (instance,instanceType)(mysql_slave_seconds_behind_master{$INSTANCE})' WHERE biz_id='88';
UPDATE t_monitor_item SET metrics_linux='sum by (instance,instanceType)(mysql_active_connections{$INSTANCE})' WHERE biz_id='89';
UPDATE t_monitor_item SET metrics_linux='sum by (instance,instanceType)(mysql_current_connection_percent{$INSTANCE})' WHERE biz_id='90';
UPDATE t_monitor_item SET metrics_linux='sum by (instance,instanceType)(mysql_qps{$INSTANCE})' WHERE biz_id='91';
UPDATE t_monitor_item SET metrics_linux='sum by (instance,instanceType)(mysql_tps{$INSTANCE})' WHERE biz_id='92';
UPDATE t_monitor_item SET metrics_linux='sum by (instance,instanceType)(mysql_select_ps{$INSTANCE})' WHERE biz_id='93';
UPDATE t_monitor_item SET metrics_linux='sum by (instance,instanceType)(mysql_update_ps{$INSTANCE})' WHERE biz_id='94';
UPDATE t_monitor_item SET metrics_linux='sum by (instance,instanceType)(mysql_insert_ps{$INSTANCE})' WHERE biz_id='95';
UPDATE t_monitor_item SET metrics_linux='sum by (instance,instanceType)(mysql_delete_ps{$INSTANCE})' WHERE biz_id='96';
UPDATE t_monitor_item SET metrics_linux='sum by (instance,instanceType)(mysql_cpu_usage{$INSTANCE})' WHERE biz_id='97';
UPDATE t_monitor_item SET metrics_linux='sum by (instance,instanceType)(mysql_mem_usage{$INSTANCE})' WHERE biz_id='98';
UPDATE t_monitor_item SET metrics_linux='sum by (instance,instanceType)(mysql_disk_usage{$INSTANCE})' WHERE biz_id='99';
UPDATE t_monitor_item SET metrics_linux='sum by (instance,instanceType)(mysql_innodb_select_ps{$INSTANCE})' WHERE biz_id='100';
UPDATE t_monitor_item SET metrics_linux='sum by (instance,instanceType)(mysql_innodb_update_ps{$INSTANCE})' WHERE biz_id='101';
UPDATE t_monitor_item SET metrics_linux='sum by (instance,instanceType)(mysql_innodb_insert_ps{$INSTANCE})' WHERE biz_id='102';
UPDATE t_monitor_item SET metrics_linux='sum by (instance,instanceType)(mysql_innodb_delete_ps{$INSTANCE})' WHERE biz_id='103';
UPDATE t_monitor_item SET metrics_linux='sum by (instance,instanceType)(mysql_innodb_cache_hit_rate{$INSTANCE})' WHERE biz_id='104';
UPDATE t_monitor_item SET metrics_linux='sum by (instance,instanceType)(mysql_innodb_reads_ps{$INSTANCE})' WHERE biz_id='105';
UPDATE t_monitor_item SET metrics_linux='sum by (instance,instanceType)(mysql_innodb_writes_ps{$INSTANCE})' WHERE biz_id='106';
UPDATE t_monitor_item SET metrics_linux='sum by (instance,instanceType)(mysql_innodb_buffer_pool_pages_dirty{$INSTANCE})' WHERE biz_id='107';
UPDATE t_monitor_item SET metrics_linux='sum by (instance,instanceType)(mysql_innodb_buffer_pool_bytes_dirty{$INSTANCE})' WHERE biz_id='108';
UPDATE t_monitor_item SET metrics_linux='sum by (instance,instanceType)(mysql_innodb_log_waits{$INSTANCE})' WHERE biz_id='109';
UPDATE t_monitor_item SET metrics_linux='sum by (instance,instanceType)(mysql_binlog_cache_disk_use{$INSTANCE})' WHERE biz_id='110';
UPDATE t_monitor_item SET metrics_linux='sum by (instance,instanceType)(mysql_slow_queries_per_min{$INSTANCE})' WHERE biz_id='111';
UPDATE t_monitor_item SET metrics_linux='sum by (instance,instanceType)(mysql_long_query_count{$INSTANCE})' WHERE biz_id='112';
UPDATE t_monitor_item SET metrics_linux='sum by (instance,instanceType)(mysql_long_query_alert_count{$INSTANCE})' WHERE biz_id='113';
UPDATE t_monitor_item SET metrics_linux='sum by (instance,instanceType)(mysql_exec_statememt_frequency{$INSTANCE})' WHERE biz_id='114';
UPDATE t_monitor_item SET metrics_linux='sum by (instance,instanceType)(mysql_read_frequency{$INSTANCE})' WHERE biz_id='115';
UPDATE t_monitor_item SET metrics_linux='sum by (instance,instanceType)(mysql_write_frequency{$INSTANCE})' WHERE biz_id='116';
UPDATE t_monitor_item SET metrics_linux='sum by (instance,instanceType)(mysql_top_statememt_avg_exec_time{$INSTANCE})' WHERE biz_id='117';
UPDATE t_monitor_item SET metrics_linux='sum by (instance,instanceType)(mysql_top_statememt_exec_err_rate{$INSTANCE})' WHERE biz_id='118';
UPDATE t_monitor_item SET metrics_linux='sum by (instance,instanceType)(mysql_current_cons_num{$INSTANCE})' WHERE biz_id='119';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType)(rate(dm_global_status_tps{$INSTANCE}[1m]))' WHERE biz_id='120';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType)(rate(dm_global_status_qps{$INSTANCE}[1m]))' WHERE biz_id='121';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType)(rate(dm_global_status_ips{$INSTANCE}[1m]))' WHERE biz_id='122';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType)(rate(dm_global_status_dps{$INSTANCE}[1m]))' WHERE biz_id='123';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType)(rate(dm_global_status_ups{$INSTANCE}[1m]))' WHERE biz_id='124';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType)(rate(dm_global_status_ddlps{$INSTANCE}[1m]))' WHERE biz_id='125';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType)(rate(dm_global_status_nioips{$INSTANCE}[1m]))' WHERE biz_id='126';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType)(rate(dm_global_status_nio_ops{$INSTANCE}[1m]))' WHERE biz_id='127';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType)(rate(dm_global_status_fio_ips{$INSTANCE}[1m]))' WHERE biz_id='128';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType)(rate(dm_global_status_fio_ops{$INSTANCE}[1m]))' WHERE biz_id='129';
UPDATE t_monitor_item SET metrics_linux='avg by(instance,instanceType)(dm_global_status_mem_used{$INSTANCE})' WHERE biz_id='130';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType)(dm_global_status_cpu_use_rate{$INSTANCE})' WHERE biz_id='131';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType)(dm_global_status_sessions{$INSTANCE})' WHERE biz_id='132';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType)(dm_global_status_active_sessions{$INSTANCE})' WHERE biz_id='133';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType)(dm_global_status_task_waiting{$INSTANCE})' WHERE biz_id='134';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType)(dm_global_status_task_ready{$INSTANCE})' WHERE biz_id='135';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType)(dm_global_status_task_total_wait_time{$INSTANCE})' WHERE biz_id='136';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType)(dm_global_status_avg_wait_time{$INSTANCE})' WHERE biz_id='137';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType)(dm_global_status_threads{$INSTANCE})' WHERE biz_id='138';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType)(pg_cpu_usage{$INSTANCE})' WHERE biz_id='139';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType)(pg_mem_usage{$INSTANCE})' WHERE biz_id='140';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType)(pg_disk_usage{$INSTANCE})' WHERE biz_id='141';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType)(pg_qps{$INSTANCE})' WHERE biz_id='142';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType)(pg_rqps{$INSTANCE})' WHERE biz_id='143';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType)(pg_wqps{$INSTANCE})' WHERE biz_id='144';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType)(pg_tps{$INSTANCE})' WHERE biz_id='145';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType)(pg_mean_exec_time{$INSTANCE})' WHERE biz_id='146';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType)(pg_open_ct_num{$INSTANCE})' WHERE biz_id='147';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType)(pg_active_ct_num{$INSTANCE})' WHERE biz_id='148';
UPDATE t_monitor_item SET metrics_linux='sum by (instance,instanceType)(dm_global_status_mem_use_rate{$INSTANCE})' WHERE biz_id='166';
UPDATE t_monitor_item SET metrics_linux='100 - (100 * (sum by(instance,instanceType) (irate(ecs_cpu_seconds_total{mode="idle",$INSTANCE}[3m])) / sum by(instance,instanceType) (irate(ecs_cpu_seconds_total{$INSTANCE}[3m]))))' WHERE biz_id='167';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType)(ecs_load1{$INSTANCE})' WHERE biz_id='168';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType)(ecs_load5{$INSTANCE})' WHERE biz_id='169';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType)(ecs_load15{$INSTANCE})' WHERE biz_id='170';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType)(ecs_memory_MemTotal_bytes{$INSTANCE}) - sum by(instance,instanceType)(ecs_memory_MemFree_bytes{$INSTANCE})' WHERE biz_id='171';
UPDATE t_monitor_item SET metrics_linux='100 * ((sum by(instance,instanceType)(ecs_memory_MemTotal_bytes{$INSTANCE}) - sum by(instance,instanceType)(ecs_memory_MemFree_bytes{$INSTANCE})) / sum by(instance,instanceType)(ecs_memory_MemTotal_bytes{$INSTANCE}))' WHERE biz_id='172';
UPDATE t_monitor_item SET metrics_linux='100 * ((sum by(instance,instanceType,device)(ecs_filesystem_size_bytes{$INSTANCE}) - sum by(instance,instanceType,device)(ecs_filesystem_free_bytes{$INSTANCE})) / sum by(instance,instanceType,device)(ecs_filesystem_size_bytes{$INSTANCE}))' WHERE biz_id='173';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,device)(irate(ecs_disk_read_bytes_total{$INSTANCE}[3m]))' WHERE biz_id='174';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,device)(irate(ecs_disk_written_bytes_total{$INSTANCE}[3m]))' WHERE biz_id='175';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,device)(irate(ecs_disk_reads_completed_total{$INSTANCE}[3m]))' WHERE biz_id='176';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,device)(irate(ecs_disk_writes_completed_total{$INSTANCE}[3m]))' WHERE biz_id='177';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,device)(irate(ecs_network_receive_bytes_total{$INSTANCE}[3m])) / 1024 / 1024 * 8' WHERE biz_id='178';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,device)(irate(ecs_network_transmit_bytes_total{$INSTANCE}[3m])) / 1024 / 1024 * 8' WHERE biz_id='179';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,device)(irate(ecs_network_receive_packets_total{$INSTANCE}[3m]))' WHERE biz_id='180';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,device)(irate(ecs_network_transmit_packets_total{$INSTANCE}[3m]))' WHERE biz_id='181';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,device)(ecs_filesystem_free_bytes{$INSTANCE})' WHERE biz_id='182';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,device)(ecs_filesystem_size_bytes{$INSTANCE}) - sum by(instance,instanceType,device)(ecs_filesystem_free_bytes{$INSTANCE})' WHERE biz_id='183';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,device)(ecs_filesystem_size_bytes{$INSTANCE})' WHERE biz_id='184';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType)(redis_cpu_usage{$INSTANCE})' WHERE biz_id='185';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType)(redis_mem_usage{$INSTANCE})' WHERE biz_id='186';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType)(redis_connected_clients{$INSTANCE})' WHERE biz_id = '187';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType)(redis_tps{$INSTANCE})' WHERE biz_id = '188';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,pod)(mongo_mongos_current_connections{$INSTANCE})' WHERE biz_id='195';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,pod)(mongo_shard_current_connections{$INSTANCE})' WHERE biz_id='196';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,pod)(mongo_config_current_connections{$INSTANCE})' WHERE biz_id='197';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType)(mongo_total_current_connections{$INSTANCE})' WHERE biz_id='198';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,pod)(mongo_mongos_memory_ratio{$INSTANCE})' WHERE biz_id='199';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType)(mongo_config_memory_ratio{$INSTANCE})' WHERE biz_id='200';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,pod)(mongo_shard_memory_ratio{$INSTANCE})' WHERE biz_id='201';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,pod)(mongo_mongos_cpu_ratio{$INSTANCE})' WHERE biz_id='202';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,pod)(mongo_shard_cpu_ratio{$INSTANCE})' WHERE biz_id='203';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType)(mongo_config_cpu_ratio{$INSTANCE})' WHERE biz_id='204';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,service,route)(rate(guard_nginx_http_current_reqs{$INSTANCE}[3m]))' WHERE biz_id='205';
UPDATE t_monitor_item SET metrics_linux='histogram_quantile(0.90, sum by(instance,instanceType,service,route,le)(rate(guard_http_latency_bucket{type="request",$INSTANCE}[3m])))' WHERE biz_id='206';
UPDATE t_monitor_item SET metrics_linux='histogram_quantile(0.95, sum by(instance,instanceType,service,route,le)(rate(guard_http_latency_bucket{type="request",$INSTANCE}[3m])))' WHERE biz_id='207';
UPDATE t_monitor_item SET metrics_linux='histogram_quantile(0.99, sum by(instance,instanceType,service,route,le)(rate(guard_http_latency_bucket{type="request",$INSTANCE}[3m])))' WHERE biz_id='208';
UPDATE t_monitor_item SET metrics_linux='histogram_quantile(0.90, sum by(instance,instanceType,service,le)(rate(guard_http_latency_bucket{type="request",$INSTANCE}[3m])))' WHERE biz_id='209';
UPDATE t_monitor_item SET metrics_linux='histogram_quantile(0.95, sum by(instance,instanceType,service,le)(rate(guard_http_latency_bucket{type="request",$INSTANCE}[3m])))' WHERE biz_id='210';
UPDATE t_monitor_item SET metrics_linux='histogram_quantile(0.99, sum by(instance,instanceType,service,le)(rate(guard_http_latency_bucket{type="request",$INSTANCE}[3m])))' WHERE biz_id='211';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,service,route)(rate(guard_nginx_url_request_succ{code="200",$INSTANCE}[3m]))/sum by(instance,instanceType,service,route)(rate(guard_nginx_url_request_succ{code="total",$INSTANCE}[3m]))*100' WHERE biz_id='212';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,service)(rate(guard_nginx_url_request_succ{code="200",$INSTANCE}[3m]))/sum by(instance,instanceType,service)(rate(guard_nginx_url_request_succ{code="total",$INSTANCE}[3m]))*100' WHERE biz_id='213';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,service,route)(rate(guard_bandwidth{type="ingress",$INSTANCE}[3m]))' WHERE biz_id='214';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,service,route)(rate(guard_bandwidth{type="egress",$INSTANCE}[3m]))' WHERE biz_id='215';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,service)(irate(guard_bandwidth{type="ingress",$INSTANCE}[3m])) ' WHERE biz_id='216';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,service)(rate(guard_bandwidth{type="egress",$INSTANCE}[3m]))' WHERE biz_id='217';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,service,route)(guard_http_apirequests{$INSTANCE})' WHERE biz_id='218';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,service)(guard_http_apirequests{$INSTANCE})' WHERE biz_id='219';
UPDATE t_monitor_item SET metrics_linux='sum by (instance,instanceType) (kafka_brokers{$INSTANCE})' WHERE biz_id='149';
UPDATE t_monitor_item SET metrics_linux='sum by (instance,instanceType) (kafka_server_replicamanager_partitioncount{$INSTANCE})' WHERE biz_id='150';
UPDATE t_monitor_item SET metrics_linux='sum by (instance,instanceType) (kafka_server_brokertopicmetrics_bytesinpersec{$INSTANCE})' WHERE biz_id='151';
UPDATE t_monitor_item SET metrics_linux='sum by (instance,instanceType) (kafka_server_brokertopicmetrics_bytesoutpersec{$INSTANCE})' WHERE biz_id='152';
UPDATE t_monitor_item SET metrics_linux='sum by (instance,instanceType) (kafka_server_brokertopicmetrics_messagesinpersec{$INSTANCE})' WHERE biz_id='153';
UPDATE t_monitor_item SET metrics_linux='sum by (instance,instanceType) (kafka_consumergroup_lag{$INSTANCE})' WHERE biz_id='154';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,cmd_line)(ecs_processes_top5Cpus{cmd_line!="",$INSTANCE})' WHERE biz_id='220';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,cmd_line)(ecs_processes_top5Mems{cmd_line!="",$INSTANCE})' WHERE biz_id='221';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType)(ecs_procs_running{$INSTANCE})' WHERE biz_id='222';
UPDATE t_monitor_item SET metrics_linux='sum by(instance,instanceType,cmd_line)(ecs_processes_top5Fds{cmd_line!="",$INSTANCE})' WHERE biz_id='223';
UPDATE t_monitor_item SET metrics_linux='sum(eip_upstream_bits_rate{$INSTANCE})/8/1024' WHERE biz_id='224';
UPDATE t_monitor_item SET metrics_linux='sum(eip_downstream_bits_rate{$INSTANCE})/8/1024' WHERE biz_id='225';



UPDATE t_monitor_item SET metrics_linux = 'ecs_processes_top5Cpus{$INSTANCE}', labels = 'instance,pid' WHERE biz_id = '220';
UPDATE t_monitor_item SET metrics_linux = 'ecs_processes_top5Mems{$INSTANCE}', labels = 'instance,pid' WHERE biz_id = '221';
UPDATE t_monitor_item SET metrics_linux = 'ecs_procs_running{$INSTANCE}' WHERE biz_id = '222';
UPDATE t_monitor_item SET metrics_linux = 'ecs_processes_top5Fds{$INSTANCE}', labels = 'instance,pid' WHERE biz_id = '223';

UPDATE t_monitor_product SET page_url = '/inner/cmq/v1/kafka/monitor/listAllCluster', status = '1' WHERE abbreviation IN ('kafka');

UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(redis_connected_clients{$INSTANCE})' WHERE biz_id = '187';
UPDATE t_monitor_item SET metrics_linux = 'sum by(instance)(redis_tps{$INSTANCE})' WHERE biz_id = '188';

UPDATE t_monitor_item SET metric_name = 'nat_snat_connection' WHERE biz_id = '55';
UPDATE t_monitor_item SET metric_name = 'nat_inbound_bandwidth' WHERE biz_id = '56';
UPDATE t_monitor_item SET metric_name = 'nat_outbound_bandwidth' WHERE biz_id = '57';
UPDATE t_monitor_item SET metric_name = 'nat_inbound_traffic' WHERE biz_id = '58';
UPDATE t_monitor_item SET metric_name = 'nat_outbound_traffic' WHERE biz_id = '59';
UPDATE t_monitor_item SET metric_name = 'nat_inbound_pps' WHERE biz_id = '60';
UPDATE t_monitor_item SET metric_name = 'nat_outbound_pps' WHERE biz_id = '61';
UPDATE t_monitor_item SET metric_name = 'nat_snat_connection_ratio' WHERE biz_id = '62';

INSERT INTO `t_monitor_item` (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('226', '15', 'NAT连接数', 'nat_e_total_connection', 'instance', 'sum by (instance)(Nat_snat_total_connection_count{$INSTANCE})', NULL, NULL, '个', NULL, NULL, 1, 1, NULL, NULL, NULL, NULL, 'chart,rule');
INSERT INTO `t_monitor_item` (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('227', '15', '入方向带宽', 'nat_e_inbound_bandwidth', 'instance', 'sum by (instance)(rate(Nat_recv_bytes_total_count{$INSTANCE}[1m])*8)', NULL, NULL, 'bps', NULL, NULL, 1, 1, NULL, NULL, NULL, NULL, 'chart,rule');
INSERT INTO `t_monitor_item` (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('228', '15', '出方向带宽', 'nat_e_outbound_bandwidth', 'instance', 'sum by (instance)(rate(Nat_send_bytes_total_count{$INSTANCE}[1m])*8)', NULL, NULL, 'bps', NULL, NULL, 1, 1, NULL, NULL, NULL, NULL, 'chart,rule');
INSERT INTO `t_monitor_item` (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('229', '15', '入方向流量', 'nat_e_inbound_traffic', 'instance', 'sum by (instance)(Nat_recv_bytes_total_count{$INSTANCE})', NULL, NULL, 'Byte', NULL, NULL, 1, 1, NULL, NULL, NULL, NULL, 'chart,rule');
INSERT INTO `t_monitor_item` (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('230', '15', '出方向流量', 'nat_e_outbound_traffic	', 'instance', 'sum by (instance)(Nat_send_bytes_total_count{$INSTANCE})', NULL, NULL, 'Byte', NULL, NULL, 1, 1, NULL, NULL, NULL, NULL, 'chart,rule');
INSERT INTO `t_monitor_item` (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('231', '15', '入方向PPS', 'nat_e_inbound_pps', 'instance', 'sum by (instance)(rate(Nat_recv_packets_total_count{$INSTANCE}[1m]))', NULL, NULL, 'pps', NULL, NULL, 1, 1, NULL, NULL, NULL, NULL, 'chart,rule');
INSERT INTO `t_monitor_item` (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('232', '15', '出方向PPS', 'nat_e_outbound_pps', 'instance', 'sum by (instance)(rate(Nat_send_packets_total_count{$INSTANCE}[1m]))', NULL, NULL, 'pps', NULL, NULL, 1, 1, NULL, NULL, NULL, NULL, 'chart,rule');
INSERT INTO `t_monitor_item` (biz_id, product_biz_id, name, metric_name, labels, metrics_linux, metrics_windows, statistics, unit, frequency, type, is_display, status, description, create_user, create_time, show_expression, display) VALUES ('233', '15', 'NAT连接数使用率', 'nat_e_total_connection_ratio', 'instance', 'sum by (instance)(Nat_snat_total_connection_count{$INSTANCE}) / avg by (instance)(Nat_nat_max_connection_count{$INSTANCE}) *100', NULL, NULL, '%', NULL, NULL, 1, 1, NULL, NULL, NULL, NULL, 'chart,rule');

ALTER TABLE t_monitor_product ADD COLUMN iam_page_url varchar(256) COMMENT 'iam请求路径';
UPDATE t_monitor_product SET iam_page_url = '/compute/ecs/instance/cbc/pageList' WHERE abbreviation IN ('ecs');
UPDATE t_monitor_product SET iam_page_url = '/eip/inner/eipInfoList' WHERE abbreviation IN ('eip');
UPDATE t_monitor_product SET iam_page_url = '/slb/inner/list' WHERE abbreviation IN ('slb');
UPDATE t_monitor_product SET iam_page_url = '/compute/monitor/vault/pageList' WHERE abbreviation IN ('cbr');
UPDATE t_monitor_product SET iam_page_url = '/nat-gw/inner/nat/page' WHERE abbreviation IN ('nat');
UPDATE t_monitor_product SET iam_page_url = '/nat-gw/nat/page' WHERE abbreviation IN ('nat-e');
UPDATE t_monitor_product SET iam_page_url = '/v1/mysql/instance' WHERE abbreviation IN ('mysql');
UPDATE t_monitor_product SET iam_page_url = '/v1/dm/instance' WHERE abbreviation IN ('dm');
UPDATE t_monitor_product SET iam_page_url = '/v1/pg/instance/' WHERE abbreviation IN ('postgresql');
UPDATE t_monitor_product SET iam_page_url = '/inner/cmq/v1/kafka/monitor/listAllCluster' WHERE abbreviation IN ('kafka');
UPDATE t_monitor_product SET iam_page_url = '/compute/bms/ops/v1/tenants/{tenantId}/servers' WHERE abbreviation IN ('bms');
UPDATE t_monitor_product SET iam_page_url = '/compute/ebms/ops/v1/tenants/{tenantId}/servers' WHERE abbreviation IN ('ebms');
UPDATE t_monitor_product SET iam_page_url = '/v1/redis/instance' WHERE abbreviation IN ('redis');
UPDATE t_monitor_product SET iam_page_url = '/v2/mongo/instance' WHERE abbreviation IN ('mongo');
UPDATE t_monitor_product SET iam_page_url = '/gateway/instance/page' WHERE abbreviation IN ('cgw');

UPDATE t_monitor_item SET metrics_linux = 'clamp_max((sum(eip_upstream_bits_rate{$INSTANCE}) by (instance,eip) / avg(eip_config_upstream_bandwidth{$INSTANCE}) by (instance,eip)) * 100, 100)' WHERE metric_name = 'eip_upstream_bandwidth_usage';
UPDATE t_monitor_item SET metrics_linux = 'clamp_max(sum by(instance,instanceType)(dm_global_status_cpu_use_rate{$INSTANCE}), 100)' WHERE metric_name = 'dm_global_status_cpu_use_rate';

INSERT INTO t_monitor_product (biz_id, name, status, description, create_user, create_time, route, cron, host, page_url, abbreviation, sort, monitor_type, iam_page_url) VALUES ('15', '增强型NAT网关', '1', 'nat-e', null, null, '/productmonitoring/nat-e', '0 0 0/1 * * ?', 'http://product-nat-controller-nat-manage.product-nat-gw', '/nat-gw/inner/nat/page', 'nat-e', '6', '云产品监控', '/nat-gw/inner/nat/page');

ALTER TABLE t_alarm_rule_resource_rel ADD INDEX alarm_rule_id ( alarm_rule_id );
