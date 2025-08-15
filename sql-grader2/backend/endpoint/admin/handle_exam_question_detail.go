package adminEndpoint

import (
	"backend/type/common"
	"backend/type/payload"
	"backend/type/response"

	"github.com/bsthun/gut"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func (r *Handler) HandleExamQuestionDetail(c *fiber.Ctx) error {
	// * get user claims
	u := c.Locals("l").(*jwt.Token).Claims.(*common.LoginClaims)

	// * parse body
	body := new(payload.ExamQuestionIdRequest)
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

	// * get exam question detail
	row, err := r.database.P().ExamQuestionDetail(c.Context(), body.ExamQuestionId)
	if err != nil {
		return gut.Err(false, "failed to get exam question detail", err)
	}

	// * prepare response
	responsePayload := &payload.ExamQuestionDetailResponse{
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
		CollectionQuestion: &payload.CollectionQuestionDetail{
			Id:           row.CollectionQuestion.Id,
			CollectionId: row.CollectionQuestion.CollectionId,
			OrderNum:     row.CollectionQuestion.OrderNum,
			Title:        row.CollectionQuestion.Title,
			Description:  row.CollectionQuestion.Description,
			CheckQuery:   row.CollectionQuestion.CheckQuery,
			CheckPrompt:  row.CollectionQuestion.CheckPrompt,
			CreatedAt:    row.CollectionQuestion.CreatedAt,
			UpdatedAt:    row.CollectionQuestion.UpdatedAt,
		},
	}

	// * response
	return c.JSON(response.Success(c, responsePayload))
}
