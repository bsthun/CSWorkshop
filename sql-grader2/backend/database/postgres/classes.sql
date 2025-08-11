-- name: ClassCreate :one
INSERT INTO classes (semester_id, code, name, register_code)
VALUES ($1, $2, $3, $4)
RETURNING *;