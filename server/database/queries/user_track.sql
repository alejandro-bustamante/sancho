-- name: AddTrackToUser :exec
INSERT INTO user_track (user_id, track_id, symlink_path)
VALUES (sqlc.arg('user_id'), sqlc.arg('track_id'), sqlc.arg('symlink_path'));

-- name: GetUserTrack :one
SELECT * FROM user_track
WHERE user_id = sqlc.arg('user_id') AND track_id = sqlc.arg('track_id');
