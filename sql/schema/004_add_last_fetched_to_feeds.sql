-- +goose Up
-- +goose StatementBegin
ALTER TABLE feeds ADD last_fetched_at TIMESTAMP;

-- +goose StatementEnd
-- +goose Down
ALTER TABLE feeds
DROP COLUMN last_fetched_at;

-- +goose StatementBegin
-- +goose StatementEnd
