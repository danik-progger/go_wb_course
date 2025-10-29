-- name: GetAllMsgs :many
SELECT * FROM messages;

-- name: GetUnprocessed :many
SELECT * FROM messages
WHERE status = "not touched";

-- name: GetMsgById :one
SELECT * FROM messages
WHERE id = $1;

-- name: AddMsg :one
INSERT INTO messages (
  body, send_at, was_sent_at, status, tries
) VALUES (
  $1, $2, $2, "not touched", 0
)
RETURNING *;

-- name: Update :exec
UPDATE messages
  set status = $2,
  tries = $3,
  was_sent_at = $4
WHERE id = $1;

-- name: DeleteMsg :exec
DELETE FROM messages
WHERE id = $1;
