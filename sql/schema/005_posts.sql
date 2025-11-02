-- +goose Up
-- +goose StatementBegin
CREATE TABLE posts (
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  name TEXT NOT NULL,
  url TEXT NOT NULL,
  description TEXT,
  published_at TIMESTAMP,
  feed_id UUID NOT NULL,
  unique (url)
);

-- +goose StatementEnd
-- +goose Down
DROP TABLE posts;

-- +goose StatementBegin
-- +goose StatementEnd
