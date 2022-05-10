RENAME TABLE  `t_config_item` TO `config_item`;
RENAME TABLE  `t_monitor_product` TO `monitor_product` ;
RENAME TABLE`t_monitor_item`  TO `monitor_item` ;
RENAME TABLE  `t_notification_record` TO `notification_record` ;
RENAME TABLE  `t_alarm_record` TO `t_alert_record` ;

ALTER TABLE `config_item` MODIFY id BIGINT UNSIGNED;
alter table `config_item` drop column `id`;
ALTER TABLE `config_item` CHANGE `biz_id` `id` VARCHAR (50) primary key  ;
ALTER TABLE `config_item` CHANGE  `p_biz_id` `pid` VARCHAR (50);

ALTER TABLE `monitor_product` MODIFY id BIGINT UNSIGNED;
alter table `monitor_product` drop column `id`;
ALTER TABLE `monitor_product` CHANGE `biz_id` `id` VARCHAR (50)  primary key ;



ALTER TABLE `monitor_item` MODIFY id BIGINT UNSIGNED;
alter table `monitor_item` drop column `id`;
ALTER TABLE `monitor_item` CHANGE `biz_id` `id` VARCHAR (50)  primary key  ;
ALTER TABLE `monitor_item` CHANGE  `product_biz_id` `product_id` VARCHAR (50);


ALTER TABLE `notification_record` MODIFY id BIGINT UNSIGNED;
alter table `notification_record` drop column `id`;
ALTER TABLE `notification_record` CHANGE `biz_id` `id` VARCHAR (50)  primary key  ;

drop index tenant_time_idx on notification_record;

drop table t_alarm_info;

ALTER TABLE `t_alert_record` CHANGE  `biz_id` `id` VARCHAR (50) primary key comment '业务Id';


