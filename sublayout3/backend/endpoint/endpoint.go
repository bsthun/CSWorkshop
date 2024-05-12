package endpoint

import (
	"github.com/gofiber/fiber/v2"

	previewEndpoint "backend/endpoint/preview"
)

func Init(router fiber.Router) {
	preview := router.Group("/preview")
	preview.Get("/state", previewEndpoint.GetStateHandler)
}
