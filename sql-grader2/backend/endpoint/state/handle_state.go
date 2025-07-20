package stateEndpoint

import (
	"backend/type/common"
	"backend/type/payload"
	"backend/type/response"
	"github.com/bsthun/gut"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func (r *Handler) HandleState(c *fiber.Ctx) error {
	// * login claims
	l := c.Locals("l").(*jwt.Token).Claims.(*common.LoginClaims)

	// * get user from database
	user, err := r.database.P().UserGetById(c.Context(), l.UserId)
	if err != nil {
		return gut.Err(false, "failed to get user", err)
	}

	// * response
	return c.JSON(response.Success(c, &payload.StateResponse{
		UserId:      user.Id,
		DisplayName: user.Firstname,
		Email:       user.Email,
		PictureUrl:  user.PictureUrl,
		IsAdmin:     user.IsAdmin,
	}))
}
