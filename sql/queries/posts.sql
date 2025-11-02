-- name: CreatePost :one
INSERT INTO
    posts (
        id,
        created_at,
        updated_at,
        name,
        url,
        description,
        published_at,
        feed_id
    )
VALUES
    ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING
    *;

-- name: GetPostsForUser :many
SELECT
    posts.name,
    posts.description,
    posts.published_at
FROM
    posts
    INNER JOIN feeds_follow ON posts.feed_id = feeds_follow.feed_id
    INNER JOIN users ON feeds_follow.user_id = users.id
WHERE
    users.id = $1
ORDER BY
    posts.published_at DESC
LIMIT
    $2;
