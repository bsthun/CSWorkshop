# Claude

General guideline:
- Use pointer for struct.
- Use r as the receiver name. Example: `func (r *Handler) HandleOrganizationCreate(c *fiber.Ctx) error`.
- Comment in format of `// * only lowercase compact action` for each step.

Endpoint implementation guideline:

- Always use `c.Locals("l").(*jwt.Token).Claims.(*common.UserClaims)` to get user claims which contains .UserId.
- Use `r.database.P()` as postgres querier and `r.database.C()` as clickhouse querier. Example: `user, err := r.database.P().GetUserById(c.Context(), u.UserId)`.
- Use gut.Iterate to iterate over array and map to another array. Example: `organizationItems, _ := gut.Iterate(organizations, func(organization sqlcpg.GetUserOrganizationsRow) (*payload.OrganizationItem, *gut.ErrorInstance)`.
- Use `response.Success(payload)` to return success response. For inline struct, always use `response.Success(c, &payload.Type{})` to avoid copy.
- Sqlc output is pointer by default, as well as payload. Use pointer as basis
- Any function that return *gut.ErrorInstance should be named `er` amd handled with `if er != nil { return er }` without new gut.Err function.

Query guideline:

- Use `select *` or `sql.embed` to fetch all column by default. as well as `returning *`.
- Query name must begin with entity, example: UserList, OrganizationGetByUserId.
- Id variable in project will have lowercase d.

Command:
- To check implementation, use `make generate` to check and generate code, if this command passed, it means the implementation is correct.