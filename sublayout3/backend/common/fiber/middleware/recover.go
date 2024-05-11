package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/sirupsen/logrus"

	"backend/common"
)

func Recover() fiber.Handler {
	if *cc.Config.Environment == 1 {
		return func(c *fiber.Ctx) error {
			defer logrus.Debug("CALL " + c.Method() + " " + c.Path() + " " + string(c.Request().Body()))
			return c.Next()
		}
	}

	return recover.New()
}
