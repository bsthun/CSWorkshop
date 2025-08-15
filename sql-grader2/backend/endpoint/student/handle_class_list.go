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

	// * get student class list with exam attempt counts
	rows, err := r.database.P().StudentClassList(c.Context(), u.UserId)
	if err != nil {
		return gut.Err(false, "failed to get classes", err)
	}

	// * map classes to payload
	classes, er := gut.Iterate(rows, func(row psql.StudentClassListRow) (*payload.StudentClassListItem, *gut.ErrorInstance) {
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
			ExamAttemptTotalCount:    row.ExamAttemptTotalCount,
			ExamAttemptFinishedCount: row.ExamAttemptFinishedCount,
		}, nil
	})
	if er != nil {
		return er
	}

	// * prepare response
	responsePayload := &payload.StudentClassListResponse{
		Classes: classes,
	}

	// * response
	return c.JSON(response.Success(c, responsePayload))
}
