-- +goose Up
-- +goose StatementBegin
CREATE TABLE gateway_config (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    default_versioning_type TEXT NOT NULL DEFAULT 'uri',
    default_version         TEXT NOT NULL DEFAULT 'latest',
    load_balancer           BOOLEAN NOT NULL DEFAULT FALSE,
    created_at              TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at              TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at              TIMESTAMP DEFAULT NULL
);
CREATE TRIGGER update_gateway_config_updated_at BEFORE UPDATE ON public.gateway_config FOR EACH ROW execute function update_updated_at_column ();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE gateway_config;
-- +goose StatementEnd
