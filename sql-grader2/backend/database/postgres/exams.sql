-- name: ExamCreate :one
INSERT INTO exams (class_id, collection_id, name, access_code, opened_at, closed_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: ExamList :many
SELECT sqlc.embed(exams), sqlc.embed(collections), COUNT(collection_questions.id) as question_count
FROM exams
JOIN collections ON exams.collection_id = collections.id
LEFT JOIN collection_questions ON collections.id = collection_questions.collection_id
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