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
WHERE (sqlc.narg(exam_attempt_id)::bigint IS NULL OR exam_submissions.exam_attempt_id = sqlc.narg(exam_attempt_id))
  AND (sqlc.narg(exam_question_id)::bigint IS NULL OR exam_submissions.exam_question_id = sqlc.narg(exam_question_id))
ORDER BY exam_submissions.created_at DESC;