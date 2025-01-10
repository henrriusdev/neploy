-- +goose Up
-- +goose StatementBegin
ALTER TABLE applications ALTER COLUMN tech_stack_id DROP NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE applications ALTER COLUMN tech_stack_id SET NOT NULL;
-- +goose StatementEnd
