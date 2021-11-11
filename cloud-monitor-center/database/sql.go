package database

var SelectAlterContact = "SELECT " +
	"ac.id AS contact_id, " +
	"ac.name AS contact_name, " +
	"acg.group_id AS group_id, " +
	"acg.group_name AS group_name, " +
	"GROUP_CONCAT( CASE aci.type WHEN 1 THEN aci.NO END ) AS phone, " +
	"GROUP_CONCAT( CASE aci.type WHEN 1 THEN aci.is_certify END ) AS phone_certify, " +
	"GROUP_CONCAT( CASE aci.type WHEN 2 THEN aci.NO END ) AS email, " +
	"GROUP_CONCAT( CASE aci.type WHEN 2 THEN aci.is_certify END ) AS email_certify, " +
	"GROUP_CONCAT( CASE aci.type WHEN 3 THEN aci.NO END ) AS lanxin, " +
	"GROUP_CONCAT( CASE aci.type WHEN 3 THEN aci.is_certify END ) AS lanxin_certify, " +
	"ac.description AS description " +
	"FROM " +
	"alert_contact AS ac " +
	"LEFT JOIN alert_contact_information AS aci ON ac.id = aci.contact_id " +
	"LEFT JOIN ( " +
	"SELECT " +
	"acgr.contact_id AS contact_id, " +
	"GROUP_CONCAT( acg.id ) AS group_id, " +
	"GROUP_CONCAT( acg.name ) AS group_name " +
	"FROM " +
	"alert_contact_group AS acg " +
	"LEFT JOIN alert_contact_group_rel AS acgr ON acg.id = acgr.group_id  " +
	"GROUP BY " +
	"acgr.contact_id ) " +
	"AS acg ON acg.contact_id = ac.id " +
	"WHERE " +
	"ac.status = 1 " +
	"AND ac.tenant_id = %s " +
	"%s" +
	"GROUP BY " +
	"ac.id " +
	"ORDER BY " +
	"ac.create_time DESC  "

var SelectAlterContactGroup = "SELECT " +
	"acg.id AS group_id, " +
	"acg.name AS group_name, " +
	"acg.description AS description, " +
	"acg.create_time AS create_time, " +
	"acg.update_time AS update_time, " +
	"COUNT( acgr.group_id ) AS contact_count " +
	"FROM " +
	"alert_contact_group AS acg " +
	"LEFT JOIN alert_contact_group_rel AS acgr ON acg.id = acgr.group_id " +
	"WHERE " +
	"acg.tenant_id = ? " +
	"AND acg.name LIKE CONCAT('%',?,'%') " +
	"GROUP BY " +
	"acg.id " +
	"ORDER BY " +
	"acg.create_time DESC "

var SelectAlterGroupContact = "SELECT " +
	"ac.id AS contact_id, " +
	"ac.name AS contact_name, " +
	"acg.group_id AS group_id, " +
	"acg.group_name AS group_name, " +
	"GROUP_CONCAT( CASE aci.type WHEN 1 THEN aci.NO END ) AS phone, " +
	"GROUP_CONCAT( CASE aci.type WHEN 1 THEN aci.is_certify END ) AS phone_certify, " +
	"GROUP_CONCAT( CASE aci.type WHEN 2 THEN aci.NO END ) AS email, " +
	"GROUP_CONCAT( CASE aci.type WHEN 2 THEN aci.is_certify END ) AS email_certify, " +
	"GROUP_CONCAT( CASE aci.type WHEN 3 THEN aci.NO END ) AS lanxin, " +
	"GROUP_CONCAT( CASE aci.type WHEN 3 THEN aci.is_certify END ) AS lanxin_certify, " +
	"ac.description AS description " +
	"FROM " +
	"alert_contact AS ac " +
	"LEFT JOIN alert_contact_information AS aci ON ac.id = aci.contact_id " +
	"LEFT JOIN ( " +
	"SELECT " +
	"acgr.contact_id AS contact_id, " +
	"GROUP_CONCAT( acg.id ) AS group_id, " +
	"GROUP_CONCAT( acg.name ) AS group_name " +
	"FROM " +
	"alert_contact_group AS acg " +
	"LEFT JOIN alert_contact_group_rel AS acgr ON acg.id = acgr.group_id  " +
	"GROUP BY " +
	"acgr.contact_id ) " +
	"AS acg ON acg.contact_id = ac.id " +
	"WHERE " +
	"ac.status = 1 " +
	"AND ac.tenant_id = ? " +
	"AND acg.group_id = ? " +
	"GROUP BY " +
	"ac.id " +
	"ORDER BY " +
	"ac.create_time DESC "
