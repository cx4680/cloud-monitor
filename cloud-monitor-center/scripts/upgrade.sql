-- 1 notify channel  拆解为t_alarm_handler
INSERT INTO t_alarm_handler (
	`alarm_rule_id`,
	`handle_type`,
	`handle_params`,
	`tenant_id`,
	`create_time`,
	update_time
) SELECT
	id,
	CASE notify_channel
WHEN 2 THEN
	1
ELSE
	2
END AS handle_type,
 '',
 tenant_id,
 create_time,
 create_time
FROM
	t_alarm_rule
WHERE
	notify_channel != 3 ;

-- 2规则与实例的关系拆解
	INSERT INTO t_alarm_rule_resource_rel (
		`alarm_rule_id`,
		`resource_id`,
		`tenant_id`,
		`create_time`,
		`update_time`
	) SELECT
		t1.id,
		t2.instance_id,
		t1.tenant_id,
		t2.create_time,
		t2.create_time
	FROM
		t_alarm_rule t1,
		t_alarm_instance t2
	WHERE
		t1.id = t2.alarm_rule_id ;

-- 3去重 删除多余实例
		DELETE
	FROM
		t_alarm_instance
	WHERE
		instance_id IN (
			SELECT
				*
			FROM
				(
					SELECT
						instance_id
					FROM
						t_alarm_instance
					GROUP BY
						instance_id
					HAVING
						count(instance_id) > 1
				) t2
		)
	AND create_time NOT IN (
		SELECT
			*
		FROM
			(
				SELECT
					min(create_time)
				FROM
					t_alarm_instance
				GROUP BY
					instance_id
				HAVING
					count(instance_id) > 1
			) t
	);
-- 4实例增加produt_type
update t_alarm_instance  t1
left JOIN t_alarm_rule t2  on t1.alarm_rule_id=t2.id
set t1.product_type = t2.product_type