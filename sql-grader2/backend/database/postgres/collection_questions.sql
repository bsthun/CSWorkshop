-- name: CollectionQuestionMaxOrderNum :one
SELECT COALESCE(MAX(order_num), 0)::int AS max_order_num
FROM collection_questions
WHERE collection_id = $1;

-- name: CollectionQuestionCreate :one
INSERT INTO collection_questions (collection_id, order_num, title, description, check_query, check_prompt)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: CollectionQuestionList :many
SELECT id, order_num, title, description
FROM collection_questions
WHERE collection_id = $1
ORDER BY order_num ASC;

-- name: CollectionQuestionDetail :one
SELECT *
FROM collection_questions
WHERE id = $1;

-- name: CollectionQuestionEdit :one
UPDATE collection_questions
SET title = $2, description = $3, check_query = $4, check_prompt = $5
WHERE id = $1
RETURNING *;

-- name: CollectionQuestionDelete :exec
DELETE FROM collection_questions
WHERE id = $1;

