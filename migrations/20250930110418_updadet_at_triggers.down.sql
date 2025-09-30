DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_alert_groups_updated_at ON alert_groups;
DROP TRIGGER IF EXISTS update_alert_group_channels_updated_at ON alert_group_channels;
DROP TRIGGER IF EXISTS update_alert_channels_updated_at ON alert_channels;
DROP TRIGGER IF EXISTS update_endpoint_ssl_updated_at ON endpoint_ssl;
DROP TRIGGER IF EXISTS update_endpoint_auth_updated_at ON endpoint_auth;
DROP TRIGGER IF EXISTS update_endpoints_updated_at ON endpoints;
DROP TRIGGER IF EXISTS update_endpoint_checks_updated_at ON endpoint_checks;

DROP FUNCTION IF EXISTS update_updated_at_column();

