# Claude Code Implementation Guide

## General Rules

- Use pointer for struct
- Receiver name: `r` (e.g., `func (r *Handler) HandleCreate(c *fiber.Ctx) error`)
- Comment format: `// * lowercase compact action`

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

### SQLC Embeds

```sql
-- name: CollectionDetail :one
SELECT sqlc.embed(collections), sqlc.embed(users), COUNT(items) as item_count
FROM collections
         JOIN users ON collections.user_id = users.id
WHERE collections.id = $1;
```

- Generated: `row.Collection.Id`, `row.User.Name`, `row.ItemCount`

## Common Fixes

- ExamQuestionMaxOrderNum: `COALESCE(MAX(order_num), 0)::bigint`, need to cast to bigint
- Use `sqlc.narg(param_name)` for named parameters
- Queries in `./database/postgres/*.sql`

## Validation

Always run `make generate` to check implementation