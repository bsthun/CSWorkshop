package endpoint

import (
	"backend/common/config"
	"backend/common/fiber/middleware"
	"backend/endpoint/admin"
	"backend/endpoint/public"
	"backend/endpoint/state"
	"github.com/gofiber/fiber/v2"
	"path/filepath"
)

func Bind(
	app *fiber.App,
	publicEndpoint *publicEndpoint.Handler,
	stateEndpoint *stateEndpoint.Handler,
	adminEndpoint *adminEndpoint.Handler,
	middleware *middleware.Middleware,
	config *config.Config,
) {
	api := app.Group("/api")
	api.Use(middleware.Id())

	// * public endpoints
	public := api.Group("/public")
	public.Get("/login/redirect", publicEndpoint.HandleLoginRedirect)
	public.Post("/login/callback", publicEndpoint.HandleLoginCallback)

	// * state endpoints
	state := api.Group("/state", middleware.Jwt(true))
	state.Post("/state", stateEndpoint.HandleState)

	// * admin endpoints
	admin := api.Group("/admin", middleware.Jwt(true))
	_ = admin

	// * static files
	app.Static("/file", ".local/file")

	// * static
	app.Static("/", *config.WebRoot)
	app.Get("/*", func(c *fiber.Ctx) error {
		return c.SendFile(filepath.Join(*config.WebRoot, "index.html"))
	})
}
