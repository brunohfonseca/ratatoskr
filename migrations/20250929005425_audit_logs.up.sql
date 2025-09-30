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

-- Função GENÉRICA para auditoria (funciona em qualquer tabela)
CREATE OR REPLACE FUNCTION log_table_changes()
RETURNS TRIGGER AS $$
DECLARE
    record_id_value TEXT;
    user_id_value INT;
BEGIN
    -- Extrai o ID do registro (prioriza 'id', mas funciona com outros campos também)
    IF TG_OP = 'DELETE' THEN
        record_id_value := COALESCE(OLD.id::text, OLD.uuid::text, 'unknown');
        user_id_value := OLD.last_modified_by;
    ELSE
        record_id_value := COALESCE(NEW.id::text, NEW.uuid::text, 'unknown');
        user_id_value := NEW.last_modified_by;
    END IF;

    -- Log da operação
    IF TG_OP = 'UPDATE' THEN
        INSERT INTO audit_log (table_name, record_id, operation, changed_data, user_id)
        VALUES (TG_TABLE_NAME, record_id_value, 'UPDATE', row_to_json(NEW)::jsonb, user_id_value);
        RETURN NEW;
    ELSIF TG_OP = 'INSERT' THEN
        INSERT INTO audit_log (table_name, record_id, operation, changed_data, user_id)
        VALUES (TG_TABLE_NAME, record_id_value, 'INSERT', row_to_json(NEW)::jsonb, user_id_value);
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        INSERT INTO audit_log (table_name, record_id, operation, changed_data, user_id)
        VALUES (TG_TABLE_NAME, record_id_value, 'DELETE', row_to_json(OLD)::jsonb, user_id_value);
        RETURN OLD;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Triggers para todas as tabelas que precisam de auditoria
CREATE TRIGGER endpoints_audit
    AFTER INSERT OR UPDATE OR DELETE ON endpoints
    FOR EACH ROW EXECUTE FUNCTION log_table_changes();

CREATE TRIGGER users_audit
    AFTER INSERT OR UPDATE OR DELETE ON users
    FOR EACH ROW EXECUTE FUNCTION log_table_changes();

CREATE TRIGGER alert_groups_audit
    AFTER INSERT OR UPDATE OR DELETE ON alert_groups
    FOR EACH ROW EXECUTE FUNCTION log_table_changes();

CREATE TRIGGER alert_channels_audit
    AFTER INSERT OR UPDATE OR DELETE ON alert_channels
    FOR EACH ROW EXECUTE FUNCTION log_table_changes();

CREATE TRIGGER endpoint_ssl_audit
    AFTER INSERT OR UPDATE OR DELETE ON endpoint_ssl
    FOR EACH ROW EXECUTE FUNCTION log_table_changes();

CREATE TRIGGER endpoint_auth_audit
    AFTER INSERT OR UPDATE OR DELETE ON endpoint_auth
    FOR EACH ROW EXECUTE FUNCTION log_table_changes();
