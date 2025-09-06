-- name: SemesterList :many
SELECT *
FROM semesters
ORDER BY name DESC;

-- name: SemesterCount :one
SELECT COUNT(*)
FROM semesters;

-- name: SemesterClassList :many
SELECT 
    classes.*,
    COUNT(class_joinees.id) as joinee_count
FROM classes
LEFT JOIN class_joinees ON classes.id = class_joinees.class_id
WHERE classes.semester_id = $1
GROUP BY classes.id, classes.semester_id, classes.code, classes.name, classes.register_code, classes.created_at, classes.updated_at
ORDER BY classes.created_at DESC;

-- name: SemesterCreate :one
INSERT INTO semesters (name)
VALUES ($1)
RETURNING *;

-- name: SemesterGetById :one
SELECT *
FROM semesters
WHERE id = $1;

-- name: SemesterUpdate :one
UPDATE semesters
SET name = $2
WHERE id = $1
RETURNING *;