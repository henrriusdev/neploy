-- +goose Up
-- +goose StatementBegin
ALTER TABLE public.application_stats DROP COLUMN environment_id;
DROP TABLE public.environments CASCADE;
-- +goose StatementEnd

-- +goose Down
