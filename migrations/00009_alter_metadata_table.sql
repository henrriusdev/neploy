-- +goose Up
-- +goose StatementBegin
ALTER TABLE metadata ADD COLUMN language TEXT NOT NULL DEFAULT 'en';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE metadata DROP COLUMN language;
-- +goose StatementEnd
