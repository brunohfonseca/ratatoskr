-- Drop child tables first to respect foreign key dependencies
DROP TABLE IF EXISTS endpoint_ssl;
DROP TABLE IF EXISTS endpoint_auth;
DROP TABLE IF EXISTS endpoint_checks;
DROP TABLE IF EXISTS sent_alerts;
DROP TABLE IF EXISTS alert_group_channels;
DROP TABLE IF EXISTS endpoints;
DROP TABLE IF EXISTS alert_channels;
DROP TABLE IF EXISTS alert_groups;
DROP TABLE IF EXISTS users;

-- Drop custom enum types created in the up migration
DROP TYPE IF EXISTS endpoint_ssl_status;
DROP TYPE IF EXISTS alert_channel_type;
DROP TYPE IF EXISTS check_status;