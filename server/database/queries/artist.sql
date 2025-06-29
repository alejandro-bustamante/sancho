-- name: GetArtistByNormalizedName :one
SELECT * FROM artist
WHERE normalized_name = sqlc.arg('normalized_name')
LIMIT 1;

-- name: InsertArtist :one
INSERT INTO artist (deezer_id, name, normalized_name)
VALUES (sqlc.arg('deezer_id'), sqlc.arg('name'), sqlc.arg('normalized_name'))
RETURNING *;

-- name: ArtistExistsByDeezerID :one
SELECT EXISTS (
  SELECT 1 FROM artist WHERE deezer_id = sqlc.arg('deezer_id')
);

-- name: GetArtistByDeezerID :one
SELECT * FROM artist
WHERE deezer_id = sqlc.arg('deezer_id')
LIMIT 1;