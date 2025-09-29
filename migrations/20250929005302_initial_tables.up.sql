-- Enum para status das checagens
CREATE TYPE check_status AS ENUM ('ok', 'down', 'ssl_expired', 'timeout', 'error');

-- Grupos de alerta
CREATE TABLE alert_groups (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT now()
);

-- Canais (slack, telegram, etc.)
CREATE TABLE alert_channels (
    id SERIAL PRIMARY KEY,
    type TEXT NOT NULL,      -- slack, telegram, email
    name TEXT NOT NULL,
    config JSONB NOT NULL,   -- {"token":"...","chat_id":"..."}
    created_at TIMESTAMPTZ DEFAULT now()
);

-- Relação N:N entre grupos e canais
CREATE TABLE alert_group_channels (
    group_id INT NOT NULL REFERENCES alert_groups(id) ON DELETE CASCADE,
    channel_id INT NOT NULL REFERENCES alert_channels(id) ON DELETE CASCADE,
    PRIMARY KEY (group_id, channel_id)
);

-- Histórico de alertas enviados
CREATE TABLE sent_alerts (
    id BIGSERIAL PRIMARY KEY,
    endpoint_id INT NOT NULL REFERENCES endpoints(id) ON DELETE CASCADE,
    channel_id INT NOT NULL REFERENCES alert_channels(id) ON DELETE CASCADE,
    message TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);


-- Endpoints monitorados
CREATE TABLE endpoints (
   id SERIAL PRIMARY KEY,
   name TEXT NOT NULL,
   domain TEXT NOT NULL,
   port INT DEFAULT 80,
   path TEXT DEFAULT '/',
   timeout_seconds INT DEFAULT 30,
   interval_seconds INT DEFAULT 300,           -- intervalo de healthcheck
   ssl_check_interval_seconds INT DEFAULT 43200, -- ex.: 12h
   check_ssl BOOLEAN DEFAULT FALSE,
   enabled BOOLEAN DEFAULT TRUE,
   alert_group_id INT NOT NULL REFERENCES alert_groups(id) ON DELETE RESTRICT,
   created_at TIMESTAMPTZ DEFAULT now(),
   updated_at TIMESTAMPTZ DEFAULT now()
);

-- Autenticação de endpoints (JSONB flexível)
CREATE TABLE endpoint_auth (
   endpoint_id INT PRIMARY KEY REFERENCES endpoints(id) ON DELETE CASCADE,
   type TEXT NOT NULL,         -- basic, bearer, api_key
   config JSONB NOT NULL,      -- {"username":"x","password":"y"} ou {"token":"..."}
   created_at TIMESTAMPTZ DEFAULT now()
);

-- Histórico das checagens (time-series)
CREATE TABLE checks (
    id BIGSERIAL PRIMARY KEY,
    endpoint_id INT NOT NULL REFERENCES endpoints(id) ON DELETE CASCADE,
    status check_status NOT NULL,
    response_time_ms INT,
    error_message TEXT,
    ssl_expiration_date TIMESTAMPTZ,
    ssl_issuer TEXT,
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE users (
   id SERIAL PRIMARY KEY,
   username TEXT NOT NULL UNIQUE,
   email TEXT NOT NULL UNIQUE,
   full_name TEXT,
   password_hash TEXT,
   auth_provider TEXT NOT NULL DEFAULT 'local',
   enabled BOOLEAN DEFAULT TRUE,
   created_at TIMESTAMPTZ DEFAULT now(),
   updated_at TIMESTAMPTZ DEFAULT now()
);