-- +goose Up
-- +goose StatementBegin
CREATE TABLE invitations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    role UUID NOT NULL REFERENCES roles(id),
    token VARCHAR(100) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    accepted_at TIMESTAMP,
    deleted_at TIMESTAMP
);
CREATE INDEX idx_invitations_token ON invitations (token);
CREATE INDEX idx_invitations_user_id ON invitations (user_id);
CREATE INDEX idx_invitations_role ON invitations (role);
CREATE UNIQUE INDEX idx_invitations_user_id_role ON invitations (user_id, role);
CREATE TRIGGER update_invitations_updated_at BEFORE UPDATE ON public.invitations FOR EACH ROW execute function update_updated_at_column (); 
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE invitations;
-- +goose StatementEnd
