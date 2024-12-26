-- +goose Up
-- +goose StatementBegin
-- Add new columns
ALTER TABLE gateways
    ADD COLUMN domain text DEFAULT '',
    ADD COLUMN subdomain text DEFAULT '',
    ADD COLUMN path text DEFAULT '',
    ADD COLUMN port text DEFAULT '',
    ADD COLUMN status text DEFAULT 'inactive';

-- Add constraints
ALTER TABLE gateways
    ALTER COLUMN endpoint_type SET NOT NULL,
    ALTER COLUMN endpoint_type TYPE text,
    ADD CONSTRAINT check_endpoint_type CHECK (endpoint_type IN ('subdomain', 'path')),
    ADD CONSTRAINT check_status CHECK (status IN ('active', 'inactive', 'error'));

-- Add indexes
CREATE INDEX idx_gateways_domain ON gateways (domain);
CREATE INDEX idx_gateways_endpoint_type ON gateways (endpoint_type);
CREATE INDEX idx_gateways_status ON gateways (status);

-- Add unique constraints that respect soft deletes
CREATE UNIQUE INDEX idx_gateways_domain_subdomain 
    ON gateways (domain, subdomain) 
    WHERE subdomain IS NOT NULL AND deleted_at IS NULL;

CREATE UNIQUE INDEX idx_gateways_domain_path 
    ON gateways (domain, path) 
    WHERE path IS NOT NULL AND deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove indexes
DROP INDEX IF EXISTS idx_gateways_domain_path;
DROP INDEX IF EXISTS idx_gateways_domain_subdomain;
DROP INDEX IF EXISTS idx_gateways_status;
DROP INDEX IF EXISTS idx_gateways_endpoint_type;
DROP INDEX IF EXISTS idx_gateways_domain;

-- Remove constraints
ALTER TABLE gateways
    DROP CONSTRAINT IF EXISTS check_endpoint_type,
    DROP CONSTRAINT IF EXISTS check_status;

-- Remove columns
ALTER TABLE gateways
    DROP COLUMN IF EXISTS domain,
    DROP COLUMN IF EXISTS subdomain,
    DROP COLUMN IF EXISTS path,
    DROP COLUMN IF EXISTS port,
    DROP COLUMN IF EXISTS status;
-- +goose StatementEnd
