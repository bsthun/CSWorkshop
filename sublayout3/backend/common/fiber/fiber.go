package fiber

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"backend/common"
	"backend/common/fiber/middleware"
	"backend/endpoint"
	"backend/type/response"
	"backend/util/text"
)

func Init() {
	// Initialize fiber instance
	app := fiber.New(fiber.Config{
		AppName:       "CSC105 Prehack 2024: Sublayout [" + text.Commit + "]",
		ErrorHandler:  ErrorHandler,
		Prefork:       false,
		StrictRouting: true,
		Network:       fiber.NetworkTCP,
	})

	// Register root endpoint
	app.All("/", func(c *fiber.Ctx) error {
		return c.JSON(response.Info("Pre-hackathon API ROOT"))
	})

	// Register API endpoints
	apiGroup := app.Group("api/")
	apiGroup.Use(middleware.Recover())
	apiGroup.Use(middleware.Cors())
	endpoint.Init(apiGroup)

	// Register not found endpoint
	app.Use(NotFoundHandler)

	// Startup
	err := app.Listen(*cc.Config.Address)
	if err != nil {
		log.Fatal("Unable to start fiber instance", err)
	}
}
