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

func (r *Handler) HandleCollectionList(c *fiber.Ctx) error {
	// * get user claims
	u := c.Locals("l").(*jwt.Token).Claims.(*common.LoginClaims)

	// * parse body
	body := new(payload.CollectionListRequest)
	if err := c.BodyParser(body); err != nil {
		return gut.Err(false, "invalid body", err)
	}

	// * validate body
	if err := gut.Validate(body); err != nil {
		return err
	}

	// * set default values
	if body.Limit == nil {
		body.Limit = gut.Ptr(int32(10))
	}
	if body.Offset == nil {
		body.Offset = gut.Ptr(int32(0))
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

	// * count collections
	count, err := r.database.P().CollectionCount(c.Context(), body.Name)
	if err != nil {
		return gut.Err(false, "failed to count collections", err)
	}

	// * list collections
	collections, err := r.database.P().CollectionList(c.Context(), &psql.CollectionListParams{
		Limit:  body.Limit,
		Offset: body.Offset,
		Name:   body.Name,
		Sort:   gut.Ptr("createdAt"),
	})
	if err != nil {
		return gut.Err(false, "failed to list collections", err)
	}

	// * map to response
	collectionPayloads, _ := gut.Iterate(collections, func(collection psql.CollectionListRow) (*payload.Collection, *gut.ErrorInstance) {
		return &payload.Collection{
			Id:            collection.Id,
			Name:          collection.Name,
			Metadata:      collection.Metadata,
			QuestionCount: collection.QuestionCount,
			CreatedAt:     collection.CreatedAt,
			UpdatedAt:     collection.UpdatedAt,
		}, nil
	})

	// * response
	return c.JSON(response.Success(c, &payload.CollectionListResponse{
		Count:       count,
		Collections: collectionPayloads,
	}))
}
