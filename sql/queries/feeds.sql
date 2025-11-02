-- name: CreateFeed :one
INSERT INTO
    feeds (id, name, created_at, updated_at, url, user_id)
VALUES
    ($1, $2, $3, $4, $5, $6)
RETURNING
    *;

-- name: ListFeeds :many
SELECT
    *
FROM
    feeds;

-- name: FeedFromUrl :one
SELECT
    *
FROM
    feeds
WHERE
    url = $1;

-- name: MarkFeedFetched :exec
UPDATE
    feeds
SET
    updated_at = $2,
    last_fetched_at = $3
WHERE
    feeds.id = $1;

-- name: GetNextFeedToFetch :one
SELECT
    *
FROM
    feeds
ORDER BY
    last_fetched_at ASC NULLS FIRST,
    updated_AT ASC
LIMIT
    1;
