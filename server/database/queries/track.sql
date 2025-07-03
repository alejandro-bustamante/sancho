-- name: InsertTrack :one
INSERT INTO track (
  title, normalized_title, artist_id, album_id, duration, track_number, disc_number,
  sample_rate, bitrate, channels, file_path, file_size, isrc, composer
) VALUES (
  sqlc.arg('title'), sqlc.arg('normalized_title'), sqlc.arg('artist_id'), 
  sqlc.arg('album_id'), sqlc.arg('duration'), sqlc.arg('track_number'),
  sqlc.arg('disc_number'), sqlc.arg('sample_rate'), sqlc.arg('bitrate'),
  sqlc.arg('channels'), sqlc.arg('file_path'), sqlc.arg('file_size'),
  sqlc.arg('isrc'), sqlc.arg('composer')
)
RETURNING *;

-- name: SearchTracksByTitle :many
SELECT * FROM track
WHERE LOWER(title) LIKE LOWER('%' || sqlc.arg('title') || '%')
ORDER BY title DESC;

-- name: SearchTracksByISRC :one
SELECT * FROM track
WHERE isrc = sqlc.arg('isrc')
LIMIT 1;

-- name: ListTracksByDate :many
SELECT * FROM track ORDER BY created_at;

-- name: TrackExistsByISRC :one
SELECT EXISTS (
  SELECT 1 FROM track WHERE isrc = sqlc.arg('isrc')
);

-- name: GetArtistByTrackID :one
SELECT artist.* FROM track
JOIN artist ON track.artist_id = artist.id
WHERE track.id = sqlc.arg('track_id');

-- name: GetAlbumByTrackID :one
SELECT album.* FROM track
JOIN album ON track.album_id = album.id
WHERE track.id = sqlc.arg('track_id');

-- name: UpdateTrackFilePath :exec
UPDATE track
SET file_path = sqlc.arg('file_path')
WHERE id = sqlc.arg('track_id');

