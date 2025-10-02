-- Enum para status das checagens
CREATE TYPE check_status AS ENUM ('online', 'offline', 'timeout', 'error', 'unknown');
CREATE TYPE alert_channel_type AS ENUM ('slack', 'telegram', 'email');
CREATE TYPE endpoint_ssl_status AS ENUM ('valid', 'warning', 'expired', 'error', 'unknown');
CREATE TYPE auth_provider AS ENUM ('local', 'keycloak');

-- Usuários
CREATE TABLE users (
    uuid UUID DEFAULT uuidv7() PRIMARY KEY,
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
    uuid UUID DEFAULT uuidv7() PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    description VARCHAR(300),
    last_modified_by UUID REFERENCES users(uuid) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- Canais (slack, telegram, etc.)
CREATE TABLE alert_channels (
    uuid UUID DEFAULT uuidv7() PRIMARY KEY,
    type alert_channel_type NOT NULL,
    name VARCHAR(50) NOT NULL,
    config JSONB NOT NULL,
    last_modified_by UUID REFERENCES users(uuid) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- Relação N:N entre grupos e canais
CREATE TABLE alert_group_channels (
    uuid UUID DEFAULT uuidv7() PRIMARY KEY,
    group_id UUID NOT NULL REFERENCES alert_groups(uuid) ON DELETE CASCADE,
    channel_id UUID NOT NULL REFERENCES alert_channels(uuid) ON DELETE CASCADE
);

-- Histórico de alertas enviados
CREATE TABLE endpoints (
    uuid UUID DEFAULT uuidv7() PRIMARY KEY,
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
    ssl_expiration_date TIMESTAMPTZ,
    ssl_issuer VARCHAR(80) NOT NULL,
    ssl_status endpoint_ssl_status NOT NULL,
    ssl_last_check TIMESTAMPTZ DEFAULT now(),
    ssl_error VARCHAR(300),
    alert_group_id UUID REFERENCES alert_groups(uuid) ON DELETE SET NULL,
    last_modified_by UUID REFERENCES users(uuid) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- Resultado da checagem SSL
CREATE TABLE endpoint_checks (
    uuid UUID DEFAULT uuidv7() PRIMARY KEY,
    endpoint_id UUID NOT NULL REFERENCES endpoints(uuid) ON DELETE CASCADE,
    status check_status NOT NULL,
    response_time_ms INT,
    response_message VARCHAR(300),
    created_at TIMESTAMPTZ DEFAULT now()
);

-- Histórico das checagens (time-series)
CREATE TABLE sent_alerts (
    uuid UUID DEFAULT uuidv7() PRIMARY KEY,
    endpoint_id UUID NOT NULL REFERENCES endpoints(uuid) ON DELETE CASCADE,
    channel_id UUID NOT NULL REFERENCES alert_channels(uuid) ON DELETE CASCADE,
    message VARCHAR(500) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);
