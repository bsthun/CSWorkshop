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

func (r *Handler) HandleClassExamList(c *fiber.Ctx) error {
	// * get user claims
	u := c.Locals("l").(*jwt.Token).Claims.(*common.LoginClaims)

	// * parse body
	body := new(payload.ClassExamListRequest)
	if err := c.BodyParser(body); err != nil {
		return gut.Err(false, "invalid body", err)
	}

	// * validate body
	if err := gut.Validate(body); err != nil {
		return err
	}

	// * check if user is member of the class
	_, err := r.database.P().ClassJoineeGetByUserAndClass(c.Context(), &psql.ClassJoineeGetByUserAndClassParams{
		UserId:  u.UserId,
		ClassId: body.ClassId,
	})
	if err != nil {
		return gut.Err(false, "access denied or class not found", err)
	}

	// * get exams for the class
	examRows, err := r.database.P().ClassExamList(c.Context(), body.ClassId)
	if err != nil {
		return gut.Err(false, "failed to get exams", err)
	}

	// * map exams to items
	examItems, _ := gut.Iterate(examRows, func(examRow psql.ClassExamListRow) (*payload.ClassExamListItem, *gut.ErrorInstance) {
		return &payload.ClassExamListItem{
			Exam: &payload.Exam{
				Id:           examRow.Exam.Id,
				ClassId:      examRow.Exam.ClassId,
				CollectionId: examRow.Exam.CollectionId,
				Name:         examRow.Exam.Name,
				AccessCode:   examRow.Exam.AccessCode,
				OpenedAt:     examRow.Exam.OpenedAt,
				ClosedAt:     examRow.Exam.ClosedAt,
				CreatedAt:    examRow.Exam.CreatedAt,
				UpdatedAt:    examRow.Exam.UpdatedAt,
			},
			QuestionCount: examRow.ExamQuestionCount,
		}, nil
	})

	// * response
	return c.JSON(response.Success(c, &payload.ClassExamListResponse{
		Exams: examItems,
	}))
}
