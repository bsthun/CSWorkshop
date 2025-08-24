-- name: UserGetById :one
SELECT *
FROM users
WHERE id = $1;

-- name: UserGetByOid :one
SELECT *
FROM users
WHERE oid = $1;

-- name: UserList :many
SELECT *
FROM users
ORDER BY created_at DESC;

-- name: UserGetByMetadataUsername :one
SELECT *
FROM users
WHERE metadata->'credential'->>'username' = sqlc.narg(username)::text;

-- name: UserCreate :one
INSERT INTO users (oid, firstname, lastname, email, picture_url, is_admin, metadata)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UserUpdateMetadata :one
UPDATE users
SET metadata = $2
WHERE id = $1
RETURNING *;
