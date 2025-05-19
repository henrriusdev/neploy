-- +goose Up
-- +goose StatementBegin
ALTER TABLE gateway_config
    DROP COLUMN default_version,
    DROP COLUMN load_balancer;
ALTER TABLE application_stats
    DROP COLUMN IF EXISTS healthy;
ALTER TABLE visitor_traces
    DROP COLUMN IF EXISTS visitor_id;
ALTER TABLE visitor_traces
    ADD COLUMN IF NOT EXISTS device     TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS browser    TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS os         TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS ip_address TEXT NOT NULL DEFAULT '';
DROP TABLE IF EXISTS visitor_info;
DROP TABLE IF EXISTS refresh_tokens;
ALTER TABLE gateways
    DROP COLUMN IF EXISTS endpoint_type,
    DROP COLUMN IF EXISTS stage,
    DROP COLUMN IF EXISTS name,
    DROP COLUMN IF EXISTS integration_type,
    DROP COLUMN IF EXISTS subdomain,
    drop column if EXISTS http_method,
    drop column if EXISTS endpoint_url,
    drop column if EXISTS logging_level;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
    ALTER table gateway_config
    ADD COLUMN default_version TEXT NOT NULL DEFAULT 'latest',
    ADD COLUMN load_balancer BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE IF EXISTS visitor_info
    DROP COLUMN IF EXISTS device,
    DROP COLUMN IF EXISTS browser,
    DROP COLUMN IF EXISTS os,
    DROP COLUMN IF EXISTS ip_address;
-- +goose StatementEnd
