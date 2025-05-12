-- +goose Up
-- +goose StatementBegin
ALTER TABLE application_stats ALTER COLUMN date SET DATA TYPE TIMESTAMP WITH TIME ZONE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE application_stats ALTER COLUMN date SET DATA TYPE DATE;
-- +goose StatementEnd
