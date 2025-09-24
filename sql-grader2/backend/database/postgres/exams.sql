-- name: ExamCreate :one
INSERT INTO exams (class_id, collection_id, name, access_code, opened_at, closed_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: ExamList :many
SELECT sqlc.embed(exams),
       sqlc.embed(collections),
       (SELECT COUNT(*)::BIGINT FROM exam_attempts WHERE exam_attempts.exam_id = exams.id) AS exam_attempt_count,
       (SELECT COUNT(*)::BIGINT FROM exam_questions WHERE exam_questions.exam_id = exams.id) AS exam_question_count,
       (SELECT COUNT(*)::BIGINT FROM collection_questions WHERE collection_questions.collection_id = collections.id) AS collection_question_count
FROM exams
JOIN collections ON exams.collection_id = collections.id
WHERE exams.class_id = $1
GROUP BY exams.id, collections.id
ORDER BY exams.created_at DESC;

-- name: ExamJoineeList :many
SELECT sqlc.embed(exam_attempts), sqlc.embed(class_joinees), sqlc.embed(users),
       COUNT(CASE WHEN latest_submissions.check_query_passed = true AND latest_submissions.check_prompt_passed = true THEN 1 END) as passed_count,
       COUNT(CASE WHEN latest_submissions.check_query_passed = true AND latest_submissions.check_prompt_passed = false THEN 1 END) as rejected_count,
       COUNT(CASE WHEN latest_submissions.check_query_passed = false THEN 1 END) as invalid_count,
       COUNT(eq.id) - COUNT(latest_submissions.id) as unsubmitted_count
FROM exam_attempts
JOIN class_joinees ON exam_attempts.class_joinee_id = class_joinees.id
JOIN users ON class_joinees.user_id = users.id
JOIN exam_questions eq ON exam_attempts.exam_id = eq.exam_id
LEFT JOIN LATERAL (
    SELECT DISTINCT ON (exam_question_id) *
    FROM exam_submissions
    WHERE exam_question_id = eq.id AND exam_attempt_id = exam_attempts.id
    ORDER BY exam_question_id, created_at DESC
) latest_submissions ON true
WHERE exam_attempts.exam_id = $1
GROUP BY exam_attempts.id, class_joinees.id, users.id
ORDER BY exam_attempts.created_at DESC;

-- name: ExamQuestionAdd :one
INSERT INTO exam_questions (exam_id, original_question_id, order_num, title, description, check_query, check_prompt)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: ExamQuestionMaxOrderNum :one
SELECT COALESCE(MAX(order_num), 0)::INT as max_order_num
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
SELECT sqlc.embed(exams),
       sqlc.embed(classes),
       sqlc.embed(collections),
       COALESCE(COUNT(CASE WHEN exam_attempts.started_at IS NOT NULL THEN 1 END), 0)::BIGINT as attempt_started_count,
       COALESCE(COUNT(CASE WHEN exam_attempts.finished_at IS NOT NULL THEN 1 END), 0)::BIGINT as attempt_finished_count,
       COALESCE(COUNT(DISTINCT exam_questions.id), 0)::BIGINT as exam_question_count,
       COALESCE(COUNT(DISTINCT collection_questions.id), 0)::BIGINT as collection_question_count
FROM exams
JOIN classes ON exams.class_id = classes.id
JOIN collections ON exams.collection_id = collections.id
LEFT JOIN exam_attempts ON exams.id = exam_attempts.exam_id
LEFT JOIN exam_questions ON exams.id = exam_questions.exam_id
LEFT JOIN collection_questions ON collections.id = collection_questions.collection_id
WHERE exams.id = $1
GROUP BY exams.id, classes.id, collections.id;

-- name: ExamScoreDistribution :many
WITH student_scores AS (
    SELECT 
        ea.id as attempt_id,
        COALESCE(COUNT(CASE WHEN es.check_query_passed = true AND es.check_prompt_passed = true THEN 1 END), 0)::BIGINT as passed_count
    FROM exam_attempts ea
    LEFT JOIN exam_submissions es ON ea.id = es.exam_attempt_id
    WHERE ea.exam_id = $1 AND ea.finished_at IS NOT NULL
    GROUP BY ea.id
),
question_count AS (
    SELECT COALESCE(COUNT(*), 0)::BIGINT as total_questions
    FROM exam_questions
    WHERE exam_id = $1
),
score_distribution AS (
    SELECT 
        passed_count as score,
        COUNT(*)::BIGINT as student_count
    FROM student_scores
    GROUP BY passed_count
    UNION ALL
    SELECT 
        generate_series(0, (SELECT total_questions FROM question_count))::BIGINT as score,
        0::BIGINT as student_count
)
SELECT 
    score,
    SUM(student_count)::BIGINT as student_count
FROM score_distribution
GROUP BY score
ORDER BY score;

-- name: ExamQuestionDetail :one
SELECT sqlc.embed(exam_questions), sqlc.embed(collection_questions)
FROM exam_questions
JOIN collection_questions ON exam_questions.original_question_id = collection_questions.id
WHERE exam_questions.id = $1;

-- name: ClassExamList :many
SELECT sqlc.embed(exams),
       sqlc.embed(exam_attempts),
       COALESCE(COUNT(DISTINCT exam_questions.id), 0)::BIGINT as exam_question_count
FROM exams
LEFT JOIN exam_questions ON exams.id = exam_questions.exam_id
LEFT JOIN class_joinees ON class_joinees.class_id = exams.class_id AND class_joinees.user_id = sqlc.narg(user_id)
LEFT JOIN exam_attempts ON exams.id = exam_attempts.exam_id AND exam_attempts.class_joinee_id = class_joinees.id
WHERE exams.class_id = $1
GROUP BY exams.id, exam_attempts.id
ORDER BY exams.created_at DESC;

-- name: ExamAttemptCreate :one
INSERT INTO exam_attempts (exam_id, class_joinee_id, database_name)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ExamAttemptGetByExamAndJoinee :one
SELECT *
FROM exam_attempts
WHERE exam_id = $1 AND class_joinee_id = $2;

-- name: ExamAttemptUpdateStartedAt :one
UPDATE exam_attempts
SET started_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: ExamAttemptUpdateUpdatedAt :one
UPDATE exam_attempts
SET updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: ExamAttemptGetById :one
SELECT sqlc.embed(exam_attempts),
       sqlc.embed(exams),
       sqlc.embed(class_joinees),
       sqlc.embed(users),
       sqlc.embed(classes),
       sqlc.embed(semesters)
FROM exam_attempts
JOIN exams ON exam_attempts.exam_id = exams.id
JOIN class_joinees ON exam_attempts.class_joinee_id = class_joinees.id
JOIN users ON class_joinees.user_id = users.id
JOIN classes ON exams.class_id = classes.id
JOIN semesters ON classes.semester_id = semesters.id
WHERE exam_attempts.id = $1;

-- name: ExamAttemptGetByDatabaseName :one
SELECT *
FROM exam_attempts
WHERE database_name = $1;

-- name: ExamQuestionListByAttempt :many
SELECT sqlc.embed(exam_questions),
       CASE 
           WHEN es.check_query_passed = true AND es.check_prompt_passed = true THEN 'passed'
           WHEN es.check_query_passed = true AND es.check_prompt_passed = false THEN 'rejected'
           WHEN es.check_query_passed = false THEN 'invalid'
           ELSE 'unsubmitted'
       END as status
FROM exam_questions
LEFT JOIN LATERAL (
    SELECT * FROM exam_submissions 
    WHERE exam_question_id = exam_questions.id 
    AND exam_attempt_id = $1
    ORDER BY created_at DESC
    LIMIT 1
) es ON true
WHERE exam_questions.exam_id = (SELECT exam_id FROM exam_attempts WHERE id = $1)
ORDER BY exam_questions.order_num ASC;

-- name: ExamEdit :one
UPDATE exams
SET name = COALESCE(sqlc.narg(name), name),
    opened_at = COALESCE(sqlc.narg(opened_at), opened_at),
    closed_at = COALESCE(sqlc.narg(closed_at), closed_at)
WHERE id = sqlc.arg(exam_id)
RETURNING *;
