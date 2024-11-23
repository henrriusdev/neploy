-- +goose Up
-- +goose StatementBegin
ALTER TABLE public.metadata
DROP COLUMN primary_color,
DROP COLUMN secondary_color;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE public.metadata
ADD COLUMN primary_color TEXT NOT NULL DEFAULT '#000000',
ADD COLUMN secondary_color TEXT NOT NULL DEFAULT '#000000';
-- +goose StatementEnd
