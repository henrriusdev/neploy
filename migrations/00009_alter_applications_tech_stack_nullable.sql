-- +goose Up
-- +goose StatementBegin
ALTER TABLE applications ALTER COLUMN tech_stack DROP NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE applications ALTER COLUMN tech_stack SET NOT NULL;
-- +goose StatementEnd
