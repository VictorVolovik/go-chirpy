-- name: CreateRefreshToken :exec
INSERT INTO refresh_tokens (
    token, created_at, updated_at, user_id, expires_at
)
VALUES (
    $1,
    now(),
    now(),
    $2,
    $3
);

-- name: GetUserFromRefreshToken :one
SELECT users.id
FROM users
INNER JOIN refresh_tokens
    ON users.id = refresh_tokens.user_id
WHERE
    refresh_tokens.token = $1
    AND refresh_tokens.revoked_at IS NULL
    AND refresh_tokens.expires_at > now();


-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET
    updated_at = now(),
    revoked_at = now()
WHERE
    token = $1;
