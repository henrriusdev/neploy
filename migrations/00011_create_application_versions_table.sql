-- +goose Up
-- +goose StatementBegin
ALTER TABLE applications
    DROP COLUMN deploy_location;
CREATE TABLE application_versions
(
    id               UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    created_at       TIMESTAMP WITH TIME ZONE  DEFAULT NOW(),
    updated_at       TIMESTAMP WITH TIME ZONE  DEFAULT NOW(),
    deleted_at       TIMESTAMP WITH TIME ZONE,

    version_tag      VARCHAR(50)      NOT NULL,
    description      TEXT,
    status           VARCHAR(20)               DEFAULT 'inactive',
    storage_location TEXT             NOT NULL,
    application_id   UUID             NOT NULL,

    CONSTRAINT fk_application
        FOREIGN KEY (application_id)
            REFERENCES applications (id)
            ON DELETE CASCADE
);
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE application_stats ADD COLUMN healthy BOOLEAN DEFAULT true;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE applications
    ADD COLUMN deploy_location TEXT;
DROP TABLE application_versions;
-- +goose StatementEnd
