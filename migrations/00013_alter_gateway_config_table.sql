-- +goose Up
-- +goose StatementBegin
ALTER TABLE gateway_config
    DROP COLUMN default_version,
    DROP COLUMN load_balancer;
ALTER TABLE application_stats
    DROP COLUMN IF EXISTS healthy;
ALTER TABLE visitor_traces DROP COLUMN IF EXISTS visitor_id;
ALTER TABLE visitor_traces
    ADD COLUMN device     TEXT NOT NULL DEFAULT '',
    ADD COLUMN browser    TEXT NOT NULL DEFAULT '',
    ADD COLUMN os         TEXT NOT NULL DEFAULT '',
    ADD COLUMN ip_address TEXT NOT NULL DEFAULT '';
DROP TABLE IF EXISTS visitor_info;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE gateway_config
    ADD COLUMN default_version TEXT    NOT NULL DEFAULT 'latest',
    ADD COLUMN load_balancer   BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE IF EXISTS visitor_info
    DROP COLUMN IF EXISTS device,
    DROP COLUMN IF EXISTS browser,
    DROP COLUMN IF EXISTS os,
    DROP COLUMN IF EXISTS ip_address;
-- +goose StatementEnd
