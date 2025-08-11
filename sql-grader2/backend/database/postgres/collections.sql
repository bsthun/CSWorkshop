-- name: CollectionList :many
SELECT *
FROM collections
WHERE (sqlc.narg('name')::text IS NULL OR LOWER(name) LIKE LOWER('%' || sqlc.narg('name') || '%'))
ORDER BY
    CASE WHEN sqlc.narg('sort') = 'name' THEN name END ASC,
    CASE WHEN sqlc.narg('sort') = 'createdAt' THEN created_at END DESC
LIMIT $1 OFFSET $2;

-- name: CollectionCount :one
SELECT COUNT(*)
FROM collections
WHERE ($1::text IS NULL OR LOWER(name) LIKE LOWER('%' || $1 || '%'));

-- name: CollectionCreate :one
INSERT INTO collections (name, metadata)
VALUES ($1, '{}'::jsonb)
RETURNING *;