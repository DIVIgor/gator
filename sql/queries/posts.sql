-- name: CreatePost :one
INSERT INTO posts(
    title, url, description, published_at,
    feed_id, created_at, updated_at
)
VALUES (
    $1, $2, $3, $4,
    $5, $6, $7
)
RETURNING *;

-- name: GetPostsForUser :many
SELECT posts.id, title, url, description, published_at,
    posts.feed_id, posts.created_at, posts.updated_at
FROM posts
JOIN feed_follows
ON posts.feed_id = feed_follows.feed_id
WHERE user_id = $1
ORDER BY published_at DESC
LIMIT $2;