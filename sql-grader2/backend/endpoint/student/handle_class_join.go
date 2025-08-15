package studentEndpoint

import (
	"backend/generate/psql"
	"backend/type/common"
	"backend/type/payload"
	"backend/type/response"
	"database/sql"

	"github.com/bsthun/gut"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func (r *Handler) HandleClassJoin(c *fiber.Ctx) error {
	// * get user claims
	u := c.Locals("l").(*jwt.Token).Claims.(*common.LoginClaims)

	// * parse body
	body := new(payload.ClassJoinRequest)
	if err := c.BodyParser(body); err != nil {
		return gut.Err(false, "invalid body", err)
	}

	// * get class by register code
	class, err := r.database.P().ClassGetByRegisterCode(c.Context(), body.RegisterCode)
	if err != nil {
		if err == sql.ErrNoRows {
			return gut.Err(false, "class not found", nil)
		}
		return gut.Err(false, "failed to get class", err)
	}

	// * check if student is already joined
	_, err = r.database.P().ClassJoineeCreate(c.Context(), &psql.ClassJoineeCreateParams{
		UserId:  u.UserId,
		ClassId: class.Id,
	})
	if err != nil {
		return gut.Err(false, "failed to join class or already joined", err)
	}

	// * response
	return c.JSON(response.Success(c, "successfully joined class"))
}
