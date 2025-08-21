-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid (),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING id, created_at, updated_at, email;

-- name: ResetUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1;

-- name: GetUserFromRefreshToken :one
SELECT *
FROM users
WHERE users.id IN (
    SELECT refresh_tokens.user_id
    FROM refresh_tokens
    WHERE token = $1
);

-- name: UpdateEmailAndPassword :one
UPDATE users
SET email = $1, hashed_password = $2, updated_at = NOW()
WHERE id = $3
RETURNING id, created_at, updated_at, email;

-- name: GetUserByChirp :one
SELECT *
FROM users
WHERE users.id IN (
    SELECT user_id
    FROM chirps
    WHERE chirps.id = $1
);
