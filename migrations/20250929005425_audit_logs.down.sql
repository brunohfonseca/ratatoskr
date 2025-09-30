-- Drop trigger first (depends on endpoints table)
DROP TRIGGER IF EXISTS endpoints_audit ON endpoints;
DROP TRIGGER IF EXISTS alert_groups_audit ON alert_groups;
DROP TRIGGER IF EXISTS alert_group_channels_audit ON alert_group_channels;
DROP TRIGGER IF EXISTS alert_channels_audit ON alert_channels;
DROP TRIGGER IF EXISTS endpoint_ssl_audit ON endpoint_ssl;
DROP TRIGGER IF EXISTS endpoint_auth_audit ON endpoint_auth;
DROP TRIGGER IF EXISTS endpoints_audit ON endpoints;
DROP TRIGGER IF EXISTS endpoint_checks_audit ON endpoint_checks;
-- Drop function used by the trigger
DROP FUNCTION IF EXISTS log_table_changes();
-- Drop audit table
DROP TABLE IF EXISTS audit_log;