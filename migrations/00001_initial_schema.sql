-- +goose Up
-- +goose StatementBegin
create
or replace function update_updated_at_column () returns trigger as $$
BEGIN
  NEW.updated_at = CURRENT_TIMESTAMP;
  RETURN NEW;
END;
$$ language plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE EXTENSION pgcrypto;
-- +goose StatementEnd

-- +goose StatementBegin
create table public.users (
  id uuid primary key default gen_random_uuid (),
  username text not null unique,
  email text not null unique,
  password_hash text not null,
  first_name text,
  last_name text,
  date_of_birth date,
  address text,
  phone text,
  created_at timestamptz default current_timestamp,
  updated_at timestamptz default current_timestamp,
  deleted_at timestamptz default null
);

create trigger update_users_updated_at before
update on public.users for each row
execute function update_updated_at_column ();
-- +goose StatementEnd

-- +goose StatementBegin
create table public.tech_stacks (
  id uuid primary key default gen_random_uuid (),
  name text not null unique,
  description text,
  created_at timestamptz default current_timestamp,
  updated_at timestamptz default current_timestamp,
  deleted_at timestamptz default null
);
-- +goose StatementEnd

create trigger update_tech_stacks_updated_at before
update on public.tech_stacks for each row
execute function update_updated_at_column ();

create table public.roles (
  id uuid primary key default gen_random_uuid (),
  name text not null unique,
  description text,
  created_at timestamptz default current_timestamp,
  updated_at timestamptz default current_timestamp,
  deleted_at timestamptz default null
);

create trigger update_roles_updated_at before
update on public.roles for each row
execute function update_updated_at_column ();

create table public.user_tech_stacks (
  user_id uuid references public.users (id),
  tech_stack_id uuid references public.tech_stacks (id),
  created_at timestamptz default current_timestamp,
  updated_at timestamptz default current_timestamp,
  deleted_at timestamptz default null,
  primary key (user_id, tech_stack_id)
);

create trigger update_user_tech_stacks_updated_at before
update on public.user_tech_stacks for each row
execute function update_updated_at_column ();

create table public.user_roles (
  user_id uuid references public.users (id),
  role_id uuid references public.roles (id),
  created_at timestamptz default current_timestamp,
  updated_at timestamptz default current_timestamp,
  deleted_at timestamptz default null,
  primary key (user_id, role_id)
);

create trigger update_user_roles_updated_at before
update on public.user_roles for each row
execute function update_updated_at_column ();

create table public.applications (
  id uuid primary key default gen_random_uuid (),
  app_name text not null unique,
  storage_location text not null,
  tech_stack_id uuid references public.tech_stacks (id),
  description text,
  created_at timestamptz default current_timestamp,
  updated_at timestamptz default current_timestamp,
  deleted_at timestamptz default null
);

create trigger update_applications_updated_at before
update on public.applications for each row
execute function update_updated_at_column ();

create table public.traces (
  id uuid primary key default gen_random_uuid (),
  user_id uuid references public.users (id),
  action text not null,
  action_timestamp timestamptz default current_timestamp,
  sql_statement text,
  created_at timestamptz default current_timestamp,
  updated_at timestamptz default current_timestamp,
  deleted_at timestamptz default null
);

create trigger update_traces_updated_at before
update on public.traces for each row
execute function update_updated_at_column ();

create table public.gateways (
  id uuid primary key default gen_random_uuid (),
  gateway_name text not null unique,
  endpoint_url text not null,
  endpoint_type text,
  stage text,
  http_method text,
  integration_type text,
  logging_level text,
  application_id uuid references public.applications (id),
  created_at timestamptz default current_timestamp,
  updated_at timestamptz default current_timestamp,
  deleted_at timestamptz default null
);

create trigger update_gateways_updated_at before
update on public.gateways for each row
execute function update_updated_at_column ();

create table public.environments (
  id uuid primary key default gen_random_uuid (),
  name text not null unique,
  created_at timestamptz default current_timestamp,
  updated_at timestamptz default current_timestamp,
  deleted_at timestamptz default null
);

create trigger update_environments_updated_at before
update on public.environments for each row
execute function update_updated_at_column ();

create table public.application_environments (
  application_id uuid references public.applications (id),
  environment_id uuid references public.environments (id),
  created_at timestamptz default current_timestamp,
  updated_at timestamptz default current_timestamp,
  deleted_at timestamptz default null,
  primary key (application_id, environment_id)
);

create trigger update_application_environments_updated_at before
update on public.application_environments for each row
execute function update_updated_at_column ();

create table public.refresh_tokens (
  id uuid primary key default gen_random_uuid (),
  user_id uuid references public.users (id),
  token text not null unique,
  expires_at timestamptz not null,
  created_at timestamptz default current_timestamp,
  updated_at timestamptz default current_timestamp,
  deleted_at timestamptz default null
);

create trigger update_refresh_tokens_updated_at before
update on public.refresh_tokens for each row
execute function update_updated_at_column ();

create table public.application_stats (
  id uuid primary key default gen_random_uuid (),
  application_id uuid references public.applications (id),
  date date not null,
  requests bigint default 0,
  errors bigint default 0,
  average_response_time numeric default 0,
  data_transferred numeric default 0,
  unique_visitors bigint default 0,
  created_at timestamptz default current_timestamp,
  updated_at timestamptz default current_timestamp,
  deleted_at timestamptz default null
);

create trigger update_application_stats_updated_at before
update on public.application_stats for each row
execute function update_updated_at_column ();

create table public.visitor_info (
  id uuid primary key default gen_random_uuid (),
  ip_address text not null,
  location text,
  visit_timestamp timestamptz default current_timestamp,
  created_at timestamptz default current_timestamp,
  updated_at timestamptz default current_timestamp,
  deleted_at timestamptz default null
);

create trigger update_visitor_info_updated_at before
update on public.visitor_info for each row
execute function update_updated_at_column ();

create table public.visitor_trace (
  id uuid primary key default gen_random_uuid (),
  visitor_id uuid references public.visitor_info (id),
  page_visited text not null,
  visit_duration numeric,
  visit_timestamp timestamptz default current_timestamp,
  application_id uuid references public.applications (id),
  created_at timestamptz default current_timestamp,
  updated_at timestamptz default current_timestamp,
  deleted_at timestamptz default null
);

create trigger update_visitor_trace_updated_at before
update on public.visitor_trace for each row
execute function update_updated_at_column ();

create table public.user_oauth (
  id uuid primary key default gen_random_uuid (),
  user_id uuid references public.users (id),
  provider text not null,
  oauth_id text not null,
  access_token text,
  created_at timestamptz default current_timestamp,
  updated_at timestamptz default current_timestamp,
  deleted_at timestamptz default null
);

create trigger update_user_oauth_updated_at before
update on public.user_oauth for each row
execute function update_updated_at_column ();

drop trigger if exists update_users_updated_at on public.users;

drop trigger if exists update_tech_stacks_updated_at on public.tech_stacks;

drop trigger if exists update_roles_updated_at on public.roles;

drop trigger if exists update_user_tech_stacks_updated_at on public.user_tech_stacks;

drop trigger if exists update_user_roles_updated_at on public.user_roles;

drop trigger if exists update_applications_updated_at on public.applications;

drop trigger if exists update_traces_updated_at on public.traces;

drop trigger if exists update_gateways_updated_at on public.gateways;

drop trigger if exists update_environments_updated_at on public.environments;

drop trigger if exists update_application_environments_updated_at on public.application_environments;

drop trigger if exists update_refresh_tokens_updated_at on public.refresh_tokens;

drop trigger if exists update_application_stats_updated_at on public.application_stats;

drop trigger if exists update_visitor_info_updated_at on public.visitor_info;

drop trigger if exists update_visitor_trace_updated_at on public.visitor_trace;

drop trigger if exists update_user_oauth_updated_at on public.user_oauth;

create trigger update_users_updated_at before
update on public.users for each row
execute function update_updated_at_column ();

create trigger update_tech_stacks_updated_at before
update on public.tech_stacks for each row
execute function update_updated_at_column ();

create trigger update_roles_updated_at before
update on public.roles for each row
execute function update_updated_at_column ();

create trigger update_user_tech_stacks_updated_at before
update on public.user_tech_stacks for each row
execute function update_updated_at_column ();

create trigger update_user_roles_updated_at before
update on public.user_roles for each row
execute function update_updated_at_column ();

create trigger update_applications_updated_at before
update on public.applications for each row
execute function update_updated_at_column ();

create trigger update_traces_updated_at before
update on public.traces for each row
execute function update_updated_at_column ();

create trigger update_gateways_updated_at before
update on public.gateways for each row
execute function update_updated_at_column ();

create trigger update_environments_updated_at before
update on public.environments for each row
execute function update_updated_at_column ();

create trigger update_application_environments_updated_at before
update on public.application_environments for each row
execute function update_updated_at_column ();

create trigger update_refresh_tokens_updated_at before
update on public.refresh_tokens for each row
execute function update_updated_at_column ();

create trigger update_application_stats_updated_at before
update on public.application_stats for each row
execute function update_updated_at_column ();

create trigger update_visitor_info_updated_at before
update on public.visitor_info for each row
execute function update_updated_at_column ();

create trigger update_visitor_trace_updated_at before
update on public.visitor_trace for each row
execute function update_updated_at_column ();

create trigger update_user_oauth_updated_at before
update on public.user_oauth for each row
execute function update_updated_at_column ();

create index idx_user_tech_stacks_user_id on public.user_tech_stacks using btree (user_id);

create index idx_user_tech_stacks_tech_stack_id on public.user_tech_stacks using btree (tech_stack_id);

create index idx_user_roles_user_id on public.user_roles using btree (user_id);

create index idx_user_roles_role_id on public.user_roles using btree (role_id);

create index idx_applications_tech_stack_id on public.applications using btree (tech_stack_id);

create index idx_traces_user_id on public.traces using btree (user_id);

create index idx_gateways_application_id on public.gateways using btree (application_id);

create index idx_application_environments_application_id on public.application_environments using btree (application_id);

create index idx_application_environments_environment_id on public.application_environments using btree (environment_id);

create index idx_refresh_tokens_user_id on public.refresh_tokens using btree (user_id);

create index idx_application_stats_application_id on public.application_stats using btree (application_id);

create index idx_visitor_trace_visitor_id on public.visitor_trace using btree (visitor_id);

create index idx_visitor_trace_application_id on public.visitor_trace using btree (application_id);

create index idx_user_oauth_user_id on public.user_oauth using btree (user_id);

alter table public.users
add constraint unique_phone unique (phone);

alter table public.users
rename column date_of_birth to dob;

alter table public.users
rename column password_hash to password;

alter table public.users
alter column address type jsonb using address::jsonb;

alter table public.applications
alter column tech_stack_id
set not null;

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
