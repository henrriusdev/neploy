-- +goose Up
-- +goose StatementBegin
alter table
  public.users
alter column
  address type text using address :: text;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
