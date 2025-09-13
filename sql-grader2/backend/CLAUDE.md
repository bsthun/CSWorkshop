# Claude Code Implementation Guide

## General Rules

- Use pointer for struct
- Receiver name: `r` (e.g., `func (r *Handler) HandleCreate(c *fiber.Ctx) error`)
- Comment format: `// * lowercase compact action`
- Use camelCase for json tags

## Endpoint Implementation

### Request

- Always use this snippet, change only payload:
    ```go
    // * parse body
    body := new(payload.HostCreateRequest)
    if err := c.BodyParser(body); err != nil {
    return gut.Err(false, "invalid body", err)
    }
    
    // * validate body
    if err := gut.Validate(body); err != nil {
    return err
    }
    ```

- Request always POST with payload
- Need to register endpoint `endpoint/endpoint.go`

### Claims & Database

```go
// * get user claims
u := c.Locals("l").(*jwt.Token).Claims.(*common.LoginClaims)

// * query postgres / clickhouse example
user, err := r.database.P().GetUserById(c.Context(), claims.UserId)
createdUser, err := r.database.P().CreateUser(c.Context(), psql.CreateUserParams{
Name:  req.Name,
Email: req.Email,
})
data, err := r.database.C().QueryClickhouse(c.Context(), params)
```

### Transactions

```go
// * begin transaction
tx, querier := r.database.Ptx(c.Context(), nil)
defer func() {
    if r := recover(); r != nil {
        _ = tx.Rollback()
    }
}()

// * query action create
action, err := querier.ActionCreate(c.Context(), &psql.ActionCreateParams{
    UserId:    u.UserId,
    ProjectId: body.ProjectId,
})
if err != nil {
    _ = tx.Rollback()
    return gut.Err(false, "failed to create action", err)
}

// * commit transaction
if err := tx.Commit(); err != nil {
    return gut.Err(false, "failed to commit transaction", err)
}
```

### Array Mapping & Response

```go
// * map organizations to items
items, _ := gut.Iterate(orgs, func(org psql.GetUserOrganizationsRow) (*payload.OrganizationItem, *gut.ErrorInstance) {
return &payload.OrganizationItem{Id: org.Id}, nil
})

// * return success response
return c.JSON(response.Success(c, &payload.OrganizationList{
	Items: items
}))
```

### Error Handling

```go
// * validate request
if er := validateRequest(req); er != nil { return er }
```

### Key Points

- Always read database schema from `generate/schema.sql`
- SQLC querier is in `./database/postgres/*.sql`, do not read migration files
- SQLC output is pointer by default
- gut.Iterate function param is NOT pointer: `func(collection psql.Collection)`
- ErrorInstance functions named `er`, handled with `if er != nil { return er }`
- Endpoint files contain only handlers, use `payload` package for types

## Query Guidelines

### Basic Structure

```sql
-- name: UserList :many
SELECT *
FROM users
WHERE ($1::text IS NULL OR name ILIKE '%' || $1 || '%')
ORDER BY CASE WHEN $2 = 'name' THEN name END;

-- name: UserCreate :one
INSERT INTO users (name, email)
VALUES ($1, $2)
RETURNING *;
```

### Naming & Patterns

- Query names: `EntityAction` (UserList, OrganizationGetByUserId)
- Use `select *` and `returning *`
- Id variables: lowercase `d` (userId not userID)
- Search: `($N::text IS NULL OR condition)`
- Dynamic ordering: `CASE WHEN $N = 'field' THEN field_name END`
- Provide both count/list for pagination
- Use `sqlc.narg(name)` instead of Column1, Column2
- Avoid use alias table when joining

### SQLC example

```sql
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
```

Key points:

- Generated: `row.Exam.Id`, `row.Collection.Name`, `row.ExamAttemptCount`
- Use subquery for counts

## Common Fixes

- ExamQuestionMaxOrderNum: `COALESCE(MAX(order_num), 0)::BIGINT`, need to cast to bigint
- Use `sqlc.narg(param_name)` for named parameters
- Queries in `./database/postgres/*.sql`

## Validation

Always run `make generate` to check implementation