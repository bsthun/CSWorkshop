package publicEndpoint

import (
	"backend/type/payload"
	"backend/type/response"
	"github.com/bsthun/gut"
	"github.com/gofiber/fiber/v2"
)

func (r *Handler) HandleLoginRedirect(c *fiber.Ctx) error {
	// * construct redirect url
	redirectUrl := r.Oauth2Config.AuthCodeURL(*gut.Random(gut.RandomSet.MixedAlphaNum, 8))

	return c.JSON(response.Success(c, &payload.OauthRedirect{
		RedirectUrl: &redirectUrl,
	}))
}
