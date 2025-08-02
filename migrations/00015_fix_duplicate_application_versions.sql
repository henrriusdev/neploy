-- +goose Up
-- +goose StatementBegin

-- First, remove duplicate application_versions keeping only the latest one
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

-- Ensure the unique constraint exists (in case the previous migration failed)
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

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Note: We cannot restore deleted duplicates, so this is a no-op
-- The unique constraint removal is handled by the original migration
-- +goose StatementEnd
