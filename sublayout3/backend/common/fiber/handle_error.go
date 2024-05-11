package fiber

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"backend/type/response"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	// Case of *fiber.Error
	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		return c.Status(fiberErr.Code).JSON(response.ErrorResponse{
			Success: false,
			Code:    strings.ReplaceAll(strings.ToUpper(fiberErr.Error()), " ", "_"),
			Message: fiberErr.Error(),
			Error:   fiberErr.Error(),
		})
	}

	// Case of ErrorInstance
	var respErr *response.ErrorInstance
	if errors.As(err, &respErr) {
		if respErr.Err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(&response.ErrorResponse{
				Success: false,
				Code:    respErr.Code,
				Message: respErr.Message,
				Error:   respErr.Err.Error(),
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(&response.ErrorResponse{
			Success: false,
			Message: respErr.Message,
			Code:    respErr.Code,
		})
	}

	// Case of validator.ValidationErrors
	var valErr validator.ValidationErrors
	if errors.As(err, &valErr) {
		var lists []string
		for _, err := range valErr {
			lists = append(lists, err.Field()+" ("+err.Tag()+")")
		}

		message := strings.Join(lists[:], ", ")

		return c.Status(fiber.StatusBadRequest).JSON(&response.ErrorResponse{
			Success: false,
			Code:    "VALIDATION_FAILED",
			Message: "VALIDATION failed on field " + message,
			Error:   valErr.Error(),
		})
	}

	return c.Status(fiber.StatusInternalServerError).JSON(&response.ErrorResponse{
		Success: false,
		Code:    "SERVER_SIDE_ERROR",
		Message: "Unknown server side error",
		Error:   err.Error(),
	})
}
