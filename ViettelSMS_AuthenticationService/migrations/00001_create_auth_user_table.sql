-- +goose Up
CREATE TABLE auth_users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    blocked BOOLEAN NOT NULL DEFAULT FALSE,
    scopes TEXT[] NOT NULL DEFAULT '{}'
);

INSERT INTO auth_users (username, password, scopes)
VALUES (
    'admin',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
    '{"user:create", "user:read", "user:update", "user:delete","user:view",
     "user:scope", "server:read", "server:update", "server:delete", "server:view", "server:import", "server:export"}'
);

-- +goose Down
DROP TABLE auth_users;
