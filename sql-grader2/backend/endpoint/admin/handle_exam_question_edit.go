package adminEndpoint

import (
	"backend/generate/psql"
	"backend/type/common"
	"backend/type/payload"
	"backend/type/response"

	"github.com/bsthun/gut"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func (r *Handler) HandleExamQuestionEdit(c *fiber.Ctx) error {
	// * get user claims
	u := c.Locals("l").(*jwt.Token).Claims.(*common.LoginClaims)

	// * parse body
	body := new(payload.ExamQuestionEditRequest)
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

	// * update exam question
	examQuestion, err := r.database.P().ExamQuestionEdit(c.Context(), &psql.ExamQuestionEditParams{
		Id:          body.ExamQuestionId,
		Title:       body.Title,
		Description: body.Description,
		CheckQuery:  body.CheckQuery,
		CheckPrompt: body.CheckPrompt,
	})
	if err != nil {
		return gut.Err(false, "failed to update exam question", err)
	}

	// * prepare response
	responsePayload := &payload.ExamQuestion{
		Id:                 examQuestion.Id,
		ExamId:             examQuestion.ExamId,
		OriginalQuestionId: examQuestion.OriginalQuestionId,
		OrderNum:           examQuestion.OrderNum,
		Title:              examQuestion.Title,
		Description:        examQuestion.Description,
		CheckQuery:         examQuestion.CheckQuery,
		CheckPrompt:        examQuestion.CheckPrompt,
		CreatedAt:          examQuestion.CreatedAt,
		UpdatedAt:          examQuestion.UpdatedAt,
	}

	// * response
	return c.JSON(response.Success(c, responsePayload))
}
