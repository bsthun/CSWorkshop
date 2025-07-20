package middleware

import (
	"backend/type/common"
	"github.com/bsthun/gut"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/jwt/v3"
)

func (r *Middleware) Jwt(require bool) fiber.Handler {
	conf := jwtware.Config{
		SigningKey:  []byte(*r.config.Secret),
		TokenLookup: "cookie:login",
		ContextKey:  "l",
		Claims:      new(common.LoginClaims),
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if !require {
				return c.Next()
			}
			return gut.Err(false, "jwt validation failure", err)
		},
	}

	return jwtware.New(conf)
}
