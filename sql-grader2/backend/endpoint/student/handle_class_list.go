package studentEndpoint

import (
	"backend/generate/psql"
	"backend/type/common"
	"backend/type/payload"
	"backend/type/response"

	"github.com/bsthun/gut"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func (r *Handler) HandleClassList(c *fiber.Ctx) error {
	// * get user claims
	u := c.Locals("l").(*jwt.Token).Claims.(*common.LoginClaims)

	// * parse body
	body := new(payload.StudentClassListRequest)
	if err := c.BodyParser(body); err != nil {
		return gut.Err(false, "invalid body", err)
	}

	// * validate body
	if err := gut.Validate(body); err != nil {
		return err
	}

	// * get student class list
	rows, err := r.database.P().StudentClassList(c.Context(), u.UserId)
	if err != nil {
		return gut.Err(false, "failed to get classes", err)
	}

	// * map classes to items
	classes, _ := gut.Iterate(rows, func(row psql.StudentClassListRow) (*payload.StudentClassListItem, *gut.ErrorInstance) {
		return &payload.StudentClassListItem{
			Class: &payload.Class{
				Id:           row.Class.Id,
				SemesterId:   row.Class.SemesterId,
				Code:         row.Class.Code,
				Name:         row.Class.Name,
				RegisterCode: row.Class.RegisterCode,
				CreatedAt:    row.Class.CreatedAt,
				UpdatedAt:    row.Class.UpdatedAt,
			},
			Semester: &payload.SemesterInfo{
				Id:        row.Semester.Id,
				Name:      row.Semester.Name,
				CreatedAt: row.Semester.CreatedAt,
				UpdatedAt: row.Semester.UpdatedAt,
			},
			ExamTotalCount:    row.ExamTotalCount,
			ExamFinishedCount: row.ExamFinishedCount,
		}, nil
	})

	// * response
	return c.JSON(response.Success(c, &payload.StudentClassListResponse{
		Classes: classes,
	}))
}
