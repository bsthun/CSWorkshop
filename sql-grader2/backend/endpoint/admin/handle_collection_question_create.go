package adminEndpoint

import (
	"backend/generate/psql"
	"backend/type/common"
	"backend/type/payload"
	"backend/type/response"
	"fmt"

	"github.com/bsthun/gut"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/tjarratt/babble"
)

func (r *Handler) HandleCollectionQuestionCreate(c *fiber.Ctx) error {
	// * get user claims
	u := c.Locals("l").(*jwt.Token).Claims.(*common.LoginClaims)

	// * parse body
	body := new(payload.CollectionQuestionCreateRequest)
	if err := c.BodyParser(body); err != nil {
		return gut.Err(false, "invalid body", err)
	}

	// * validate body
	if err := gut.Validate(body); err != nil {
		return err
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

	// * get max order num
	maxOrderResult, err := r.database.P().CollectionQuestionMaxOrderNum(c.Context(), gut.Ptr(uint64(*body.CollectionId)))
	if err != nil {
		return gut.Err(false, "failed to get max order num", err)
	}

	// * generate title using babble
	babbler := babble.NewBabbler()
	babbler.Count = 3
	babbler.Separator = " "
	nextOrderNum := *maxOrderResult + 1
	title := fmt.Sprintf("#%d %s", nextOrderNum, babbler.Babble())

	// * create collection question
	_, err = r.database.P().CollectionQuestionCreate(c.Context(), &psql.CollectionQuestionCreateParams{
		CollectionId: gut.Ptr(*body.CollectionId),
		OrderNum:     gut.Ptr(nextOrderNum),
		Title:        gut.Ptr(title),
		Description:  body.Description,
		CheckQuery:   gut.Ptr(""),
		CheckPrompt:  gut.Ptr(""),
	})
	if err != nil {
		return gut.Err(false, "failed to create collection question", err)
	}

	// * response
	return c.JSON(response.Success(c, "successfully created collection question"))
}
