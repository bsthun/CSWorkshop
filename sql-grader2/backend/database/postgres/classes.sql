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
SELECT sqlc.embed(class_joinees), 
       sqlc.embed(users),
       (SELECT COUNT(*)::BIGINT FROM exam_attempts WHERE exam_attempts.class_joinee_id = class_joinees.id) AS exam_attempt_count
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
       (SELECT COUNT(*)::BIGINT FROM exams WHERE exams.class_id = classes.id) AS exam_total_count,
       (SELECT COUNT(*)::BIGINT FROM exam_attempts WHERE exam_attempts.class_joinee_id = class_joinees.id AND exam_attempts.finished_at IS NOT NULL) AS exam_finished_count
FROM class_joinees
JOIN classes ON class_joinees.class_id = classes.id
JOIN semesters ON classes.semester_id = semesters.id
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

-- name: ClassJoineeGetById :one
SELECT *
FROM class_joinees
WHERE id = $1;