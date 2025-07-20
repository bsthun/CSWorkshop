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

-- name: UserCreate :one
INSERT INTO users (oid, firstname, lastname, email, picture_url)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;
