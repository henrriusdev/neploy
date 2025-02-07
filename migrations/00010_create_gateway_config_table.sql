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

-- +goose StatementBegin
CREATE TABLE gateway_versions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    version_name    TEXT NOT NULL DEFAULT 'v1',
    status          TEXT NOT NULL DEFAULT 'active',
    versioning_type TEXT NOT NULL DEFAULT 'uri',
    gateway_id      UUID NOT NULL REFERENCES public.gateways,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMP DEFAULT NULL
);
CREATE INDEX idx_gateway_version_gateway_id ON gateway_versions(gateway_id);
CREATE TRIGGER update_gateway_versions_updated_at BEFORE UPDATE ON public.gateway_versions FOR EACH ROW execute function update_updated_at_column ();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE gateway_config;
DROP TABLE gateway_version;
-- +goose StatementEnd
