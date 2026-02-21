-- name: AddTrackToUser :exec
INSERT INTO user_track (user_id, track_id, symlink_path)
VALUES (sqlc.arg('user_id'), sqlc.arg('track_id'), sqlc.arg('symlink_path'));

-- name: GetUserTrack :one
SELECT * FROM user_track
WHERE user_id = sqlc.arg('user_id') AND track_id = sqlc.arg('track_id');

-- name: IsTrackLinkedToUserByUsernameAndISRC :one
SELECT EXISTS (
  SELECT 1
  FROM user_track ut
  JOIN track t ON ut.track_id = t.id
  JOIN user u ON ut.user_id = u.id
  WHERE u.username = sqlc.arg('username') AND t.isrc = sqlc.arg('isrc')
);

-- name: DeleteUserTrack :exec
DELETE FROM user_track
WHERE user_id = sqlc.arg('user_id') AND track_id = sqlc.arg('track_id');

-- name: CountUsersForTrack :one
SELECT COUNT(*) FROM user_track
WHERE track_id = sqlc.arg('track_id');

-- name: ListTracksByUsername :many
SELECT
  t.id,
  t.title,
  t.duration,
  art.name AS artist,
  alb.title AS album,
  alb.album_art_path
FROM track AS t
JOIN user_track AS ut ON t.id = ut.track_id
JOIN "user" AS u ON ut.user_id = u.id
LEFT JOIN artist AS art ON t.artist_id = art.id
LEFT JOIN album AS alb ON t.album_id = alb.id
WHERE u.username = ?
ORDER BY art.name, alb.title, t.track_number;