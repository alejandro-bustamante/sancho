-- name: InsertDownloadHistory :one
INSERT INTO download_history (
  id, user_id, track_id, quality, 
  status, service, completed_at, error_message
) VALUES (
  sqlc.arg('id'), sqlc.arg('user_id'), sqlc.arg('track_id'),
  sqlc.arg('quality'), sqlc.arg('status'), sqlc.arg('service'), 
  sqlc.arg('completed_at'), sqlc.arg('error_message')
)
RETURNING *;

-- name: UpdateDownloadCompletion :exec
UPDATE download_history
SET completed_at = CURRENT_TIMESTAMP, status = sqlc.arg('status'), error_message = sqlc.arg('error_message')
WHERE id = sqlc.arg('id');
