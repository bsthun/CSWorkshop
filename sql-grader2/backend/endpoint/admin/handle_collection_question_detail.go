package adminEndpoint

import (
	"backend/type/common"
	"backend/type/payload"
	"backend/type/response"
	"strconv"

	"github.com/bsthun/gut"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func (r *Handler) HandleCollectionQuestionDetail(c *fiber.Ctx) error {
	// * get user claims
	u := c.Locals("l").(*jwt.Token).Claims.(*common.LoginClaims)

	// * get question id from query params
	questionIdStr := c.Query("id")
	if questionIdStr == "" {
		return gut.Err(false, "question id is required", nil)
	}

	questionId, err := strconv.ParseUint(questionIdStr, 10, 64)
	if err != nil {
		return gut.Err(false, "invalid question id", err)
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

	// * get collection question detail
	question, err := r.database.P().CollectionQuestionDetail(c.Context(), gut.Ptr(questionId))
	if err != nil {
		return gut.Err(false, "failed to get collection question", err)
	}

	// * map to payload
	questionDetail := &payload.CollectionQuestionDetail{
		Id:           question.Id,
		CollectionId: question.CollectionId,
		OrderNum:     question.OrderNum,
		Title:        question.Title,
		Description:  question.Description,
		CheckQuery:   question.CheckQuery,
		CheckPrompt:  question.CheckPrompt,
		CreatedAt:    question.CreatedAt,
		UpdatedAt:    question.UpdatedAt,
	}

	// * response
	return c.JSON(response.Success(c, questionDetail))
}
