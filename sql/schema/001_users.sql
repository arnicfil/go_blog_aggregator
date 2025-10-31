-- +goose Up
CREATE TABLE users (
  id UUID PRIMARY KEY,
  name TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  UNIQUE (name)
);

-- +goose Down
DROP TABLE users;
