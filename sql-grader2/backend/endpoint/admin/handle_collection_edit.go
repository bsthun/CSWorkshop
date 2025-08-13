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

func (r *Handler) HandleCollectionEdit(c *fiber.Ctx) error {
	// * get user claims
	u := c.Locals("l").(*jwt.Token).Claims.(*common.LoginClaims)

	// * parse body
	body := new(payload.CollectionEditRequest)
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

	// * update collection
	collection, err := r.database.P().CollectionEdit(c.Context(), &psql.CollectionEditParams{
		Id:   body.Id,
		Name: body.Name,
	})
	if err != nil {
		return gut.Err(false, "failed to update collection", err)
	}

	// * response
	return c.JSON(response.Success(c, &payload.Collection{
		Id:        collection.Id,
		Name:      collection.Name,
		Metadata:  collection.Metadata,
		CreatedAt: collection.CreatedAt,
		UpdatedAt: collection.UpdatedAt,
	}))
}