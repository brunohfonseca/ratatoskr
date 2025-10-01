-- Enum para status das checagens
CREATE TYPE check_status AS ENUM ('ok', 'down', 'timeout', 'error', 'unknown');
CREATE TYPE alert_channel_type AS ENUM ('slack', 'telegram', 'email');
CREATE TYPE endpoint_ssl_status AS ENUM ('ok', 'warning', 'expired');
CREATE TYPE auth_provider AS ENUM ('local', 'keycloak');

-- Usuários
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    uuid UUID DEFAULT uuidv7(),
    email VARCHAR(70) NOT NULL UNIQUE,
    full_name VARCHAR(80),
    password_hash VARCHAR(500),
    auth_provider auth_provider NOT NULL DEFAULT 'local',
    enabled BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- Grupos de alerta
CREATE TABLE alert_groups (
    id SERIAL PRIMARY KEY,
    uuid UUID DEFAULT uuidv7(),
    name VARCHAR(50) NOT NULL,
    description VARCHAR(300),
    last_modified_by INT REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- Canais (slack, telegram, etc.)
CREATE TABLE alert_channels (
    id SERIAL PRIMARY KEY,
    uuid UUID DEFAULT uuidv7(),
    type alert_channel_type NOT NULL,
    name VARCHAR(50) NOT NULL,
    config JSONB NOT NULL,
    last_modified_by INT REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- Relação N:N entre grupos e canais
CREATE TABLE alert_group_channels (
    id SERIAL PRIMARY KEY,
    uuid UUID DEFAULT uuidv7(),
    group_id INT NOT NULL REFERENCES alert_groups(id) ON DELETE CASCADE,
    channel_id INT NOT NULL REFERENCES alert_channels(id) ON DELETE CASCADE
);

-- Histórico de alertas enviados
CREATE TABLE endpoints (
    id SERIAL PRIMARY KEY,
    uuid UUID DEFAULT uuidv7(),
    name VARCHAR(25) NOT NULL,
    domain VARCHAR(75) NOT NULL,
    path VARCHAR(30) DEFAULT '/',
    check_ssl BOOLEAN DEFAULT FALSE,
    enabled BOOLEAN DEFAULT TRUE,
    status check_status NOT NULL DEFAULT 'unknown',
    expected_response_code INT,
    response_code INT,
    response_message VARCHAR(300),
    response_time_ms INT,
    timeout_seconds INT DEFAULT 30,
    interval_seconds INT DEFAULT 300,
    alert_group_id INT REFERENCES alert_groups(id) ON DELETE SET NULL,
    last_modified_by INT REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- Endpoints monitorados
CREATE TABLE endpoint_ssl (
    id SERIAL PRIMARY KEY,
    uuid UUID DEFAULT uuidv7(),
    endpoint_id INT NOT NULL REFERENCES endpoints(id) ON DELETE CASCADE,
    expiration_date TIMESTAMPTZ NOT NULL,
    issuer VARCHAR(40) NOT NULL,
    status endpoint_ssl_status NOT NULL,
    last_check TIMESTAMPTZ DEFAULT now(),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- Resultado da checagem SSL
CREATE TABLE endpoint_checks (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID DEFAULT uuidv7(),
    endpoint_id INT NOT NULL REFERENCES endpoints(id) ON DELETE CASCADE,
    status check_status NOT NULL,
    response_time_ms INT,
    response_message VARCHAR(300),
    created_at TIMESTAMPTZ DEFAULT now()
);

-- Autenticação de endpoints (JSONB flexível)
CREATE TABLE endpoint_auth (
    uuid UUID DEFAULT uuidv7(),
    endpoint_id INT PRIMARY KEY REFERENCES endpoints(id) ON DELETE CASCADE,
    config JSONB NOT NULL,
    last_modified_by INT REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- Histórico das checagens (time-series)
CREATE TABLE sent_alerts (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID DEFAULT uuidv7(),
    endpoint_id INT NOT NULL REFERENCES endpoints(id) ON DELETE CASCADE,
    channel_id INT NOT NULL REFERENCES alert_channels(id) ON DELETE CASCADE,
    message VARCHAR(500) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);
