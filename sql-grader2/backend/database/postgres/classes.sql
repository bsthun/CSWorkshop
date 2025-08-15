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

-- name: StudentClassList :many
SELECT sqlc.embed(classes),
       sqlc.embed(semesters),
       COALESCE(COUNT(exam_attempts.id), 0)::bigint as exam_attempt_total_count,
       COALESCE(COUNT(CASE WHEN exam_attempts.finished_at IS NOT NULL THEN 1 END), 0)::bigint as exam_attempt_finished_count
FROM class_joinees
JOIN classes ON class_joinees.class_id = classes.id
JOIN semesters ON classes.semester_id = semesters.id
LEFT JOIN exam_attempts ON class_joinees.id = exam_attempts.class_joinee_id
WHERE class_joinees.user_id = $1
GROUP BY classes.id, semesters.id, class_joinees.id
ORDER BY class_joinees.created_at DESC;

-- name: ClassGetByRegisterCode :one
SELECT *
FROM classes
WHERE register_code = $1;

-- name: ClassJoineeCreate :one
INSERT INTO class_joinees (user_id, class_id)
VALUES ($1, $2)
RETURNING *;

-- name: ClassJoineeGetByUserAndClass :one
SELECT *
FROM class_joinees
WHERE user_id = $1 AND class_id = $2;