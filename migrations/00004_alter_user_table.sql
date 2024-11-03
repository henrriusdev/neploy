-- +goose Up
-- +goose StatementBegin
alter table
  public.users
alter column
  address type text using address :: text;

-- +goose StatementEnd
-- +goose StatementBegin
alter table
  public.user_oauth drop column access_token;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
alter table
  public.users
alter column
  address type varchar(255) using address :: varchar(255);

alter table
  public.user_oauth
add
  column if not exists access_token text;

-- +goose StatementEnd