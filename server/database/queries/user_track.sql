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