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

func (r *Handler) HandleExamQuestionList(c *fiber.Ctx) error {
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

	// * get exam question list
	rows, err := r.database.P().ExamQuestionList(c.Context(), body.ExamId)
	if err != nil {
		return gut.Err(false, "failed to get exam questions", err)
	}

	// * map questions to payload
	questions, er := gut.Iterate(rows, func(row psql.ExamQuestionListRow) (*payload.CollectionQuestionItem, *gut.ErrorInstance) {
		return &payload.CollectionQuestionItem{
			Id:          row.Id,
			OrderNum:    row.OrderNum,
			Title:       row.Title,
			Description: row.Description,
		}, nil
	})
	if er != nil {
		return er
	}

	// * prepare response
	responsePayload := &payload.ExamQuestionListResponse{
		Questions: questions,
	}

	// * response
	return c.JSON(response.Success(c, responsePayload))
}
