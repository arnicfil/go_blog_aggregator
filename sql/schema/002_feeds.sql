-- +goose Up
-- +goose StatementBegin
CREATE TABLE feeds (
  id UUID PRIMARY KEY,
  name TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  url TEXT NOT NULL,
  user_id UUID NOT NULL,
  unique (url),
  FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE feeds
-- +goose StatementEnd
