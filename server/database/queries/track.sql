-- name: InsertTrack :one
INSERT INTO track (
  title, normalized_title, artist_id, album_id, duration, track_number, disc_number,
  sample_rate, bit_depth, bitrate, channels, codec,
  file_path, file_size, isrc
) VALUES (
  sqlc.arg('title'), sqlc.arg('normalized_title'), sqlc.arg('artist_id'), 
  sqlc.arg('album_id'), sqlc.arg('duration'), sqlc.arg('track_number'),
  sqlc.arg('disc_number'), sqlc.arg('sample_rate'), sqlc.arg('bit_depth'),
  sqlc.arg('bitrate'), sqlc.arg('channels'), sqlc.arg('codec'),
  sqlc.arg('file_path'), sqlc.arg('file_size'), sqlc.arg('isrc')
)
RETURNING *;

-- name: SearchTracksByTitle :many
SELECT * FROM track
WHERE LOWER(title) LIKE LOWER('%' || sqlc.arg('title') || '%')
ORDER BY title DESC;

-- name: ListTracksByDate :many
SELECT * FROM track ORDER BY created_at;


