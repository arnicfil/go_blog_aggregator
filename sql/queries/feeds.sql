-- name: CreateFeed :one
INSERT INTO feeds (id,name, created_at, updated_at, url, user_id)
VALUES (
	$1,
	$2,
	$3,
	$4,
	$5,
	$6
)
returning *;

-- name: ListFeeds :many
SELECT * FROM feeds;

-- name: RetrieveFeedUser :one
SELECT users.name FROM users
WHERE users.id IN (
	SELECT feeds.user_id FROM feeds
	WHERE feeds.name = $1
);
