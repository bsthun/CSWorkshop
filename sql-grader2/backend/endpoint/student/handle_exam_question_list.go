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

func (r *Handler) HandleStudentExamQuestionList(c *fiber.Ctx) error {
	// * get user claims
	u := c.Locals("l").(*jwt.Token).Claims.(*common.LoginClaims)

	// * parse body
	body := new(payload.ExamAttemptIdRequest)
	if err := c.BodyParser(body); err != nil {
		return gut.Err(false, "invalid body", err)
	}

	// * validate body
	if err := gut.Validate(body); err != nil {
		return err
	}

	// * check if exam is opened
	if er := r.submissionProcedure.ServeOpened(c.Context(), u.UserId, body.ExamAttemptId); er != nil {
		return er
	}

	// * get exam attempt to verify ownership
	attemptData, err := r.database.P().ExamAttemptGetById(c.Context(), body.ExamAttemptId)
	if err != nil {
		return gut.Err(false, "exam attempt not found", err)
	}

	// * verify user owns this attempt
	if *attemptData.User.Id != *u.UserId {
		return gut.Err(false, "access denied", nil)
	}

	// * get exam questions with status
	questionRows, err := r.database.P().ExamQuestionListByAttempt(c.Context(), body.ExamAttemptId)
	if err != nil {
		return gut.Err(false, "failed to get exam questions", err)
	}

	// * map questions to items
	items, _ := gut.Iterate(questionRows, func(row psql.ExamQuestionListByAttemptRow) (*payload.ExamQuestionWithStatus, *gut.ErrorInstance) {
		return &payload.ExamQuestionWithStatus{
			ExamQuestion: &payload.ExamQuestion{
				Id:                 row.ExamQuestion.Id,
				ExamId:             row.ExamQuestion.ExamId,
				OriginalQuestionId: row.ExamQuestion.OriginalQuestionId,
				OrderNum:           row.ExamQuestion.OrderNum,
				Title:              row.ExamQuestion.Title,
				Description:        row.ExamQuestion.Description,
				CheckQuery:         row.ExamQuestion.CheckQuery,
				CheckPrompt:        row.ExamQuestion.CheckPrompt,
				CreatedAt:          row.ExamQuestion.CreatedAt,
				UpdatedAt:          row.ExamQuestion.UpdatedAt,
			},
			Status: row.Status,
		}, nil
	})

	// * return success response
	return c.JSON(response.Success(c, &payload.ExamQuestionListByAttemptResponse{
		Questions: items,
	}))
}
