-- Log de auditoria
CREATE TABLE audit_log (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID DEFAULT uuidv7() NOT NULL,
    table_name TEXT NOT NULL,
    record_id TEXT NOT NULL,
    operation TEXT NOT NULL,              -- INSERT, UPDATE, DELETE
    changed_data JSONB NOT NULL,
    changed_at TIMESTAMPTZ DEFAULT now(),
    user_id INT REFERENCES users(id)
);

-- Índice para melhorar performance de consultas
CREATE INDEX idx_audit_log_table_record ON audit_log(table_name, record_id);
CREATE INDEX idx_audit_log_user ON audit_log(user_id);
CREATE INDEX idx_audit_log_changed_at ON audit_log(changed_at DESC);