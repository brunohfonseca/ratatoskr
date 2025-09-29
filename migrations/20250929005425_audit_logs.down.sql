-- Drop trigger first (depends on endpoints table)
DROP TRIGGER IF EXISTS endpoints_audit ON endpoints;
-- Drop function used by the trigger
DROP FUNCTION IF EXISTS log_endpoint_changes();
-- Drop audit table
DROP TABLE IF EXISTS audit_log;