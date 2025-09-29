-- Log de auditoria
CREATE TABLE audit_log (
    id BIGSERIAL PRIMARY KEY,
    table_name TEXT NOT NULL,
    record_id TEXT NOT NULL,
    operation TEXT NOT NULL,              -- INSERT, UPDATE, DELETE
    changed_data JSONB NOT NULL,
    changed_at TIMESTAMPTZ DEFAULT now(),
    user_id INT REFERENCES users(id)
);

-- Função para auditoria de endpoints (pode ser adaptada para outras tabelas)
CREATE OR REPLACE FUNCTION log_endpoint_changes()
RETURNS TRIGGER AS $$
DECLARE
acting_user INT;
BEGIN
    -- pega usuário atual da sessão (setado pela app após login)
    acting_user := current_setting('app.current_user', true)::int;

    IF TG_OP = 'UPDATE' THEN
        INSERT INTO audit_log (table_name, record_id, operation, changed_data, user_id)
        VALUES ('endpoints', NEW.id::text, 'UPDATE', row_to_json(NEW)::jsonb, acting_user);
RETURN NEW;
ELSIF TG_OP = 'INSERT' THEN
        INSERT INTO audit_log (table_name, record_id, operation, changed_data, user_id)
        VALUES ('endpoints', NEW.id::text, 'INSERT', row_to_json(NEW)::jsonb, acting_user);
RETURN NEW;
ELSIF TG_OP = 'DELETE' THEN
        INSERT INTO audit_log (table_name, record_id, operation, changed_data, user_id)
        VALUES ('endpoints', OLD.id::text, 'DELETE', row_to_json(OLD)::jsonb, acting_user);
RETURN OLD;
END IF;
RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER endpoints_audit
    AFTER INSERT OR UPDATE OR DELETE ON endpoints
    FOR EACH ROW EXECUTE FUNCTION log_endpoint_changes();
