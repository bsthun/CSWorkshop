-- name: ExamCreate :one
INSERT INTO exams (class_id, collection_id, name, access_code, opened_at, closed_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: ExamList :many
SELECT sqlc.embed(exams),
       sqlc.embed(collections),
       COALESCE(COUNT(exam_questions.id), 0)::int as exam_question_count,
       COALESCE(COUNT(collection_questions.id), 0)::int as collection_question_count
FROM exams
JOIN collections ON exams.collection_id = collections.id
LEFT JOIN collection_questions ON collections.id = collection_questions.collection_id
LEFT JOIN exam_questions ON exams.id = exam_questions.exam_id
WHERE exams.class_id = $1
GROUP BY exams.id, collections.id
ORDER BY exams.created_at DESC;

-- name: ExamJoineeList :many
SELECT sqlc.embed(exam_attempts), sqlc.embed(class_joinees), sqlc.embed(users),
       COUNT(CASE WHEN es.check_query_passed = true AND es.check_prompt_passed = true THEN 1 END) as passed_count,
       COUNT(CASE WHEN es.check_query_passed = true AND es.check_prompt_passed = false THEN 1 END) as rejected_count,
       COUNT(CASE WHEN es.check_query_passed = false THEN 1 END) as invalid_count,
       COUNT(eq.id) - COUNT(es.id) as unsubmitted_count
FROM exam_attempts
JOIN class_joinees ON exam_attempts.class_joinee_id = class_joinees.id
JOIN users ON class_joinees.user_id = users.id
JOIN exam_questions eq ON exam_attempts.exam_id = eq.exam_id
LEFT JOIN exam_submissions es ON eq.id = es.exam_question_id AND es.exam_attempt_id = exam_attempts.id
WHERE exam_attempts.exam_id = $1
GROUP BY exam_attempts.id, class_joinees.id, users.id
ORDER BY exam_attempts.created_at DESC;

-- name: ExamQuestionAdd :one
INSERT INTO exam_questions (exam_id, original_question_id, order_num, title, description, check_query, check_prompt)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: ExamQuestionMaxOrderNum :one
SELECT COALESCE(MAX(order_num), 0)::int as max_order_num
FROM exam_questions
WHERE exam_id = $1;

-- name: ExamQuestionDelete :exec
DELETE FROM exam_questions
WHERE id = $1;

-- name: ExamQuestionEdit :one
UPDATE exam_questions
SET title = $2, description = $3, check_query = $4, check_prompt = $5
WHERE id = $1
RETURNING *;

-- name: ExamQuestionList :many
SELECT id, order_num, title, description
FROM exam_questions
WHERE exam_id = $1
ORDER BY order_num ASC;

-- name: CollectionQuestionGetById :one
SELECT sqlc.embed(collection_questions), sqlc.embed(collections)
FROM collection_questions
JOIN collections ON collection_questions.collection_id = collections.id
WHERE collection_questions.id = $1;

-- name: ExamGetById :one
SELECT *
FROM exams
WHERE id = $1;

-- name: ExamDetail :one
SELECT sqlc.embed(exams), sqlc.embed(classes), sqlc.embed(collections)
FROM exams
JOIN classes ON exams.class_id = classes.id
JOIN collections ON exams.collection_id = collections.id
WHERE exams.id = $1;