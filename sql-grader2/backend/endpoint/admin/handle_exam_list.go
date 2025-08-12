package adminEndpoint

import (
	"backend/generate/psql"
	"backend/type/common"
	"backend/type/payload"
	"backend/type/response"
	"strconv"

	"github.com/bsthun/gut"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func (r *Handler) HandleExamList(c *fiber.Ctx) error {
	// * get user claims
	u := c.Locals("l").(*jwt.Token).Claims.(*common.LoginClaims)

	// * get class id from params
	classIdParam := c.Params("id")
	classId, err := strconv.ParseUint(classIdParam, 10, 64)
	if err != nil {
		return gut.Err(false, "invalid class id", err)
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

	// * get exam list with collection info and question count
	examRows, err := r.database.P().ExamList(c.Context(), &classId)
	if err != nil {
		return gut.Err(false, "failed to get exams", err)
	}

	// * map exams to payload
	exams, er := gut.Iterate(examRows, func(row psql.ExamListRow) (*payload.ExamListItem, *gut.ErrorInstance) {
		return &payload.ExamListItem{
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
			Collection: &payload.Collection{
				Id:        row.Collection.Id,
				Name:      row.Collection.Name,
				Metadata:  row.Collection.Metadata,
				CreatedAt: row.Collection.CreatedAt,
				UpdatedAt: row.Collection.UpdatedAt,
			},
			QuestionCount: row.QuestionCount,
		}, nil
	})
	if er != nil {
		return er
	}

	// * prepare response
	responsePayload := &payload.ExamListResponse{
		Exams: exams,
	}

	// * response
	return c.JSON(response.Success(c, responsePayload))
}
