-- +goose Up
-- +goose StatementBegin
-- Add provider column to users table
ALTER TABLE public.users
ADD COLUMN provider text;

-- Update users with provider data from user_oauth table
UPDATE public.users u
SET provider = uo.provider
FROM public.user_oauth uo
WHERE u.id = uo.user_id;

-- Drop user_oauth table and related indexes
DROP INDEX IF EXISTS idx_user_oauth_user_id;
DROP TABLE IF EXISTS public.user_oauth CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Recreate user_oauth table
CREATE TABLE public.user_oauth (
  id uuid primary key default gen_random_uuid(),
  user_id uuid references public.users (id),
  provider text not null,
  oauth_id text not null,
  access_token text,
  created_at timestamptz default current_timestamp,
  updated_at timestamptz default current_timestamp,
  deleted_at timestamptz default null
);

-- Create trigger for updated_at
CREATE TRIGGER update_user_oauth_updated_at BEFORE
UPDATE ON public.user_oauth FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Create index
CREATE INDEX idx_user_oauth_user_id ON public.user_oauth USING btree (user_id);

-- Move provider data from users table to user_oauth table
-- Note: This is a best-effort migration, as oauth_id will be set to 'unknown'
-- in the rollback scenario, which may not be ideal
INSERT INTO public.user_oauth (user_id, provider, oauth_id)
SELECT id, provider, 'unknown'
FROM public.users
WHERE provider IS NOT NULL;

-- Remove provider column from users table
ALTER TABLE public.users
DROP COLUMN provider;
-- +goose StatementEnd
