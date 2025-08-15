package adminEndpoint

import (
	"backend/type/common"
	"backend/type/payload"
	"backend/type/response"

	"github.com/bsthun/gut"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func (r *Handler) HandleExamDetail(c *fiber.Ctx) error {
	// * get user claims
	u := c.Locals("l").(*jwt.Token).Claims.(*common.LoginClaims)

	// * parse body
	body := new(payload.ExamIdRequest)
	if err := c.BodyParser(body); err != nil {
		return gut.Err(false, "invalid body", err)
	}

	// * get user from database
	user, err := r.database.P().UserGetById(c.Context(), u.UserId)
	if err != nil {
		return gut.Err(false, "failed to get user", err)
	}

	// * check if user is admin
	if user.IsAdmin == nil || !*user.IsAdmin {
		return gut.Err(false, "access denied", nil)
	}

	// * get exam detail with class and collection
	row, err := r.database.P().ExamDetail(c.Context(), body.ExamId)
	if err != nil {
		return gut.Err(false, "failed to get exam detail", err)
	}

	// * prepare response
	responsePayload := &payload.ExamDetailResponse{
		Exam: &payload.Exam{
			Id:           row.Exam.Id,
			ClassId:      row.Exam.ClassId,
			CollectionId: row.Exam.CollectionId,
			Name:         row.Exam.Name,
			AccessCode:   row.Exam.AccessCode,
			OpenedAt:     row.Exam.OpenedAt,
			ClosedAt:     row.Exam.ClosedAt,
			CreatedAt:    row.Exam.CreatedAt,
			UpdatedAt:    row.Exam.UpdatedAt,
		},
		Class: &payload.Class{
			Id:           row.Class.Id,
			SemesterId:   row.Class.SemesterId,
			Code:         row.Class.Code,
			Name:         row.Class.Name,
			RegisterCode: row.Class.RegisterCode,
			CreatedAt:    row.Class.CreatedAt,
			UpdatedAt:    row.Class.UpdatedAt,
		},
		Collection: &payload.Collection{
			Id:        row.Collection.Id,
			Name:      row.Collection.Name,
			Metadata:  row.Collection.Metadata,
			CreatedAt: row.Collection.CreatedAt,
			UpdatedAt: row.Collection.UpdatedAt,
		},
		AttemptCount: &payload.ExamAttemptCount{
			OpenedCount:   row.AttemptOpenedCount,
			StartedCount:  row.AttemptStartedCount,
			FinishedCount: row.AttemptFinishedCount,
		},
	}

	// * response
	return c.JSON(response.Success(c, responsePayload))
}
