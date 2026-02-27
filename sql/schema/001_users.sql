-- +goose Up
CREATE TABLE users (
    id INT PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    email VARCHAR(255) UNIQUE NOT NULL
);

-- +goose Down
DROP TABLE users;