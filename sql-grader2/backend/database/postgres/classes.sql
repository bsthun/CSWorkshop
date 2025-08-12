-- name: ClassCreate :one
INSERT INTO classes (semester_id, code, name, register_code)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ClassGetById :one
SELECT sqlc.embed(classes), sqlc.embed(semesters)
FROM classes
JOIN semesters ON classes.semester_id = semesters.id
WHERE classes.id = $1;

-- name: ClassJoineeList :many
SELECT sqlc.embed(class_joinees), sqlc.embed(users)
FROM class_joinees
JOIN users ON class_joinees.user_id = users.id
WHERE class_joinees.class_id = $1
ORDER BY class_joinees.created_at DESC;

-- name: ClassUpdate :one
UPDATE classes
SET code = $2, name = $3, register_code = $4
WHERE id = $1
RETURNING *;