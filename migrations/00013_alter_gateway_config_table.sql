-- +goose Up
-- +goose StatementBegin
ALTER TABLE gateway_config
    DROP COLUMN default_version,
    DROP COLUMN load_balancer;
ALTER TABLE application_stats DROP COLUMN healthy;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE gateway_config
    ADD COLUMN default_version TEXT    NOT NULL DEFAULT 'latest',
    ADD COLUMN load_balancer   BOOLEAN NOT NULL DEFAULT FALSE;
-- +goose StatementEnd
