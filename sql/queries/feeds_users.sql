-- name: RetrieveFeedUser :one
SELECT users.name FROM users
WHERE users.id IN (
	SELECT feeds.user_id FROM feeds
	WHERE feeds.name = $1
);

-- name: CreateFeedFollow :one
with inserted_feed_follow as (
	INSERT INTO feeds_follow(id, created_at, updated_at, user_id, feed_id)
	VALUES (
		$1,
		$2,
		$3,
		$4,
		$5
	)
	returning *
)
SELECT
	inserted_feed_follow.*,
	feeds.name AS feed_name,
	users.name AS user_ame
FROM
	inserted_feed_follow
INNER JOIN
	users ON inserted_feed_follow.user_id = users.id
INNER JOIN
	feeds ON inserted_feed_follow.feed_id = feeds.id;

-- name: GetFeedFollowsForUser :many
SELECT users.name, feeds.name
FROM
	users
INNER JOIN
	feeds_follow on users.id = feeds_follow.user_id
INNER JOIN
	feeds on feeds.id = feeds_follow.feed_id
where users.name = $1;

-- name: DeleteFeedFollow :exec
DELETE FROM
  feeds_follow
WHERE
  id IN (
    SELECT
      feeds_follow.id
    FROM
      feeds_follow
    INNER JOIN
      feeds ON feeds_follow.feed_id = feeds.id
    INNER JOIN
      users ON feeds_follow.user_id = users.id
    WHERE
      users.id = $1 AND feeds.id = $2
  );
