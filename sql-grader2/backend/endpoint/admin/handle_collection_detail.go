package adminEndpoint

import (
	"backend/type/common"
	"backend/type/payload"
	"backend/type/response"

	"github.com/bsthun/gut"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func (r *Handler) HandleCollectionDetail(c *fiber.Ctx) error {
	// * get user claims
	u := c.Locals("l").(*jwt.Token).Claims.(*common.LoginClaims)

	// * parse body
	body := new(payload.CollectionIdRequest)
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

	// * get collection detail
	collectionDetail, err := r.database.P().CollectionDetail(c.Context(), body.CollectionId)
	if err != nil {
		return gut.Err(false, "failed to get collection", err)
	}

	// * map to response
	responsePayload := &payload.CollectionDetailResponse{
		Collection: &payload.Collection{
			Id:            collectionDetail.Id,
			Name:          collectionDetail.Name,
			Metadata:      collectionDetail.Metadata,
			QuestionCount: collectionDetail.QuestionCount,
			CreatedAt:     collectionDetail.CreatedAt,
			UpdatedAt:     collectionDetail.UpdatedAt,
		},
	}

	// * response
	return c.JSON(response.Success(c, responsePayload))
}
