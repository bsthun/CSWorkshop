# Claude

General guideline:
- Use pointer for struct.
- Use r as the receiver name. Example: `func (r *Handler) HandleOrganizationCreate(c *fiber.Ctx) error`.
- Comment in format of `// * only lowercase compact action` for each step.

Endpoint implementation guideline:

- Always use `c.Locals("l").(*jwt.Token).Claims.(*common.LoginClaims)` to get user claims which contains .UserId.
- Use `r.database.P()` as postgres querier and `r.database.C()` as clickhouse querier. Example: `user, err := r.database.P().GetUserById(c.Context(), u.UserId)`.
- Use gut.Iterate to iterate over array and map to another array. Example: `organizationItems, _ := gut.Iterate(organizations, func(organization sqlcpg.GetUserOrganizationsRow) (*payload.OrganizationItem, *gut.ErrorInstance)`.
- Use `response.Success(payload)` to return success response. For inline struct, always use `response.Success(c, &payload.Type{})` to avoid copy.
- Sqlc output is pointer by default, as well as payload. Use pointer as basis
- When using gut.Iterate with SQLC results, the function parameter should NOT be a pointer (e.g., `func(collection psql.Collection)` not `func(collection *psql.Collection)`) because SQLC returns `[]Collection` not `[]*Collection`
- Any function that return *gut.ErrorInstance should be named `er` amd handled with `if er != nil { return er }` without new gut.Err function.
- The endpoint file should have only handler functions, types defined in `payload` package, and no other logic.

Query guideline:

- Use `select *` or `sql.embed` to fetch all column by default. as well as `returning *`.
- Query name must begin with entity, example: UserList, OrganizationGetByUserId.
- Id variable in project will have lowercase d.
- For search queries with nullable parameters, use `($N::text IS NULL OR condition)` pattern
- For ordering with dynamic fields, use `CASE WHEN $N = 'field' THEN field_name END` pattern
- Always provide both count and list queries for paginated endpoints (e.g., `CollectionCount` and `CollectionList`)
- Do not uses table alias in query if not necessary.
- Use `sqlc.narg` instead of `Column1, Column2`. Use name of `sort` for order by.

SQLC Generate Type Mapping Notes:

- When using `sqlc.embed(table_name)` in queries, the generated row struct will use singular table names (e.g., `row.Collection` not `row.Collections`, `row.Exam` not `row.Exams`)
- For queries with embedded structs, access fields like: `row.Collection.Id`, `row.User.Name`, `row.Class.RegisterCode`
- For queries with joins and embeds, the generated struct follows pattern: `type QueryNameRow struct { EntityName Entity, AnotherEntity Entity, CustomField *Type }`
- Always check generated types in `/generate/psql/*.sql.go` after writing queries to understand the exact field structure
- When using `COUNT()` in queries with embed, the count field will be directly on the row struct (e.g., `row.QuestionCount`)

Parameter Type Guidelines:

- All SQLC query parameters should be pointers (e.g., `&classId` not `classId`) when the database field is nullable or when following the project's pointer convention
- Route parameters need to be converted from string: `strconv.ParseUint(c.Params("id"), 10, 64)`
- Payload validation should use `validate:"required"` for required fields
- Time fields in request payloads should be `*time.Time` and validated as required if database field is NOT NULL

Common Implementation Pattern for Detail/List Endpoints:

- Detail endpoints: Get single entity with related data using joins and embeds
- List endpoints: Get multiple entities with aggregated data (counts, etc.)
- Always use ORDER BY with logical sorting (usually created_at DESC)
- Map SQLC results to payload structs using gut.Iterate for arrays

Command:
- To check implementation, use `make generate` to check and generate code, if this command passed, it means the implementation is correct.

Common SQLC Generation Issues (Manual Fixes Required):

- **ExamQuestionMaxOrderNum Query**: The auto-generated function returns `interface{}` instead of typed result. Need to manually edit the SQL query to cast result as specific type:
  - Original: `SELECT COALESCE(MAX(order_num), 0) as max_order_num`
  - Fixed: `SELECT COALESCE(MAX(order_num), 0)::int as max_order_num`
  - This makes SQLC generate `*int32` return type instead of `interface{}`

- **Parameter Struct Issues**: When SQLC generates functions with multiple parameters, it sometimes creates generic `Column1, Column2` instead of named parameters. Use `sqlc.narg(param_name)` to fix:
  - Original: `WHERE ($1::type IS NULL OR condition)` 
  - Fixed: `WHERE (sqlc.narg(param_name)::type IS NULL OR condition)`

- **Always regenerate SQLC after SQL changes**: Use `make sqlc` to regenerate types before implementing Go code to see correct function signatures and parameter structures.