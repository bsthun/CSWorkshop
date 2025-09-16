-- name: SubmissionDetail :one
SELECT sqlc.embed(exam_submissions), 
       sqlc.embed(exam_questions), 
       sqlc.embed(exam_attempts)
FROM exam_submissions
JOIN exam_questions ON exam_submissions.exam_question_id = exam_questions.id
JOIN exam_attempts ON exam_submissions.exam_attempt_id = exam_attempts.id
WHERE exam_submissions.id = $1;

-- name: SubmissionList :many
SELECT sqlc.embed(exam_submissions), 
       sqlc.embed(exam_questions), 
       sqlc.embed(exam_attempts),
       sqlc.embed(class_joinees),
       sqlc.embed(users)
FROM exam_submissions
JOIN exam_questions ON exam_submissions.exam_question_id = exam_questions.id
JOIN exam_attempts ON exam_submissions.exam_attempt_id = exam_attempts.id
JOIN class_joinees ON exam_attempts.class_joinee_id = class_joinees.id
JOIN users ON class_joinees.user_id = users.id
WHERE (sqlc.narg(exam_attempt_id)::BIGINT IS NULL OR exam_submissions.exam_attempt_id = sqlc.narg(exam_attempt_id))
  AND (sqlc.narg(exam_question_id)::BIGINT IS NULL OR exam_submissions.exam_question_id = sqlc.narg(exam_question_id))
ORDER BY exam_submissions.created_at DESC;

-- name: ExamSubmissionCreate :one
INSERT INTO exam_submissions (exam_question_id, exam_attempt_id, answer, result, check_query_passed, check_query_at, check_prompt_passed, check_prompt_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: ExamQuestionGetById :one
SELECT *
FROM exam_questions
WHERE id = $1;

-- name: ExamQuestionDetailWithSubmissions :many
SELECT sqlc.embed(exam_questions),
       sqlc.embed(exam_submissions),
       CASE 
           WHEN exam_submissions.check_query_passed = true AND exam_submissions.check_prompt_passed = true THEN 'passed'
           WHEN exam_submissions.check_query_passed = true AND exam_submissions.check_prompt_passed = false THEN 'rejected'
           WHEN exam_submissions.check_query_passed = false THEN 'invalid'
           ELSE 'unsubmitted'
       END as status
FROM exam_questions
LEFT JOIN exam_submissions ON exam_submissions.exam_question_id = exam_questions.id 
                            AND exam_submissions.exam_attempt_id = $2
WHERE exam_questions.id = $1
ORDER BY exam_submissions.created_at DESC;