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

func (r *Handler) HandleCollectionQuestionEdit(c *fiber.Ctx) error {
	// * get user claims
	u := c.Locals("l").(*jwt.Token).Claims.(*common.LoginClaims)

	// * parse body
	body := new(payload.CollectionQuestionEditRequest)
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

	// * update collection question
	_, err = r.database.P().CollectionQuestionEdit(c.Context(), &psql.CollectionQuestionEditParams{
		Id:          body.Id,
		Title:       body.Title,
		Description: body.Description,
		CheckQuery:  body.CheckQuery,
		CheckPrompt: body.CheckPrompt,
	})
	if err != nil {
		return gut.Err(false, "failed to update collection question", err)
	}

	// * response
	return c.JSON(response.Success(c, "successfully updated collection question"))
}
