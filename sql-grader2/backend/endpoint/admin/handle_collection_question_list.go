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

func (r *Handler) HandleCollectionQuestionList(c *fiber.Ctx) error {
	// * get user claims
	u := c.Locals("l").(*jwt.Token).Claims.(*common.LoginClaims)

	// * get collection id from query params
	collectionIdStr := c.Query("collectionId")
	if collectionIdStr == "" {
		return gut.Err(false, "collection id is required", nil)
	}

	collectionId, err := strconv.ParseUint(collectionIdStr, 10, 64)
	if err != nil {
		return gut.Err(false, "invalid collection id", err)
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

	// * get collection questions
	questions, err := r.database.P().CollectionQuestionList(c.Context(), gut.Ptr(collectionId))
	if err != nil {
		return gut.Err(false, "failed to get collection questions", err)
	}

	// * map to payload
	questionItems, _ := gut.Iterate(questions, func(question psql.CollectionQuestionListRow) (*payload.CollectionQuestionListItem, *gut.ErrorInstance) {
		return &payload.CollectionQuestionListItem{
			Id:       question.Id,
			OrderNum: question.OrderNum,
			Title:    question.Title,
		}, nil
	})

	// * response
	return c.JSON(response.Success(c, &payload.CollectionQuestionListResponse{
		Questions: questionItems,
	}))
}
