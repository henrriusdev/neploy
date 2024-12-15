-- +goose Up
-- +goose StatementBegin
CREATE TABLE invitations (
    id VARCHAR(36) PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    team_id VARCHAR(36) NOT NULL,
    role VARCHAR(50) NOT NULL,
    token VARCHAR(100) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    accepted_at TIMESTAMP,
    FOREIGN KEY (team_id) REFERENCES teams(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE invitations;
-- +goose StatementEnd
