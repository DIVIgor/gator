-- name: CreateFeed :one
INSERT INTO feeds (name, url, user_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetFeeds :many
SELECT f.id, f.name, url, u.name AS username, f.created_at, f.updated_at
FROM feeds f
JOIN users u
ON f.user_id = u.id;

-- name: GetFeed :one
SELECT *
FROM feeds
WHERE url = $1;