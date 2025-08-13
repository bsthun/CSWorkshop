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

func (r *Handler) HandleCollectionQuestionList(c *fiber.Ctx) error {
	// * get user claims
	u := c.Locals("l").(*jwt.Token).Claims.(*common.LoginClaims)

	// * parse body
	body := new(payload.CollectionIdRequest)
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

	// * get collection questions
	questions, err := r.database.P().CollectionQuestionList(c.Context(), body.CollectionId)
	if err != nil {
		return gut.Err(false, "failed to get collection questions", err)
	}

	// * map to payload
	questionItems, _ := gut.Iterate(questions, func(question psql.CollectionQuestionListRow) (*payload.QuestionListItem, *gut.ErrorInstance) {
		return &payload.QuestionListItem{
			Id:          question.Id,
			OrderNum:    question.OrderNum,
			Title:       question.Title,
			Description: question.Description,
		}, nil
	})

	// * response
	return c.JSON(response.Success(c, &payload.CollectionQuestionListResponse{
		Questions: questionItems,
	}))
}
