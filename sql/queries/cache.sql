-- name: InsertCache :one
INSERT INTO Cache (id, song_id, cached_url)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetAllSongs :many
SELECT * FROM Songs;

-- name: GetSongByTitle :one
SELECT * FROM Songs WHERE title = $1;

-- name: ClearCache :exec
TRUNCATE TABLE Cache;

