-- +goose Up
-- +goose StatementBegin
ALTER TABLE roles ADD COLUMN icon TEXT NOT NULL DEFAULT '';
ALTER TABLE roles ADD COLUMN color TEXT NOT NULL DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
