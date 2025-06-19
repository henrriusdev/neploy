-- +goose Up
ALTER TABLE application_versions
ADD CONSTRAINT unique_app_version UNIQUE (application_id, version_tag);

-- +goose Down
ALTER TABLE application_versions
DROP CONSTRAINT unique_app_version;
