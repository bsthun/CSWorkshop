-- name: CollectionList :many
SELECT collections.*,
       COALESCE(COUNT(collection_questions.collection_id), 0)::bigint as question_count
FROM collections
LEFT JOIN collection_questions ON collections.id = collection_questions.collection_id
WHERE (sqlc.narg('name')::text IS NULL OR LOWER(collections.name) LIKE LOWER('%' || sqlc.narg('name') || '%'))
GROUP BY collections.id
ORDER BY
    CASE WHEN sqlc.narg('sort') = 'name' THEN collections.name END ASC,
    CASE WHEN sqlc.narg('sort') = 'createdAt' THEN collections.created_at END DESC
LIMIT $1 OFFSET $2;

-- name: CollectionCount :one
SELECT COUNT(*)
FROM collections
WHERE ($1::text IS NULL OR LOWER(name) LIKE LOWER('%' || $1 || '%'));

-- name: CollectionCreate :one
INSERT INTO collections (name, metadata)
VALUES ($1, '{}'::jsonb)
RETURNING *;

-- name: CollectionDetail :one
SELECT collections.*,
       COALESCE(COUNT(collection_questions.collection_id), 0)::bigint as question_count
FROM collections
         LEFT JOIN collection_questions ON collections.id = collection_questions.collection_id
WHERE collections.id = $1
GROUP BY collections.id;

-- name: CollectionEdit :one
UPDATE collections
SET name = $2
WHERE id = $1
RETURNING *;

-- name: CollectionUpdateMetadata :one
UPDATE collections
SET metadata = $2
WHERE id = $1
RETURNING *;