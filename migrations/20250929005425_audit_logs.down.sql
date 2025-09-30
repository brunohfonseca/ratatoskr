-- Drop índices
DROP INDEX IF EXISTS idx_audit_log_changed_at;
DROP INDEX IF EXISTS idx_audit_log_user;
DROP INDEX IF EXISTS idx_audit_log_table_record;

-- Drop audit table
DROP TABLE IF EXISTS audit_log;