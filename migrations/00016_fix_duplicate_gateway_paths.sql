-- +goose Up
-- +goose StatementBegin

-- Remove duplicate application_versions keeping only the latest one
WITH duplicates AS (
    SELECT 
        id,
        ROW_NUMBER() OVER (
            PARTITION BY application_id, version_tag 
            ORDER BY created_at DESC
        ) as rn
    FROM application_versions 
    WHERE deleted_at IS NULL
)
DELETE FROM application_versions 
WHERE id IN (
    SELECT id FROM duplicates WHERE rn > 1
);

-- Ensure the unique constraint exists for application versions
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint 
        WHERE conname = 'unique_app_version'
    ) THEN
        ALTER TABLE application_versions
        ADD CONSTRAINT unique_app_version UNIQUE (application_id, version_tag);
    END IF;
END $$;

-- Remove duplicate gateways with same domain+path, keeping the latest one
WITH duplicates AS (
    SELECT 
        id,
        ROW_NUMBER() OVER (
            PARTITION BY domain, path 
            ORDER BY created_at DESC
        ) as rn
    FROM gateways 
    WHERE deleted_at IS NULL 
    AND path IS NOT NULL 
    AND path != ''
)
DELETE FROM gateways 
WHERE id IN (
    SELECT id FROM duplicates WHERE rn > 1
);

-- Remove duplicate gateways with same domain+subdomain, keeping the latest one  
WITH duplicates AS (
    SELECT 
        id,
        ROW_NUMBER() OVER (
            PARTITION BY domain 
            ORDER BY created_at DESC
        ) as rn
    FROM gateways 
    WHERE deleted_at IS NULL
)
DELETE FROM gateways 
WHERE id IN (
    SELECT id FROM duplicates WHERE rn > 1
);

-- Ensure the unique indexes exist for gateways (in case they were dropped)
DO $$ 
BEGIN
    -- Check and create unique index for domain+path
    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes 
        WHERE indexname = 'idx_gateways_domain_path'
    ) THEN
        CREATE UNIQUE INDEX idx_gateways_domain_path 
            ON gateways (domain, path) 
            WHERE path IS NOT NULL AND path != '' AND deleted_at IS NULL;
    END IF;
END $$;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Note: We cannot restore deleted duplicates, so this is a no-op
-- Remove the unique constraint for application versions if needed
ALTER TABLE application_versions
DROP CONSTRAINT IF EXISTS unique_app_version;
-- The unique index removal for gateways is handled by the original migration
-- +goose StatementEnd
