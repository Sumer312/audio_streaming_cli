-- name: InsertSong :one
INSERT INTO Songs (id, title, url, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetAllSongs :many
SELECT * FROM Songs;

-- name: GetUrlByTitle :one
SELECT song.url FROM Songs song WHERE song.title = $1;

-- name: GetCachedUrlByTitle :one
SELECT cache.cached_url FROM Cache cache 
JOIN Songs song ON cache.song_id = song.id
WHERE song.id = (SELECT song.id FROM Songs song WHERE song.title = $1);

-- name: DeleteBySongId :exec
DELETE FROM Songs WHERE id = $1;
