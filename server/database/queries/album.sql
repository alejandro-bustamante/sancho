-- name: GetAlbumByNormalizedTitleAndArtist :one
SELECT * FROM album
WHERE normalized_title = sqlc.arg('normalized_title')
  AND artist_id = sqlc.arg('artist_id')
LIMIT 1;

-- name: InsertAlbum :one
INSERT INTO album (
  deezer_id, title, normalized_title, artist_id, release_date,
  album_art_path, genre, total_tracks
) VALUES (
  sqlc.arg('deezer_id'), sqlc.arg('title'), sqlc.arg('normalized_title'),
  sqlc.arg('artist_id'), sqlc.arg('release_date'), sqlc.arg('album_art_path'),
  sqlc.arg('genre'), sqlc.arg('total_tracks')
)
RETURNING *;

-- name: AlbumExistsByDeezerID :one
SELECT EXISTS (
  SELECT 1 FROM album WHERE deezer_id = sqlc.arg('deezer_id')
);

-- name: GetAlbumByDeezerID :one
SELECT * FROM album
WHERE deezer_id = sqlc.arg('deezer_id')
LIMIT 1;